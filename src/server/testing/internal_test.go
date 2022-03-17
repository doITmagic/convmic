package testing

import (
	

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/doitmagic/convmic/src/server/helper"
	"github.com/doitmagic/convmic/src/server/internal"
)

//The test coverage is small, just to demonstrate how is implemented
var _ = Describe("Internal tests", func() {
	Context("test appcontext", func() {

		appcontext := internal.GetInstance()
		It("currencies variable must be initialisated", func() {
			currencies := appcontext.GetAllCurrencies()
			Expect(currencies).ToNot(BeNil())
		})

		It("you must add currency record to currencies variable", func() {

			//populate dummy data
			totalRecordsNr :=  helper.PopulateData(appcontext)
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

