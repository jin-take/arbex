package domain

import (
	"math/big"
	"testing"
)

func TestAmountCopiesInputAndOutput(t *testing.T) {
	in := big.NewInt(1_000_000)
	a, err := NewAmount(in, 6)
	if err != nil {
		t.Fatal(err)
	}
	in.SetInt64(0)
	if got := a.Units().String(); got != "1000000" {
		t.Fatalf("input mutation leaked into amount: %s", got)
	}
	out := a.Units()
	out.SetInt64(0)
	if got := a.Units().String(); got != "1000000" {
		t.Fatalf("output mutation leaked into amount: %s", got)
	}
}

func TestAmountDecimals(t *testing.T) {
	for _, decimals := range []uint8{6, 18} {
		a, err := NewAmount(big.NewInt(1), decimals)
		if err != nil {
			t.Fatalf("decimals=%d: %v", decimals, err)
		}
		if a.Decimals() != decimals {
			t.Fatalf("want decimals=%d got=%d", decimals, a.Decimals())
		}
	}
}

func TestAmountRejectsNegative(t *testing.T) {
	if _, err := NewAmount(big.NewInt(-1), 18); err == nil {
		t.Fatal("expected negative amount error")
	}
}

func TestAmountArithmetic(t *testing.T) {
	a := MustAmount("100", 6)
	b := MustAmount("40", 6)
	c, err := a.Add(b)
	if err != nil || c.Units().String() != "140" {
		t.Fatalf("add: %v %s", err, c.Units())
	}
	d, err := a.Sub(b)
	if err != nil || d.Units().String() != "60" {
		t.Fatalf("sub: %v %s", err, d.Units())
	}
	if _, err := b.Sub(a); err == nil {
		t.Fatal("expected underflow to be rejected")
	}
}

func TestAmountLargeUint256Range(t *testing.T) {
	max := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	a, err := NewAmount(max, 18)
	if err != nil {
		t.Fatal(err)
	}
	if a.Units().Cmp(max) != 0 {
		t.Fatal("uint256-sized amount changed")
	}
}
