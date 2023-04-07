package service

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpcservice"
	"github.com/oleksiivelychko/go-microservice/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CurrencyService struct {
	logger           hclog.Logger
	exchangerClient  grpcservice.ExchangerClient
	currency         string
	cachedRates      map[string]float64
	subscriberClient grpcservice.Exchanger_SubscriberClient
}

func NewCurrencyService(logger hclog.Logger, exchangerClient grpcservice.ExchangerClient, currency string) *CurrencyService {
	service := &CurrencyService{
		logger,
		exchangerClient,
		currency,
		make(map[string]float64),
		nil,
	}

	go service.handleUpdates()
	return service
}

func (service *CurrencyService) GetRate() (float64, *errors.GRPCServiceError) {
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
		service.logger.Error("unable to send exchange request", "error", err)
	}

	return exchangeResponse.GetRate(), nil
}

func (service *CurrencyService) SetCurrency(currency string) {
	service.currency = currency
}

func (service *CurrencyService) handleUpdates() {
	subscribedClient, err := service.exchangerClient.Subscriber(context.Background())
	if err != nil {
		service.logger.Error("unable to subscribe for updates", "error", err)
	}

	service.subscriberClient = subscribedClient

	for {
		streamExchangeResponse, recvErr := subscribedClient.Recv()
		if grpcErr := streamExchangeResponse.GetError(); grpcErr != nil {
			service.logger.Error("grpcservice.Exchanger_SubscriberClient", "error", grpcErr)
			continue
		}

		if exchangeResponse := streamExchangeResponse.GetExchangeResponse(); exchangeResponse != nil {
			if recvErr != nil {
				service.logger.Error("unable to receive the message", "error", recvErr)
				return
			}

			service.logResponse(exchangeResponse)
			service.cachedRates[exchangeResponse.GetTo().String()] = exchangeResponse.GetRate()
		}
	}
}

func (service *CurrencyService) logResponse(response *grpcservice.ExchangeResponse) {
	service.logger.Info("got gRPC response",
		"from", response.GetFrom(),
		"to", response.GetTo(),
		"rate", response.GetRate(),
		"createdAt", response.GetCreatedAt().AsTime().Format("2006-01-02"),
	)
}
