package service

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
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
	if rate, ok := cs.cachedRates[cs.currency]; ok {
		return rate, nil
	}

	exchangeRequest := &grpc_service.ExchangeRequest{
		From: grpc_service.Currencies_EUR,
		To:   grpc_service.Currencies(grpc_service.Currencies_value[cs.currency]),
	}

	// get initial rate
	response, err := cs.client.MakeExchange(context.Background(), exchangeRequest)
	if err != nil {
		if grpcStatus, ok := status.FromError(err); ok {
			details := grpcStatus.Details()[0].(*grpc_service.ExchangeRequest)
			if grpcStatus.Code() == codes.InvalidArgument {
				return -1, fmt.Errorf(
					"unable to get exchange request from gRPC server, base '%s' and destination '%s' cannot be the same",
					details.GetFrom().String(),
					details.GetTo().String(),
				)
			}
		}

		return -1, err
	}

	cs.logger.Info("got gRPC response", "rate", response.Rate, "createdAt", response.CreatedAt)
	cs.cachedRates[cs.currency] = response.Rate

	// subscribe for updates
	cs.subscriber.Send(exchangeRequest)

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
		exchangeResponse, recvErr := subscribedClient.Recv()
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
