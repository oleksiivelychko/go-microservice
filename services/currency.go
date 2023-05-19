package services

import (
	"context"
	"fmt"
	"github.com/oleksiivelychko/go-grpc-service/logger"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"github.com/oleksiivelychko/go-microservice/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Currency struct {
	logger           *logger.Logger
	exchangerClient  grpcservice.ExchangerClient
	currency         string
	cachedRates      map[string]float64
	subscriberClient grpcservice.Exchanger_SubscriberClient
}

func NewCurrency(exchangerClient grpcservice.ExchangerClient, currency string, logger *logger.Logger) *Currency {
	service := &Currency{
		logger,
		exchangerClient,
		currency,
		make(map[string]float64),
		nil,
	}

	go service.handleUpdates()
	return service
}

func (service *Currency) GetRate() (float64, *errors.GRPCServiceError) {
	exchangeRequest := &grpcservice.ExchangeRequest{
		From: grpcservice.Currencies_EUR,
		To:   grpcservice.Currencies(grpcservice.Currencies_value[service.currency]),
	}

	exchangeResponse, err := service.exchangerClient.MakeExchange(context.Background(), exchangeRequest)
	if err != nil {
		// convert the gRPC error message
		grpcErr, ok := status.FromError(err)
		if !ok {
			return -1, &errors.GRPCServiceError{Message: err.Error()}
		}

		if grpcErr.Code() == codes.InvalidArgument {
			return -1, &errors.GRPCServiceError{
				Message: fmt.Sprintf("unable to retrive exchange request from gRPC server: %s", grpcErr.Message()),
			}
		}
	}

	service.logResponse(exchangeResponse)
	service.cachedRates[service.currency] = exchangeResponse.GetRate()

	// subscribe for updates
	if err = service.subscriberClient.Send(exchangeRequest); err != nil {
		service.logger.Error("unable to send exchange request: %s", err)
	}

	return exchangeResponse.GetRate(), nil
}

func (service *Currency) SetCurrency(currency string) {
	service.currency = currency
}

func (service *Currency) handleUpdates() {
	subscribedClient, err := service.exchangerClient.Subscriber(context.Background())
	if err != nil {
		service.logger.Error("unable to subscribe for updates: %s", err)
	}

	service.subscriberClient = subscribedClient

	for {
		streamExchangeResponse, recvErr := subscribedClient.Recv()
		if grpcErr := streamExchangeResponse.GetError(); grpcErr != nil {
			service.logger.Error("grpcservice.Exchanger_SubscriberClient: %s", grpcErr)
			continue
		}

		if exchangeResponse := streamExchangeResponse.GetExchangeResponse(); exchangeResponse != nil {
			if recvErr != nil {
				service.logger.Error("unable to receive the message: %s", recvErr)
				return
			}

			service.logResponse(exchangeResponse)
			service.cachedRates[exchangeResponse.GetTo().String()] = exchangeResponse.GetRate()
		}
	}
}

func (service *Currency) logResponse(response *grpcservice.ExchangeResponse) {
	service.logger.Info("got gRPC response: from=%s, to=%s, rate=%f, createdAt=%s",
		response.GetFrom(),
		response.GetTo(),
		response.GetRate(),
		response.GetCreatedAt().AsTime().Format("2006-01-02"),
	)
}
