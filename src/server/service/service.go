package service

import (
	"context"

	"github.com/doitmagic/convmic/src/server/model"
)

type Provider interface {
	PopulateProviderCurrencies(ctx context.Context) error
	Convert(ctx context.Context, from []model.CurrencyConvert, to string) ([]model.CurrencyConverted, error)
	SyncCurrencies(limitPage int) (bool, error)
}
