package pair

import (
	"strings"

	"github.com/thrasher-/gocryptotrader/common"
)

// CurrencyItem is an exported string with methods to manipulate the data instead
// of using array/slice access modifiers
type CurrencyItem string

// Lower converts the CurrencyItem object c to lowercase
func (c CurrencyItem) Lower() CurrencyItem {
	return CurrencyItem(strings.ToLower(string(c)))
}

// Upper converts the CurrencyItem object c to uppercase
func (c CurrencyItem) Upper() CurrencyItem {
	return CurrencyItem(strings.ToUpper(string(c)))
}

// String converts the CurrencyItem object c to string
func (c CurrencyItem) String() string {
	return string(c)
}

// CurrencyPair holds currency pair information
type CurrencyPair struct {
	Delimiter      string       `json:"delimiter"`
	FirstCurrency  CurrencyItem `json:"first_currency"`
	SecondCurrency CurrencyItem `json:"second_currency"`
}

// Pair returns a currency pair string
func (c CurrencyPair) Pair() CurrencyItem {
	return c.FirstCurrency + CurrencyItem(c.Delimiter) + c.SecondCurrency
}

// Display formats and returns the currency based on user preferences,
// overriding the default Pair() display
func (c CurrencyPair) Display(delimiter string, uppercase bool) CurrencyItem {
	var pair CurrencyItem

	if delimiter != "" {
		pair = c.FirstCurrency + CurrencyItem(delimiter) + c.SecondCurrency
	} else {
		pair = c.FirstCurrency + c.SecondCurrency
	}

	if uppercase {
		return pair.Upper()
	}
	return pair.Lower()
}

// Equal compares two currency pairs and returns whether or not they are equal
func (c CurrencyPair) Equal(p CurrencyPair, exact bool) bool {
	if !exact {
		if c.FirstCurrency.Upper() == p.FirstCurrency.Upper() &&
			c.SecondCurrency.Upper() == p.SecondCurrency.Upper() ||
			c.FirstCurrency.Upper() == p.SecondCurrency.Upper() &&
				c.SecondCurrency.Upper() == p.FirstCurrency.Upper() {
			return true
		}
	} else {
		if c.FirstCurrency.Upper() == p.FirstCurrency.Upper() &&
			c.SecondCurrency.Upper() == p.SecondCurrency.Upper() {
			return true
		}
	}
	return false
}

// Swap swaps the pairs first and second currencies
func (c CurrencyPair) Swap() CurrencyPair {
	p := c
	p.FirstCurrency = c.SecondCurrency
	p.SecondCurrency = c.FirstCurrency
	return p
}

// NewCurrencyPairDelimiter splits the desired currency string at delimeter,
// the returns a CurrencyPair struct
func NewCurrencyPairDelimiter(currency, delimiter string) CurrencyPair {
	result := strings.Split(currency, delimiter)
	return CurrencyPair{
		Delimiter:      delimiter,
		FirstCurrency:  CurrencyItem(result[0]),
		SecondCurrency: CurrencyItem(result[1]),
	}
}

// NewCurrencyPair returns a CurrencyPair without a delimiter
func NewCurrencyPair(firstCurrency, secondCurrency string) CurrencyPair {
	return CurrencyPair{
		FirstCurrency:  CurrencyItem(firstCurrency),
		SecondCurrency: CurrencyItem(secondCurrency),
	}
}

// NewCurrencyPairFromIndex returns a CurrencyPair via a currency string and
// specific index
func NewCurrencyPairFromIndex(currency, index string) CurrencyPair {
	i := strings.Index(currency, index)
	if i == 0 {
		return NewCurrencyPair(currency[0:len(index)], currency[len(index):])
	}
	return NewCurrencyPair(currency[0:i], currency[i:])
}

// NewCurrencyPairFromString converts currency string into a new CurrencyPair
// with or without delimeter
func NewCurrencyPairFromString(currency string) CurrencyPair {
	delimiters := []string{"_", "-"}
	var delimiter string
	for _, x := range delimiters {
		if strings.Contains(currency, x) {
			delimiter = x
			return NewCurrencyPairDelimiter(currency, delimiter)
		}
	}
	return NewCurrencyPair(currency[0:3], currency[3:])
}

// Contains checks to see if a specified pair exists inside a currency pair
// array
func Contains(pairs []CurrencyPair, p CurrencyPair, exact bool) bool {
	for x := range pairs {
		if pairs[x].Equal(p, exact) {
			return true
		}
	}
	return false
}

// ContainsCurrency checks to see if a pair contains a specific currency
func ContainsCurrency(p CurrencyPair, c string) bool {
	return p.FirstCurrency.Upper().String() == common.StringToUpper(c) ||
		p.SecondCurrency.Upper().String() == common.StringToUpper(c)
}

// RemovePairsByFilter checks to see if a pair contains a specific currency
// and removes it from the list of pairs
func RemovePairsByFilter(p []CurrencyPair, filter string) []CurrencyPair {
	var pairs []CurrencyPair
	for x := range p {
		if ContainsCurrency(p[x], filter) {
			continue
		}
		pairs = append(pairs, p[x])
	}
	return pairs
}

// FormatPairs formats a string array to a list of currency pairs with the
// supplied currency pair format
func FormatPairs(pairs []string, delimiter, index string) []CurrencyPair {
	var result []CurrencyPair
	for x := range pairs {
		if pairs[x] == "" {
			continue
		}
		var p CurrencyPair
		if delimiter != "" {
			p = NewCurrencyPairDelimiter(pairs[x], delimiter)
		} else {
			if index != "" {
				p = NewCurrencyPairFromIndex(pairs[x], index)
			} else {
				p = NewCurrencyPair(pairs[x][0:3], pairs[x][3:])
			}
		}
		result = append(result, p)
	}
	return result
}

// CopyPairFormat copies the pair format from a list of pairs once matched
func CopyPairFormat(p CurrencyPair, pairs []CurrencyPair, exact bool) CurrencyPair {
	for x := range pairs {
		if p.Equal(pairs[x], exact) {
			return pairs[x]
		}
	}
	return CurrencyPair{}
}

// FindPairDifferences returns pairs which are new or have been removed
func FindPairDifferences(oldPairs, newPairs []string) ([]string, []string) {
	var newPs, removedPs []string
	for x := range newPairs {
		if newPairs[x] == "" {
			continue
		}
		if !common.StringDataCompareUpper(oldPairs, newPairs[x]) {
			newPs = append(newPs, newPairs[x])
		}
	}
	for x := range oldPairs {
		if oldPairs[x] == "" {
			continue
		}
		if !common.StringDataCompareUpper(newPairs, oldPairs[x]) {
			removedPs = append(removedPs, oldPairs[x])
		}
	}
	return newPs, removedPs
}

// PairsToStringArray returns a list of pairs as a string array
func PairsToStringArray(pairs []CurrencyPair) []string {
	var p []string
	for x := range pairs {
		p = append(p, pairs[x].Pair().String())
	}
	return p
}
