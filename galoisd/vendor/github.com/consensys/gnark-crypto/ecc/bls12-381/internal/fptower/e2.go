// Copyright 2020 ConsenSys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fptower

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fp"
	"math/big"
)

// E2 is a degree two finite field extension of fp.Element
type E2 struct {
	A0, A1 fp.Element
}

// Equal returns true if z equals x, false otherwise
func (z *E2) Equal(x *E2) bool {
	return z.A0.Equal(&x.A0) && z.A1.Equal(&x.A1)
}

// Bits
// TODO @gbotrel fixme this shouldn't return a E2
func (z *E2) Bits() E2 {
	r := E2{}
	r.A0 = z.A0.Bits()
	r.A1 = z.A1.Bits()
	return r
}

// Cmp compares (lexicographic order) z and x and returns:
//
//	-1 if z <  x
//	 0 if z == x
//	+1 if z >  x
func (z *E2) Cmp(x *E2) int {
	if a1 := z.A1.Cmp(&x.A1); a1 != 0 {
		return a1
	}
	return z.A0.Cmp(&x.A0)
}

// LexicographicallyLargest returns true if this element is strictly lexicographically
// larger than its negation, false otherwise
func (z *E2) LexicographicallyLargest() bool {
	// adapted from github.com/zkcrypto/bls12_381
	if z.A1.IsZero() {
		return z.A0.LexicographicallyLargest()
	}
	return z.A1.LexicographicallyLargest()
}

// SetString sets a E2 element from strings
func (z *E2) SetString(s1, s2 string) *E2 {
	z.A0.SetString(s1)
	z.A1.SetString(s2)
	return z
}

// SetZero sets an E2 elmt to zero
func (z *E2) SetZero() *E2 {
	z.A0.SetZero()
	z.A1.SetZero()
	return z
}

// Set sets an E2 from x
func (z *E2) Set(x *E2) *E2 {
	z.A0 = x.A0
	z.A1 = x.A1
	return z
}

// SetOne sets z to 1 in Montgomery form and returns z
func (z *E2) SetOne() *E2 {
	z.A0.SetOne()
	z.A1.SetZero()
	return z
}

// SetRandom sets a0 and a1 to random values
func (z *E2) SetRandom() (*E2, error) {
	if _, err := z.A0.SetRandom(); err != nil {
		return nil, err
	}
	if _, err := z.A1.SetRandom(); err != nil {
		return nil, err
	}
	return z, nil
}

// IsZero returns true if the two elements are equal, false otherwise
func (z *E2) IsZero() bool {
	return z.A0.IsZero() && z.A1.IsZero()
}

func (z *E2) IsOne() bool {
	return z.A0.IsOne() && z.A1.IsZero()
}

// Add adds two elements of E2
func (z *E2) Add(x, y *E2) *E2 {
	addE2(z, x, y)
	return z
}

// Sub two elements of E2
func (z *E2) Sub(x, y *E2) *E2 {
	subE2(z, x, y)
	return z
}

// Double doubles an E2 element
func (z *E2) Double(x *E2) *E2 {
	doubleE2(z, x)
	return z
}

// Neg negates an E2 element
func (z *E2) Neg(x *E2) *E2 {
	negE2(z, x)
	return z
}

// String implements Stringer interface for fancy printing
func (z *E2) String() string {
	return z.A0.String() + "+" + z.A1.String() + "*u"
}

// MulByElement multiplies an element in E2 by an element in fp
func (z *E2) MulByElement(x *E2, y *fp.Element) *E2 {
	var yCopy fp.Element
	yCopy.Set(y)
	z.A0.Mul(&x.A0, &yCopy)
	z.A1.Mul(&x.A1, &yCopy)
	return z
}

// Conjugate conjugates an element in E2
func (z *E2) Conjugate(x *E2) *E2 {
	z.A0 = x.A0
	z.A1.Neg(&x.A1)
	return z
}

// Halve sets z = z / 2
func (z *E2) Halve() {
	z.A0.Halve()
	z.A1.Halve()
}

// Legendre returns the Legendre symbol of z
func (z *E2) Legendre() int {
	var n fp.Element
	z.norm(&n)
	return n.Legendre()
}

// Exp sets z=xᵏ (mod q²) and returns it
func (z *E2) Exp(x E2, k *big.Int) *E2 {
	if k.IsUint64() && k.Uint64() == 0 {
		return z.SetOne()
	}

	e := k
	if k.Sign() == -1 {
		// negative k, we invert
		// if k < 0: xᵏ (mod q²) == (x⁻¹)ᵏ (mod q²)
		x.Inverse(&x)

		// we negate k in a temp big.Int since
		// Int.Bit(_) of k and -k is different
		e = bigIntPool.Get().(*big.Int)
		defer bigIntPool.Put(e)
		e.Neg(k)
	}

	z.SetOne()
	b := e.Bytes()
	for i := 0; i < len(b); i++ {
		w := b[i]
		for j := 0; j < 8; j++ {
			z.Square(z)
			if (w & (0b10000000 >> j)) != 0 {
				z.Mul(z, &x)
			}
		}
	}

	return z
}

func init() {
	q := fp.Modulus()
	tmp := big.NewInt(3)
	sqrtExp1.Set(q).Sub(&sqrtExp1, tmp).Rsh(&sqrtExp1, 2)

	tmp.SetUint64(1)
	sqrtExp2.Set(q).Sub(&sqrtExp2, tmp).Rsh(&sqrtExp2, 1)
}

var sqrtExp1, sqrtExp2 big.Int

// Sqrt sets z to the square root of and returns z
// The function does not test wether the square root
// exists or not, it's up to the caller to call
// Legendre beforehand.
// cf https://eprint.iacr.org/2012/685.pdf (algo 9)
func (z *E2) Sqrt(x *E2) *E2 {

	var a1, alpha, b, x0, minusone E2

	minusone.SetOne().Neg(&minusone)

	a1.Exp(*x, &sqrtExp1)
	alpha.Square(&a1).
		Mul(&alpha, x)
	x0.Mul(x, &a1)
	if alpha.Equal(&minusone) {
		var c fp.Element
		c.Set(&x0.A0)
		z.A0.Neg(&x0.A1)
		z.A1.Set(&c)
		return z
	}
	a1.SetOne()
	b.Add(&a1, &alpha)

	b.Exp(b, &sqrtExp2).Mul(&x0, &b)
	z.Set(&b)
	return z
}

// BatchInvertE2 returns a new slice with every element inverted.
// Uses Montgomery batch inversion trick
//
// if a[i] == 0, returns result[i] = a[i]
func BatchInvertE2(a []E2) []E2 {
	res := make([]E2, len(a))
	if len(a) == 0 {
		return res
	}

	zeroes := make([]bool, len(a))
	var accumulator E2
	accumulator.SetOne()

	for i := 0; i < len(a); i++ {
		if a[i].IsZero() {
			zeroes[i] = true
			continue
		}
		res[i].Set(&accumulator)
		accumulator.Mul(&accumulator, &a[i])
	}

	accumulator.Inverse(&accumulator)

	for i := len(a) - 1; i >= 0; i-- {
		if zeroes[i] {
			continue
		}
		res[i].Mul(&res[i], &accumulator)
		accumulator.Mul(&accumulator, &a[i])
	}

	return res
}

func (z *E2) Select(cond int, caseZ *E2, caseNz *E2) *E2 {
	//Might be able to save a nanosecond or two by an aggregate implementation

	z.A0.Select(cond, &caseZ.A0, &caseNz.A0)
	z.A1.Select(cond, &caseZ.A1, &caseNz.A1)

	return z
}

func (z *E2) Div(x *E2, y *E2) *E2 {
	var r E2
	r.Inverse(y).Mul(x, &r)
	return z.Set(&r)
}