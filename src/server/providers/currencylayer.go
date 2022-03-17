package providers

import (
	"context"

	"github.com/doitmagic/convmic/src/server/internal"
	"github.com/doitmagic/convmic/src/server/model"
)

const (
	CurrencylayerBaseURL = "https://api.coingecko.com/api/v3/"
)

type CurrencylayerProvider struct {
	Name   string
	c      *internal.Client
	market string
	ctx    context.Context
}

func NewCurrencylayerProvider(ctx context.Context, APIKey, SecretKey string) *CoingeckoProvider {
	return &CoingeckoProvider{
		Name:   "Currencylayer",
		c:      internal.NewClient(APIKey, SecretKey, CurrencylayerBaseURL),
		market: "usd",
	}
}

func (p *CurrencylayerProvider) Convert(ctx context.Context, from []model.CurrencyConvert, to string) ([]model.CurrencyConverted, error) {
	//must be implemented
	return []model.CurrencyConverted{}, nil
}

func (p *CurrencylayerProvider) SyncCurrencies(period, page int) (bool, error) {
	//must be implemented
	return true, nil
}

func (p *CurrencylayerProvider) PopulateProviderCurrencies(ctx context.Context) error {
	//must be implemented
	return nil
}
