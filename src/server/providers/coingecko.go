package providers

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/doitmagic/convmic/src/server/internal"
	"github.com/doitmagic/convmic/src/server/model"
	log "github.com/sirupsen/logrus"
)

const (
	CoingeckoBaseURL = "https://api.coingecko.com/api/v3/"
)

type CoingeckoProvider struct {
	Name   string
	c      *internal.Client
	market string
}

type CoinsMarket []CoinsMarketItem

type coinBaseStruct struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type CoinsListItem struct {
	coinBaseStruct
}

type CoinList []CoinsListItem

type CoinsMarketItem struct {
	coinBaseStruct
	Image        string  `json:"image"`
	CurrentPrice float64 `json:"current_price"`
}

func NewCoingeckoProvider() *CoingeckoProvider {
	log.Info("created new provider CoingeckoProvider")
	return &CoingeckoProvider{
		Name:   "coingecko",
		c:      internal.NewClient("", "", CoingeckoBaseURL),
		market: "usd",
	}
}

func (p *CoingeckoProvider) Convert(ctx context.Context, from []model.CurrencyConvert, to string) ([]model.CurrencyConverted, error) {
	return []model.CurrencyConverted{}, nil
}

func (p *CoingeckoProvider) SyncCurrencies(period int) (bool, error) {
	//must be implemented
	return true, nil
}

func (p *CoingeckoProvider) GetCurrencies(ctx context.Context) ([]model.Currency, error) {
	log.Println(p.coinsList())
	//must be implemented
	return []model.Currency{}, nil
}

// func (p *CoingeckoProvider) GetMarketCurrencies(period int) (*CoinsMarket, error) {

// 	params := url.Values{}
// 	params.Add("vs_currencies", p.market)

// 	resp, err := p.c.MakeReq("coins/markets", params)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var data *CoinsMarket
// 	err = json.Unmarshal(resp, &data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }

func (p *CoingeckoProvider) getMultiplePrices(ids []string, vsCurrencies []string) (*map[string]map[string]float32, error) {
	params := url.Values{}
	idsParam := strings.Join(ids[:], ",")
	vsCurrenciesParam := strings.Join(vsCurrencies[:], ",")

	params.Add("ids", idsParam)
	params.Add("vs_currencies", vsCurrenciesParam)

	resp, err := p.c.MakeReq("simple/price", params)
	if err != nil {
		return nil, err
	}

	t := make(map[string]map[string]float32)
	err = json.Unmarshal(resp, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (p *CoingeckoProvider) coinsList() (*CoinList, error) {

	resp, err := p.c.MakeReq("coins/list", nil)
	if err != nil {
		return nil, err
	}

	var data *CoinList
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
