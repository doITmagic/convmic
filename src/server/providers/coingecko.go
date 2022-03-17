package providers

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
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
	ctx    context.Context
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

func NewCoingeckoProvider(ctx context.Context) *CoingeckoProvider {
	log.Info("created new provider CoingeckoProvider")

	return &CoingeckoProvider{
		Name:   "coingecko",
		c:      internal.NewClient("", "", CoingeckoBaseURL),
		market: "usd",
		ctx:    ctx,
	}
}

func (p *CoingeckoProvider) Convert(ctx context.Context, from []model.CurrencyConvert, to string) ([]model.CurrencyConverted, error) {
	return []model.CurrencyConverted{}, nil
}

//SyncCurrencies get the currencies values with limited pages,
//the page limit is required to not exceed the limits set by the API provider
func (p *CoingeckoProvider) SyncCurrencies(limitPage int) (bool, error) {

	//get all currencies from provider
	coinList, err := p.coinsList()
	if err != nil {
		return false, err
	}

	if limitPage == 0 {
		limitPage = 10
	}

	//calculate total numbers of currencies page for provider
	totalMarketPageNr := len(*coinList) / 150

	//for each page, if is lower then limit page execute
	for i := 1; i <= totalMarketPageNr; i++ {
		//limit page number request because of free API request number per minute
		if i < limitPage {

			//request is executed in goroutines
			//no problem with concurrent access of context variable currencies 
			// because we use sync.map{}
			go func(j int) {
				var tIds []string
				var tmpCoinList CoinList

				start := (j - 1) * 150
				stop := start + 150

				tmpCoinList = (*coinList)[start:stop]
				for _, coinItem := range tmpCoinList {
					tIds = append(tIds, coinItem.ID)
				}

				coinsMarket, err := p.GetMarketCurrencies(tIds, strconv.Itoa(j))
				if err != nil {
					log.Error(err)
				}

				for _, v := range *coinsMarket {
					internal.GetInstance().SetCurrency(v.Name, v.CurrentPrice)
				}

			}(i)
		}
	}

	return true, nil
}

func (p *CoingeckoProvider) PopulateProviderCurrencies(ctx context.Context) error {
	coinList, err := p.coinsList()
	if err != nil {
		return err
	}

	appContext := internal.GetInstance()

	for _, curency := range *coinList {
		appContext.SetCurrency(curency.Name, 0)
	}

	return nil
}

func (p *CoingeckoProvider) GetMarketCurrencies(idsParam []string, page string) (*CoinsMarket, error) {

	params := url.Values{}
	params.Add("vs_currency", "usd")
	params.Add("ids", strings.Join(idsParam, ","))
	//params.Add("per_page", "150")
	//params.Add("page", page)
	params.Add("start", "0")
	params.Add("limit", "150")

	resp, err := p.c.MakeReq("coins/markets", params)
	if err != nil {
		return nil, err
	}
	var data *CoinsMarket
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// func (p *CoingeckoProvider) getMultiplePrices(ids []string, vsCurrencies []string) (*map[string]map[string]float32, error) {
// 	params := url.Values{}
// 	idsParam := strings.Join(ids[:], ",")
// 	vsCurrenciesParam := strings.Join(vsCurrencies[:], ",")

// 	params.Add("ids", idsParam)
// 	params.Add("vs_currencies", vsCurrenciesParam)

// 	resp, err := p.c.MakeReq("simple/price", params)
// 	if err != nil {
// 		return nil, err
// 	}

// 	t := make(map[string]map[string]float32)
// 	err = json.Unmarshal(resp, &t)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &t, nil
// }

func (p *CoingeckoProvider) coinsList() (*CoinList, error) {

	params := url.Values{}
	params.Add("include_platform", "false")

	resp, err := p.c.MakeReq("coins/list", params)
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
