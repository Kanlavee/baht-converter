# Thai Baht Text Converter

A production-ready Go service that converts numeric `decimal.Decimal` values into formal Thai Baht text representation — suitable for cheques, invoices, and payment systems.

## Design Decisions

- Supports negative numbers (prefix with "ลบ")
- Fractional values are truncated to 2 digits (no rounding)
- Thai grammar rules applied:
  - "ยี่" for tens
  - "เอ็ด" for last digit

## Features

- 🛡️ **Error handling** — returns `ErrAmountTooLarge` for values exceeding 9,999,999,999,999.99
- 🎯 **Precision-safe** — uses [`shopspring/decimal`](https://github.com/shopspring/decimal), no floats
- 📜 **Correct Thai grammar** — `เอ็ด`, `ยี่สิบ`, `สิบ`, and all positional words
- ➖ **Negative number support** — prefix `"ลบ"` on negative amounts
- 💰 **Full currency support** — satang precision to 2 decimal places
- 🔢 **Arbitrarily large numbers** — recursive ล้าน grouping handles billions and beyond
- 🧩 **Struct-based API** — `Converter` struct designed for service/DI integration
- ✅ **79 unit tests** at ~97% code coverage

## Project Structure

```
baht-converter/
├── cmd/
│   └── main.go                    # Entry point with example conversions
├── internal/
│   └── converter/
│       ├── converter.go            # Converter struct + Thai Baht logic
│       └── converter_test.go       # Table-driven unit tests (79 cases)
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

- Go **1.25.7** or later

## Installation

```bash
git clone https://github.com/Kanlavee/baht-converter.git
cd baht-converter
go mod download
```

## How to Run

```bash
go run cmd/main.go
```

**Expected output:**

```
Input: 1234            -> Output: หนึ่งพันสองร้อยสามสิบสี่บาทถ้วน
Input: 33333.75        -> Output: สามหมื่นสามพันสามร้อยสามสิบสามบาทเจ็ดสิบห้าสตางค์
Input: 0               -> Output: ศูนย์บาทถ้วน
Input: 11.01           -> Output: สิบเอ็ดบาทหนึ่งสตางค์
Input: 21.5            -> Output: ยี่สิบเอ็ดบาทห้าสิบสตางค์
Input: 1000001         -> Output: หนึ่งล้านเอ็ดบาทถ้วน
Input: -51             -> Output: ลบห้าสิบเอ็ดบาทถ้วน
Input: -0.5            -> Output: ลบห้าสิบสตางค์
Input: 0.0099999999    -> Output: ศูนย์บาทถ้วน
Input: 1.999           -> Output: หนึ่งบาทเก้าสิบเก้าสตางค์
Input: 9999999.99      -> Output: เก้าล้านเก้าแสนเก้าหมื่นเก้าพันเก้าร้อยเก้าสิบเก้าบาทเก้าสิบเก้าสตางค์
```

## How to Run Unit Tests

```bash
# Run all tests
go test ./internal/converter/

# Verbose output (see each test case name)
go test -v ./internal/converter/

```

## How to Use the Package

```go
import (
    "github.com/shopspring/decimal"
    "github.com/Kanlavee/baht-converter/internal/converter"
)

c := converter.NewConverter()

result, err := c.Convert(decimal.NewFromFloat(1234))
// result = "หนึ่งพันสองร้อยสามสิบสี่บาทถ้วน", err = nil

result, err = c.Convert(decimal.NewFromFloat(33333.75))
// result = "สามหมื่นสามพันสามร้อยสามสิบสามบาทเจ็ดสิบห้าสตางค์", err = nil

result, err = c.Convert(decimal.NewFromFloat(-51))
// result = "ลบห้าสิบเอ็ดบาทถ้วน", err = nil

// Amount exceeds maximum:
result, err = c.Convert(decimal.NewFromFloat(10_000_000_000_000))
// result = "", err = converter.ErrAmountTooLarge
```

## Error Handling

`Convert` returns `(string, error)`. The only error currently returned is:

| Error | Condition |
|-------|-----------|
| `converter.ErrAmountTooLarge` | Absolute value > `9,999,999,999,999.99` |

**Maximum supported value:** `9,999,999,999,999.99` (~10 trillion baht).  
This cap prevents silent `int64` overflow inside the internal integer conversion and covers all realistic financial amounts.

```go
result, err := c.Convert(amount)
if err != nil {
    // handle: log, return HTTP 422, etc.
}
```

## Decimal Handling (Truncation)

Satang is derived by **truncating** (not rounding) the fractional part to 2 decimal places. Any digits beyond the second decimal place are silently discarded:

| Input | Satang | Output |
|-------|--------|--------|
| `1.999` | 99 (truncated from 99.9) | หนึ่งบาท**เก้าสิบเก้า**สตางค์ |
| `0.009` | 0 (truncated from 0.9) | **ศูนย์บาทถ้วน** |
| `0.006` | 0 (truncated from 0.6) | **ศูนย์บาทถ้วน** |
| `0.004` | 0 (truncated from 0.4) | **ศูนย์บาทถ้วน** |

> Use `decimal.NewFromString("1.99")` for exact values rather than `decimal.NewFromFloat(1.999)`.

| Rule | Example Input | Output |
|------|--------------|--------|
| Ones digit `1` in a compound number → **เอ็ด** | 11, 21, 101, 1,000,001 | สิบ**เอ็ด**, ยี่สิบ**เอ็ด**, ร้อย**เอ็ด**, ล้าน**เอ็ด** |
| Standalone `1` → **หนึ่ง** | 1, 0.01 satang | **หนึ่ง**บาท, **หนึ่ง**สตางค์ |
| Tens digit `2` → **ยี่สิบ** | 20, 21, 25 | **ยี่สิบ**, **ยี่สิบ**เอ็ด, **ยี่สิบ**ห้า |
| Tens digit `1` → **สิบ** (no prefix) | 10–19 | **สิบ**, **สิบ**เอ็ด, **สิบ**เก้า |
| No fractional part → **บาทถ้วน** | 100, 1234 | หนึ่งร้อย**บาทถ้วน** |
| With satang → **สตางค์** | 1.50, 33333.75 | ...บาท**ห้าสิบสตางค์** |

## Negative Number Handling

Negative amounts are prefixed with `"ลบ"` applied to the absolute-value result.

| Input | Absolute Text | Final Output |
|-------|--------------|--------------|
| `-51` | ห้าสิบเอ็ดบาทถ้วน | **ลบ**ห้าสิบเอ็ดบาทถ้วน |
| `-0.5` | ห้าสิบสตางค์ | **ลบ**ห้าสิบสตางค์ |
| `-33333.75` | สามหมื่น...สตางค์ | **ลบ**สามหมื่น...สตางค์ |

> **Note:** Satang-only amounts (zero baht, non-zero satang) omit the "บาท" portion entirely.
> e.g., `0.25` → `"ยี่สิบห้าสตางค์"` (not `"ศูนย์บาทยี่สิบห้าสตางค์"`)

## Conversion Reference

| Input | Thai Text |
|-------|-----------|
| `0` | ศูนย์บาทถ้วน |
| `1` | หนึ่งบาทถ้วน |
| `11` | สิบเอ็ดบาทถ้วน |
| `21` | ยี่สิบเอ็ดบาทถ้วน |
| `101` | หนึ่งร้อยเอ็ดบาทถ้วน |
| `1,000` | หนึ่งพันบาทถ้วน |
| `1,234` | หนึ่งพันสองร้อยสามสิบสี่บาทถ้วน |
| `33,333.75` | สามหมื่นสามพันสามร้อยสามสิบสามบาทเจ็ดสิบห้าสตางค์ |
| `1,000,001` | หนึ่งล้านเอ็ดบาทถ้วน |
| `1,000,000,000` | หนึ่งพันล้านบาทถ้วน |
| `0.25` | ยี่สิบห้าสตางค์ |
| `-51` | ลบห้าสิบเอ็ดบาทถ้วน |
| `-0.5` | ลบห้าสิบสตางค์ |



