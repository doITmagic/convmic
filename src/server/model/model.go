package model

import "time"

type Currency struct {
    Name  string
    Value string
}

type CurrencyConvert struct {
    Name   string
    Amount float64
}

type CurrencyConverted struct {
    from           string
    from_amount    float64
    to             string
    to_amount      float64
    converted_time time.Time
}
