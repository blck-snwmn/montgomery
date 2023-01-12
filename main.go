package main

import (
	"fmt"
	"math/big"
)

type montgomery struct {
	n    *big.Int
	n_   *big.Int
	r    *big.Int
	mask *big.Int
	r2   *big.Int
}

func NewMontgomery() *montgomery {
	r := new(big.Int).Lsh(big.NewInt(1), 130)
	mask := new(big.Int).Sub(r, big.NewInt(1))
	n := new(big.Int).Sub(r, big.NewInt(5))
	r2 := new(big.Int).Mul(r, r)
	r2 = r2.Mod(r2, n)

	var (
		result = big.NewInt(0)
		t      = big.NewInt(0)
		rr     = new(big.Int).Set(r)
		i      = big.NewInt(1)
	)

	for rr.Cmp(big.NewInt(1)) > 0 {
		if t.Bit(0) == 0 {
			t.Add(t, n)
			result.Add(result, i)
		}
		t.Rsh(t, 1)
		rr.Rsh(rr, 1)
		i.Lsh(i, 1)
	}

	return &montgomery{
		n:    n,
		n_:   result,
		r:    r,
		mask: mask,
		r2:   r2,
	}
}

func (m *montgomery) mul(a, b *big.Int) *big.Int {
	tmpa := new(big.Int).Set(a)
	tmpb := new(big.Int).Set(b)
	aa := m.mr(tmpa.Mul(tmpa, m.r2))
	bb := m.mr(tmpb.Mul(tmpb, m.r2))
	cc := m.mr(tmpa.Mul(aa, bb))
	return m.mr(cc)
}

func (m *montgomery) mul_(a, b *big.Int) *big.Int {
	tmpa := new(big.Int).Set(a)
	tmpb := new(big.Int).Set(b)
	aa := m.mr_(tmpa.Mul(tmpa, m.r2))
	bb := m.mr_(tmpb.Mul(tmpb, m.r2))
	cc := m.mr_(tmpa.Mul(aa, bb))
	return m.mr_(cc)
}

// mr do montgomery reduction
func (m *montgomery) mr(t *big.Int) *big.Int {
	tmp := new(big.Int)
	tmp = tmp.
		Mul(t, m.n_).
		And(tmp, m.mask).
		Mul(tmp, m.n).
		Add(tmp, t).
		Rsh(tmp, uint(m.r.BitLen()-1))
	if tmp.Cmp(m.n) > 0 {
		tmp.Sub(tmp, m.n)
	}
	return tmp
}

// mr do montgomery reduction
func (m *montgomery) mr_(t *big.Int) *big.Int {
	tmp := new(big.Int)
	tmp = tmp.
		Mul(t, m.n_).
		Mod(tmp, m.r). //todo 怪しい
		Mul(tmp, m.n).
		Add(tmp, t).
		Div(tmp, m.r) //todo 怪しい
	if tmp.Cmp(m.n) > 0 {
		tmp.Sub(tmp, m.n)
	}
	return tmp
}

func main() {
	var (
		a = big.NewInt(41221321312312)
		b = big.NewInt(89798798234)
	)
	m := NewMontgomery()
	{
		c := m.mul_(a, b)
		fmt.Println("mul_", c)
	}
	{
		c := m.mul(a, b)
		fmt.Println("mul ", c)
	}
	{
		t := new(big.Int)
		c := t.Mul(a, b).Mod(t, m.n)
		fmt.Println("calc", c)
	}
	{
		aa := m.mr(m.mr(new(big.Int).Mul(a, m.r2)))
		fmt.Println(aa)
		fmt.Println(a)
		fmt.Println(m.n)
		fmt.Println(aa.String() == a.String())
	}
}
