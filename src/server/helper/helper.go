package helper

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/doitmagic/convmic/src/server/internal"
)


func PopulateData(appcontext *internal.AppContext) int {

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
