package service

import (
	"context"

	"github.com/doitmagic/convmic/src/server/model"
)

type Provider interface {
	GetCurrencies(ctx context.Context) ([]model.Currency, error)
	Convert(ctx context.Context, from []model.CurrencyConvert, to string) ([]model.CurrencyConverted, error)
	SyncCurrencies(period int) (bool,error)
}
