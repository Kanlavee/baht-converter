package converter

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestConverter_Convert(t *testing.T) {
	c := NewConverter()

	tests := []struct {
		name    string
		input   string // use string to avoid float imprecision
		want    string
		wantErr bool
	}{
		// ── Required cases from specification ────────────────────────────────
		{name: "zero", input: "0", want: "ศูนย์บาทถ้วน"},
		{name: "one", input: "1", want: "หนึ่งบาทถ้วน"},
		{name: "eleven", input: "11", want: "สิบเอ็ดบาทถ้วน"},
		{name: "twenty_one", input: "21", want: "ยี่สิบเอ็ดบาทถ้วน"},
		{name: "one_hundred_one", input: "101", want: "หนึ่งร้อยเอ็ดบาทถ้วน"},
		{name: "one_thousand", input: "1000", want: "หนึ่งพันบาทถ้วน"},
		{name: "one_thousand_two_three_four", input: "1234", want: "หนึ่งพันสองร้อยสามสิบสี่บาทถ้วน"},
		{name: "thirty_three_three_three_seventy_five_satang", input: "33333.75", want: "สามหมื่นสามพันสามร้อยสามสิบสามบาทเจ็ดสิบห้าสตางค์"},
		{name: "negative_fifty_one", input: "-51", want: "ลบห้าสิบเอ็ดบาทถ้วน"},
		{name: "negative_zero_point_five", input: "-0.5", want: "ลบห้าสิบสตางค์"},

		// ── Tens special grammar ──────────────────────────────────────────────
		{name: "ten", input: "10", want: "สิบบาทถ้วน"},
		{name: "twenty", input: "20", want: "ยี่สิบบาทถ้วน"},
		{name: "nineteen", input: "19", want: "สิบเก้าบาทถ้วน"},
		{name: "ninety_nine", input: "99", want: "เก้าสิบเก้าบาทถ้วน"},

		// ── Hundreds / thousands ──────────────────────────────────────────────
		{name: "one_hundred", input: "100", want: "หนึ่งร้อยบาทถ้วน"},
		{name: "one_thousand_one", input: "1001", want: "หนึ่งพันเอ็ดบาทถ้วน"},
		{name: "ten_thousand", input: "10000", want: "หนึ่งหมื่นบาทถ้วน"},
		{name: "one_hundred_thousand", input: "100000", want: "หนึ่งแสนบาทถ้วน"},

		// ── Millions ──────────────────────────────────────────────────────────
		{name: "one_million", input: "1000000", want: "หนึ่งล้านบาทถ้วน"},
		{name: "one_million_one", input: "1000001", want: "หนึ่งล้านเอ็ดบาทถ้วน"},
		{name: "one_billion", input: "1000000000", want: "หนึ่งพันล้านบาทถ้วน"},

		// ── Satang handling ───────────────────────────────────────────────────
		{name: "satang_only_positive", input: "0.25", want: "ยี่สิบห้าสตางค์"},
		{name: "satang_only_one", input: "0.01", want: "หนึ่งสตางค์"},
		{name: "eleven_point_zero_one", input: "11.01", want: "สิบเอ็ดบาทหนึ่งสตางค์"},
		{name: "one_million_fifty_satang", input: "1000000.50", want: "หนึ่งล้านบาทห้าสิบสตางค์"},

		// ── Negative numbers ──────────────────────────────────────────────────
		{name: "negative_one", input: "-1", want: "ลบหนึ่งบาทถ้วน"},
		{name: "negative_satang_only", input: "-0.11", want: "ลบสิบเอ็ดสตางค์"},
		{name: "negative_with_satang", input: "-33333.75", want: "ลบสามหมื่นสามพันสามร้อยสามสิบสามบาทเจ็ดสิบห้าสตางค์"},

		// ── Edge cases ────────────────────────────────────────────────────────
		{name: "max_precision", input: "9999999.99", want: "เก้าล้านเก้าแสนเก้าหมื่นเก้าพันเก้าร้อยเก้าสิบเก้าบาทเก้าสิบเก้าสตางค์"},

		// ── Single digits ────────────────────────────────────────────────
		{name: "two", input: "2", want: "สองบาทถ้วน"},
		{name: "three", input: "3", want: "สามบาทถ้วน"},
		{name: "nine", input: "9", want: "เก้าบาทถ้วน"},

		// ── Teens edge ───────────────────────────────────────────────────
		{name: "twelve", input: "12", want: "สิบสองบาทถ้วน"},
		{name: "thirteen", input: "13", want: "สิบสามบาทถ้วน"},
		{name: "fifteen", input: "15", want: "สิบห้าบาทถ้วน"},

		// ── Tens grammar ─────────────────────────────────────────────────
		{name: "twenty_two", input: "22", want: "ยี่สิบสองบาทถ้วน"},
		{name: "twenty_nine", input: "29", want: "ยี่สิบเก้าบาทถ้วน"},
		{name: "thirty_one", input: "31", want: "สามสิบเอ็ดบาทถ้วน"},
		{name: "forty_one", input: "41", want: "สี่สิบเอ็ดบาทถ้วน"},
		{name: "fifty_one", input: "51", want: "ห้าสิบเอ็ดบาทถ้วน"},
		{name: "sixty_one", input: "61", want: "หกสิบเอ็ดบาทถ้วน"},
		{name: "seventy_one", input: "71", want: "เจ็ดสิบเอ็ดบาทถ้วน"},
		{name: "eighty_one", input: "81", want: "แปดสิบเอ็ดบาทถ้วน"},
		{name: "ninety_one", input: "91", want: "เก้าสิบเอ็ดบาทถ้วน"},

		// ── Hundreds ─────────────────────────────────────────────────────
		{name: "one_zero_one", input: "101", want: "หนึ่งร้อยเอ็ดบาทถ้วน"},
		{name: "one_one_one", input: "111", want: "หนึ่งร้อยสิบเอ็ดบาทถ้วน"},
		{name: "two_hundred_ten", input: "210", want: "สองร้อยสิบบาทถ้วน"},
		{name: "nine_hundred_ninety_nine", input: "999", want: "เก้าร้อยเก้าสิบเก้าบาทถ้วน"},

		// ── Thousands ────────────────────────────────────────────────────
		{name: "two_thousand", input: "2000", want: "สองพันบาทถ้วน"},
		{name: "two_thousand_twenty_one", input: "2021", want: "สองพันยี่สิบเอ็ดบาทถ้วน"},
		{name: "nine_thousand_nine", input: "9009", want: "เก้าพันเก้าบาทถ้วน"},
		{name: "ten_thousand_one", input: "10001", want: "หนึ่งหมื่นเอ็ดบาทถ้วน"},

		// ── Hundred thousands ────────────────────────────────────────────
		{name: "two_hundred_thousand", input: "200000", want: "สองแสนบาทถ้วน"},
		{name: "two_hundred_thousand_one", input: "200001", want: "สองแสนเอ็ดบาทถ้วน"},
		{name: "nine_hundred_ninety_nine_thousand", input: "999000", want: "เก้าแสนเก้าหมื่นเก้าพันบาทถ้วน"},

		// ── Millions complex ─────────────────────────────────────────────
		{name: "two_million", input: "2000000", want: "สองล้านบาทถ้วน"},
		{name: "two_million_twenty_one", input: "2000021", want: "สองล้านยี่สิบเอ็ดบาทถ้วน"},
		{name: "two_million_one_hundred", input: "2000100", want: "สองล้านหนึ่งร้อยบาทถ้วน"},
		{name: "nine_million_nine", input: "9000009", want: "เก้าล้านเก้าบาทถ้วน"},

		// ── Large mixed numbers ──────────────────────────────────────────
		{name: "complex_number_1", input: "1234567", want: "หนึ่งล้านสองแสนสามหมื่นสี่พันห้าร้อยหกสิบเจ็ดบาทถ้วน"},
		{name: "complex_number_2", input: "7654321", want: "เจ็ดล้านหกแสนห้าหมื่นสี่พันสามร้อยยี่สิบเอ็ดบาทถ้วน"},

		// ── Satang variations ────────────────────────────────────────────
		{name: "point_one", input: "0.1", want: "สิบสตางค์"},
		{name: "point_two", input: "0.2", want: "ยี่สิบสตางค์"},
		{name: "point_nine_nine", input: "0.99", want: "เก้าสิบเก้าสตางค์"},
		{name: "one_point_one", input: "1.1", want: "หนึ่งบาทสิบสตางค์"},
		{name: "two_point_two_five", input: "2.25", want: "สองบาทยี่สิบห้าสตางค์"},
		{name: "ten_point_zero_one", input: "10.01", want: "สิบบาทหนึ่งสตางค์"},

		// ── Satang rounding-like behavior (no rounding expected) ─────────
		{name: "truncate_more_than_two_decimal", input: "1.999", want: "หนึ่งบาทเก้าสิบเก้าสตางค์"},
		{name: "small_fraction", input: "0.004", want: "ศูนย์บาทถ้วน"},

		// ── Negative variations ──────────────────────────────────────────
		{name: "negative_twenty_one", input: "-21", want: "ลบยี่สิบเอ็ดบาทถ้วน"},
		{name: "negative_one_hundred_one", input: "-101", want: "ลบหนึ่งร้อยเอ็ดบาทถ้วน"},
		{name: "negative_large", input: "-1234567", want: "ลบหนึ่งล้านสองแสนสามหมื่นสี่พันห้าร้อยหกสิบเจ็ดบาทถ้วน"},
		{name: "negative_point_nine_nine", input: "-0.99", want: "ลบเก้าสิบเก้าสตางค์"},
		{name: "negative_one_point_zero_one", input: "-1.01", want: "ลบหนึ่งบาทหนึ่งสตางค์"},

		// ── Edge formatting ──────────────────────────────────────────────
		{name: "exact_zero_point_zero", input: "0.00", want: "ศูนย์บาทถ้วน"},
		{name: "one_point_zero_zero", input: "1.00", want: "หนึ่งบาทถ้วน"},
		{name: "large_exact", input: "10000000", want: "สิบล้านบาทถ้วน"},

		// ── Error cases ──────────────────────────────────────────────────
		{name: "too_large", input: "10000000000000", wantErr: true},
		{name: "too_large_negative", input: "-10000000000000", wantErr: true},
		{name: "max_supported", input: "9999999999999.99", want: "เก้าล้านเก้าแสนเก้าหมื่นเก้าพันเก้าร้อยเก้าสิบเก้าล้านเก้าแสนเก้าหมื่นเก้าพันเก้าร้อยเก้าสิบเก้าบาทเก้าสิบเก้าสตางค์"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			val, err := decimal.NewFromString(tc.input)
			if err != nil {
				t.Fatalf("invalid decimal %q: %v", tc.input, err)
			}
			got, err := c.Convert(val)
			if tc.wantErr {
				if err == nil {
					t.Errorf("Convert(%s): expected error, got %q", tc.input, got)
				}
				return
			}
			if err != nil {
				t.Errorf("Convert(%s): unexpected error: %v", tc.input, err)
				return
			}
			if got != tc.want {
				t.Errorf("Convert(%s)\n  got:  %q\n  want: %q", tc.input, got, tc.want)
			}
		})
	}
}
