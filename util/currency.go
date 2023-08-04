package util

const (
	USD = "USD"
	EUR = "EUR"
	CAD = "CAD"
	YEN = "YEN"
)

// IsSupportedCurrency return true if the currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, EUR, CAD, YEN:
		return true
	default:
		return false
	}
}
