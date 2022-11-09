package service

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CurrencyService struct {
	logger      hclog.Logger
	client      grpc_service.CurrencyClient
	currency    string
	cachedRates map[string]float64
	subscriber  grpc_service.Currency_SubscriberClient
}

func NewCurrencyService(l hclog.Logger, cc grpc_service.CurrencyClient, c string) *CurrencyService {
	cs := &CurrencyService{l, cc, c, make(map[string]float64), nil}
	go cs.handleUpdates()
	return cs
}

func (cs *CurrencyService) GetRate() (float64, error) {
	exchangeRequest := &grpc_service.ExchangeRequest{
		From: grpc_service.Currencies_EUR,
		To:   grpc_service.Currencies(grpc_service.Currencies_value[cs.currency]),
	}

	exchangeResponse, err := cs.client.MakeExchange(context.Background(), exchangeRequest)
	if err != nil {
		// convert the gRPC error message
		grpcErr, ok := status.FromError(err)
		if !ok {
			return -1, err
		}

		if grpcErr.Code() == codes.InvalidArgument {
			return -1, fmt.Errorf("unable to retrive exchange request from gRPC server: '%s'", grpcErr.Message())
		}
	}

	cs.logResponse(exchangeResponse)
	cs.cachedRates[cs.currency] = exchangeResponse.GetRate()

	// subscribe for updates
	if err = cs.subscriber.Send(exchangeRequest); err != nil {
		cs.logger.Error("unable to send exchange request", "error", err)
	}

	return exchangeResponse.GetRate(), nil
}

func (cs *CurrencyService) SetCurrency(currency string) {
	cs.currency = currency
}

func (cs *CurrencyService) handleUpdates() {
	subscribedClient, err := cs.client.Subscriber(context.Background())
	if err != nil {
		cs.logger.Error("unable to subscribe for updates", "error", err)
	}

	cs.subscriber = subscribedClient

	for {
		streamExchangeResponse, recvErr := subscribedClient.Recv()
		if grpcErr := streamExchangeResponse.GetError(); grpcErr != nil {
			cs.logger.Error("subscriber", "error", grpcErr)
			continue
		}

		if exchangeResponse := streamExchangeResponse.GetExchangeResponse(); exchangeResponse != nil {
			if recvErr != nil {
				cs.logger.Error("unable to receive the message", "error", recvErr)
				return
			}

			cs.logResponse(exchangeResponse)
			cs.cachedRates[exchangeResponse.GetTo().String()] = exchangeResponse.GetRate()
		}
	}
}

func (cs *CurrencyService) logResponse(response *grpc_service.ExchangeResponse) {
	cs.logger.Info("got gRPC response",
		"from", response.GetFrom(),
		"to", response.GetTo(),
		"rate", response.GetRate(),
		"createdAt", response.GetCreatedAt().AsTime().Format("2006-01-02"),
	)
}
