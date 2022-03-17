package testing

import (
	"fmt"
	"strconv"
	"unicode"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/doitmagic/convmic/src/server/internal"
)

var _ = Describe("Internal tests", func() {
	Context("test appcontext", func() {

		appcontext := internal.GetInstance()
		It("currencies variable must be initialisated", func() {
			currencies := appcontext.GetAllCurrencies()
			Expect(currencies).ToNot(BeNil())
		})

		It("you must add currency record to currencies variable", func() {

			//populate dummy data
			totalRecordsNr := populateData(appcontext)
			//get all currencies
			currenciesNr := appcontext.CountCurrencies()

			Expect(currenciesNr).To(Equal(totalRecordsNr))

		})

		It("expect correct pagination results", func() {
			paginatedCurrencies := appcontext.GetCurrenciesByPage(1, 10)
			Expect(len(paginatedCurrencies)).To(Equal(10))
		})

	})
})

func populateData(appcontext *internal.AppContext) int {

	count := 0
	for r := 'a'; r < 'g'; r++ {
		R := unicode.ToUpper(r)
		for j := 0; j < 11; j++ {
			appcontext.SetCurrency(fmt.Sprintf("%c", R)+"test"+strconv.Itoa(count+1), 100)
			count = count + 1
		}
	}

	return count
}
