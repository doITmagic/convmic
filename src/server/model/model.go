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
    From           string
    FromAmount    float64
    To             string
    ToAmount      float64
    ConvertedTime time.Time
}
