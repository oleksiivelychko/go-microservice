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

func NewCurrencyService(logger hclog.Logger, client grpc_service.CurrencyClient, currency string) *CurrencyService {
	cs := &CurrencyService{logger, client, currency, make(map[string]float64), nil}

	go cs.handleUpdates()

	return cs
}

func (cs *CurrencyService) GetRate() (float64, error) {
	exchangeRequest := &grpc_service.ExchangeRequest{
		From: grpc_service.Currencies_EUR,
		To:   grpc_service.Currencies(grpc_service.Currencies_value[cs.currency]),
	}

	// get initial rate
	response, err := cs.client.MakeExchange(context.Background(), exchangeRequest)
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

	cs.logger.Info("got gRPC response", "rate", response.Rate, "createdAt", response.CreatedAt.AsTime().Format("2006-01-02"))
	cs.cachedRates[cs.currency] = response.Rate

	// subscribe for updates
	if err = cs.subscriber.Send(exchangeRequest); err != nil {
		cs.logger.Error("unable to send exchange request", "error", err)
	}

	return response.Rate, nil
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
			cs.logger.Error("internal subscriber error", "error", grpcErr)
			continue
		}

		if exchangeResponse := streamExchangeResponse.GetExchangeResponse(); exchangeResponse != nil {
			if recvErr != nil {
				cs.logger.Error("unable to receive the message", "error", recvErr)
				return
			}

			cs.logger.Info("received the update from server",
				"from", exchangeResponse.GetFrom(),
				"to", exchangeResponse.GetTo(),
				"rate", exchangeResponse.GetRate(),
				"createdAt", exchangeResponse.GetCreatedAt().AsTime().Format("2006-01-02"),
			)

			cs.cachedRates[exchangeResponse.GetTo().String()] = exchangeResponse.GetRate()
		}
	}
}
