package internal

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
)

type AppContext struct {
	currencies *sync.Map
}

//GetAllCurrencies return currencies ordered by name , no paginations
func (c *AppContext) GetAllCurrencies() *map[string]string {

	all := make(map[string]string)
	result := make(map[string]string)

	//a variable just to store name an order them
	names := make([]string, 0, len(all))

	//get from sync map
	c.currencies.Range(func(k, v interface{}) bool {
		names = append(names, k.(string))
		all[k.(string)] = v.(string)
		return true
	})

	sort.Strings(names) //sort by key

	//populate result
	for _, name := range names {
		result[name] = all[name]
	}

	return &result
}

//GetAllCurrencies return currencies ordered by name with pagination by page nr and perPage parameters
func (c *AppContext) GetCurrenciesByPage(page, perPage int32) map[string]string {

	all := make(map[string]string)

	//a variable just to store name an order them
	names := make([]string, 0, len(all))

	//get from sync map
	c.currencies.Range(func(k, v interface{}) bool {
		names = append(names, k.(string))
		all[k.(string)] = v.(string)
		return true
	})

	if len(names) > 0 {
		//paginate result
		return paginate(names, all, int(page), int(perPage))

	}

	return nil
}

//paginate private paginate functionality
func paginate(names []string, allCurrencies map[string]string, page, perPage int) map[string]string {

	result := make(map[string]string)

	//sort the slice of name
	sort.Strings(names)

	if page == 0 {
		page = 1
	}

	//calculate start and stop index
	start := (page - 1) * perPage
	stop := start + perPage

	if stop > len(names) {
		stop = len(names)
	}

	//add to map only the record between start and stop
	for _, name := range (names)[start:stop] {
		result[name] = (allCurrencies)[name]
	}

	return result
}

//GetCurrencyValue get the value for the named currency
func (c *AppContext) GetCurrencyValue(name string) (float64, error) {
	value, ok := c.currencies.Load(name)
	if ok {
		valueFloat, _ := strconv.ParseFloat(value.(string), 64)
		return valueFloat, nil
	}
	return 0, fmt.Errorf("the currency %s does not exist", name)
}

//SetCurrency set or update currency value
func (c *AppContext) SetCurrency(name string, val float64) {
	//convert float64 to string
	valueString := strconv.FormatFloat(val, 'f', 2, 64)
	c.currencies.Store(name, valueString)
}

//CountCurrencies count currecties length
func (c *AppContext) CountCurrencies() int {

	length := 0
	c.currencies.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length
}

var instance *AppContext
var once sync.Once

// GetInstance returns the same context every time
func GetInstance() *AppContext {
	once.Do(func() {
		instance = &AppContext{currencies: &sync.Map{}}
	})
	return instance
}
