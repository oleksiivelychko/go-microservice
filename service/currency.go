package service

import (
	"context"
	"github.com/oleksiivelychko/go-grpc-protobuf/proto/grpc_service"
)

type CurrencyService struct {
	client          grpc_service.CurrencyClient
	defaultCurrency string
}

func NewCurrencyService(client grpc_service.CurrencyClient, currency string) *CurrencyService {
	return &CurrencyService{client, currency}
}

func (cs *CurrencyService) GetRate() (float64, error) {
	er := &grpc_service.ExchangeRequest{
		From: grpc_service.Currencies_EUR.String(),
		To:   cs.defaultCurrency,
	}

	response, err := cs.client.MakeExchange(context.Background(), er)
	if err != nil {
		return -1, err
	}

	return response.Rate, nil
}
