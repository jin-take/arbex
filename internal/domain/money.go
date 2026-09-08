package domain

import (
	"errors"
	"math/big"
)

var (
	ErrNegativeAmount = errors.New("amount must be non-negative")
	ErrInvalidDecimals = errors.New("decimals out of range")
)

type Amount struct {
	units    *big.Int
	decimals uint8
}

func NewAmount(units *big.Int, decimals uint8) (Amount, error) {
	if units == nil {
		units = new(big.Int)
	}
	if units.Sign() < 0 {
		return Amount{}, ErrNegativeAmount
	}
	if decimals > 36 {
		return Amount{}, ErrInvalidDecimals
	}
	return Amount{units: new(big.Int).Set(units), decimals: decimals}, nil
}

func MustAmount(units string, decimals uint8) Amount {
	v, ok := new(big.Int).SetString(units, 10)
	if !ok {
		panic("invalid integer amount")
	}
	a, err := NewAmount(v, decimals)
	if err != nil {
		panic(err)
	}
	return a
}

func (a Amount) Units() *big.Int {
	if a.units == nil {
		return new(big.Int)
	}
	return new(big.Int).Set(a.units)
}

func (a Amount) Decimals() uint8 { return a.decimals }
func (a Amount) IsZero() bool    { return a.Units().Sign() == 0 }

func (a Amount) Add(b Amount) (Amount, error) {
	if a.decimals != b.decimals {
		return Amount{}, ErrInvalidDecimals
	}
	return NewAmount(new(big.Int).Add(a.Units(), b.Units()), a.decimals)
}

func (a Amount) Sub(b Amount) (Amount, error) {
	if a.decimals != b.decimals {
		return Amount{}, ErrInvalidDecimals
	}
	v := new(big.Int).Sub(a.Units(), b.Units())
	if v.Sign() < 0 {
		return Amount{}, ErrNegativeAmount
	}
	return NewAmount(v, a.decimals)
}
