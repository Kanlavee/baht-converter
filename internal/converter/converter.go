// Package converter provides financial-grade Thai Baht text conversion,
// designed for use inside a service (e.g., invoice or payment systems).
package converter

import (
	"errors"
	"strings"

	"github.com/shopspring/decimal"
)

// ErrAmountTooLarge is returned when the absolute value of the amount
// exceeds the maximum supported value (9,999,999,999,999.99 baht).
var ErrAmountTooLarge = errors.New("amount exceeds maximum supported value (9,999,999,999,999.99)")

// maxAmount caps inputs to prevent silent int64 overflow in IntPart().
// ~10 trillion baht covers all realistic Thai financial transactions.
var maxAmount = decimal.RequireFromString("9999999999999.99")



// thaiDigits maps digits 1–9 to their formal Thai words.
// Index 0 is intentionally empty — zero never appears mid-number in Thai.
var thaiDigits = []string{
	"", "หนึ่ง", "สอง", "สาม", "สี่", "ห้า", "หก", "เจ็ด", "แปด", "เก้า",
}

// Thai positional magnitude words.
const (
	unitMillion  = "ล้าน"
	unitHundredK = "แสน"
	unitTenK     = "หมื่น"
	unitThousand = "พัน"
	unitHundred  = "ร้อย"
	unitTen      = "สิบ"
)

// Thai special-case and currency words.
const (
	// wordYi: formal Thai uses "ยี่สิบ" for 20–29, not "สองสิบ".
	wordYi = "ยี่สิบ"

	// wordEt: when ones digit is 1 AND higher-magnitude digits precede it,
	// Thai formal grammar requires "เอ็ด" instead of "หนึ่ง".
	// e.g., 11 → "สิบเอ็ด", 21 → "ยี่สิบเอ็ด", 101 → "หนึ่งร้อยเอ็ด"
	wordEt = "เอ็ด"

	wordBaht     = "บาท"
	wordThuean   = "ถ้วน" // exact suffix: no satang
	wordSatang   = "สตางค์"
	wordZeroBaht = "ศูนย์บาทถ้วน"
	wordNegative = "ลบ"
)

// Converter converts decimal.Decimal amounts to formal Thai Baht text.
// Use NewConverter() to obtain an instance.
type Converter struct{}

// NewConverter returns a ready-to-use Converter.
func NewConverter() *Converter {
	return &Converter{}
}

// Convert converts a decimal.Decimal amount to formal Thai Baht text.
func (c *Converter) Convert(amount decimal.Decimal) (string, error) {
	if err := validateAmount(amount); err != nil {
		return "", err
	}

	// Zero returns immediately without further processing.
	if amount.IsZero() {
		return wordZeroBaht, nil
	}

	isNegative, normalized := extractSign(amount)
	result := formatBahtText(normalized)

	if isNegative {
		return wordNegative + result, nil
	}
	return result, nil
}

// validateAmount returns ErrAmountTooLarge when the amount is outside
// the supported range, preventing silent int64 overflow downstream.
func validateAmount(amount decimal.Decimal) error {
	if amount.Abs().GreaterThan(maxAmount) {
		return ErrAmountTooLarge
	}
	return nil
}

// extractSign returns whether the amount is negative and its absolute value.
func extractSign(amount decimal.Decimal) (isNegative bool, normalized decimal.Decimal) {
	if amount.IsNegative() {
		return true, amount.Abs()
	}
	return false, amount
}

// splitAmount decomposes a positive amount into its integer value and satang (0–99).
func splitAmount(amount decimal.Decimal) (integerValue int64, satangValue int64) {
	integerPart := amount.Floor()
	integerValue = integerPart.IntPart()

	// Fractional part is truncated to 2 decimal places (no rounding).
	satangValue = amount.Sub(integerPart).Mul(decimal.NewFromInt(100)).Floor().IntPart()

	// Defensive clamp: satang must always be in [0, 99].
	if satangValue < 0 {
		satangValue = 0
	} else if satangValue > 99 {
		satangValue = 99
	}
	return integerValue, satangValue
}

// formatBahtText assembles the complete Thai currency string for a
// positive, validated, non-zero amount.
func formatBahtText(amount decimal.Decimal) string {
	integerValue, satangValue := splitAmount(amount)

	// Satang-only: integer part is zero but satang is non-zero.
	if integerValue == 0 {
		if satangValue == 0 {
			return wordZeroBaht
		}
		return convertIntegerPartToThai(satangValue, false) + wordSatang
	}

	var sb strings.Builder
	sb.WriteString(convertIntegerPartToThai(integerValue, false))
	sb.WriteString(wordBaht)
	if satangValue == 0 {
		sb.WriteString(wordThuean)
	} else {
		sb.WriteString(convertIntegerPartToThai(satangValue, false))
		sb.WriteString(wordSatang)
	}
	return sb.String()
}

// convertIntegerPartToThai converts a positive int64 to Thai text.
func convertIntegerPartToThai(n int64, hasPrefixDigits bool) string {
	if n == 0 {
		return ""
	}
	if n >= 1_000_000 {
		millions := n / 1_000_000
		remainder := n % 1_000_000

		var sb strings.Builder
		sb.WriteString(convertIntegerPartToThai(millions, hasPrefixDigits))
		sb.WriteString(unitMillion)
		if remainder > 0 {
			// Pass true: the lan group provides prefix digits for the remainder.
			sb.WriteString(convertSubMillionToThai(int(remainder), true))
		}
		return sb.String()
	}
	return convertSubMillionToThai(int(n), hasPrefixDigits)
}

// convertSubMillionToThai converts n (1–999,999) to Thai positional text.
func convertSubMillionToThai(n int, hasPrefixDigits bool) string {
	if n <= 0 {
		return ""
	}

	hundredThousands := (n / 100_000) % 10
	tenThousands := (n / 10_000) % 10
	thousands := (n / 1_000) % 10
	hundreds := (n / 100) % 10
	tens := (n / 10) % 10
	ones := n % 10

	var sb strings.Builder

	if hundredThousands > 0 {
		sb.WriteString(thaiDigits[hundredThousands])
		sb.WriteString(unitHundredK)
	}
	if tenThousands > 0 {
		sb.WriteString(thaiDigits[tenThousands])
		sb.WriteString(unitTenK)
	}
	if thousands > 0 {
		sb.WriteString(thaiDigits[thousands])
		sb.WriteString(unitThousand)
	}
	if hundreds > 0 {
		sb.WriteString(thaiDigits[hundreds])
		sb.WriteString(unitHundred)
	}

	switch {
	case tens == 1:
		sb.WriteString(unitTen) // "sip" — no digit prefix for 10–19
	case tens == 2:
		sb.WriteString(wordYi) // "yisip" — formal Thai for 20–29
	case tens > 2:
		sb.WriteString(thaiDigits[tens])
		sb.WriteString(unitTen)
	}

	if ones > 0 {
		// hasPrefixDigits from caller OR something already written in this call.
		sb.WriteString(formatOnesDigit(ones, hasPrefixDigits || sb.Len() > 0))
	}

	return sb.String()
}

// formatOnesDigit returns the Thai word for a ones digit (1–9).
func formatOnesDigit(digit int, hasPrefixDigits bool) string {
	if digit == 1 && hasPrefixDigits {
		return wordEt
	}
	return thaiDigits[digit]
}
