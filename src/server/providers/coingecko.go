package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

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
	log.Info("Privider Convert")
	appContext := internal.GetInstance()
	toCurrencyValue, err := appContext.GetCurrencyValue(to)
	if err == nil {
		return []model.CurrencyConverted{}, fmt.Errorf("can not convert to currency %s because does not exist", to)
	}

	for _, currencyConvert := range from {
		value, err := appContext.GetCurrencyValue(currencyConvert.Name)
		if err == nil {
			fromCurrencyTotalValue := currencyConvert.Amount * value
			if fromCurrencyTotalValue > 0 {
				rez := fromCurrencyTotalValue / toCurrencyValue
				fmt.Printf("%v ammount of currency %s represent %v of currency %s  \n", currencyConvert.Amount, currencyConvert.Name, rez, to)
			}
		}
	}

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

				coinsMarket, err := p.GetMarketCurrencies(tIds, start, stop)
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

//PopulateProviderCurrencies add all currencies names to context currencies
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

//GetMarketCurrencies get all market currency values from start to stop records
func (p *CoingeckoProvider) GetMarketCurrencies(idsParam []string, start, stop int) (*CoinsMarket, error) {

	params := url.Values{}
	params.Add("vs_currency", "usd")
	//params.Add("ids", strings.Join(idsParam, ","))
	params.Add("start", strconv.Itoa(start))
	params.Add("limit", strconv.Itoa(stop))

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

//coinsList private method to list all currencies from provider
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
