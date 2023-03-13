package service

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	grpcService "github.com/oleksiivelychko/go-grpc-service/proto/grpc_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CurrencyService struct {
	logger                   hclog.Logger
	currencyClient           grpcService.CurrencyClient
	currency                 string
	cachedRates              map[string]float64
	currencySubscriberClient grpcService.Currency_SubscriberClient
}

func NewCurrencyService(
	logger hclog.Logger,
	currencyClient grpcService.CurrencyClient,
	currency string,
) *CurrencyService {
	service := &CurrencyService{
		logger,
		currencyClient,
		currency,
		make(map[string]float64),
		nil,
	}
	go service.handleUpdates()
	return service
}

func (currencyService *CurrencyService) GetRate() (float64, error) {
	exchangeRequest := &grpcService.ExchangeRequest{
		From: grpcService.Currencies_EUR,
		To:   grpcService.Currencies(grpcService.Currencies_value[currencyService.currency]),
	}

	exchangeResponse, err := currencyService.currencyClient.MakeExchange(context.Background(), exchangeRequest)
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

	currencyService.logResponse(exchangeResponse)
	currencyService.cachedRates[currencyService.currency] = exchangeResponse.GetRate()

	// subscribe for updates
	if err = currencyService.currencySubscriberClient.Send(exchangeRequest); err != nil {
		currencyService.logger.Error("unable to send exchange request", "error", err)
	}

	return exchangeResponse.GetRate(), nil
}

func (currencyService *CurrencyService) SetCurrency(currency string) {
	currencyService.currency = currency
}

func (currencyService *CurrencyService) handleUpdates() {
	subscribedClient, err := currencyService.currencyClient.Subscriber(context.Background())
	if err != nil {
		currencyService.logger.Error("unable to subscribe for updates", "error", err)
	}

	currencyService.currencySubscriberClient = subscribedClient

	for {
		streamExchangeResponse, recvErr := subscribedClient.Recv()
		if grpcErr := streamExchangeResponse.GetError(); grpcErr != nil {
			currencyService.logger.Error("currencySubscriberClient", "error", grpcErr)
			continue
		}

		if exchangeResponse := streamExchangeResponse.GetExchangeResponse(); exchangeResponse != nil {
			if recvErr != nil {
				currencyService.logger.Error("unable to receive the message", "error", recvErr)
				return
			}

			currencyService.logResponse(exchangeResponse)
			currencyService.cachedRates[exchangeResponse.GetTo().String()] = exchangeResponse.GetRate()
		}
	}
}

func (currencyService *CurrencyService) logResponse(response *grpcService.ExchangeResponse) {
	currencyService.logger.Info("got gRPC response",
		"from", response.GetFrom(),
		"to", response.GetTo(),
		"rate", response.GetRate(),
		"createdAt", response.GetCreatedAt().AsTime().Format("2006-01-02"),
	)
}
