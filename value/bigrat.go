// Copyright 2014 Rob Pike. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package value // import "robpike.io/ivy/value"

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type BigRat struct {
	*big.Rat
}

func setBigRatString(s string) (br BigRat, err error) {
	base := conf.InputBase()
	r := big.NewRat(0, 1)
	var ok bool
	slash := strings.IndexByte(s, '/')
	if slash < 0 {
		r, ok = r.SetString(s)
	} else {
		switch base {
		case 0, 10: // Fine as is.
			r, ok = r.SetString(s)
		default:
			// big.Rat doesn't handle arbitrary bases, but big.Int does,
			// so do the numerator and denominator separately.
			var num, denom BigInt
			num, err = setBigIntString(s[:slash])
			if err == nil {
				denom, err = setBigIntString(s[slash+1:])
			}
			if err != nil {
				return
			}
			return BigRat{r.SetFrac(num.Int, denom.Int)}, nil
		}
	}
	if !ok {
		return BigRat{}, errors.New("rational number syntax")
	}
	return BigRat{r}, nil
}

func (r BigRat) String() string {
	format := conf.Format()
	if format != "" {
		return fmt.Sprintf(conf.RatFormat(), r.Num(), r.Denom())
	}
	num := BigInt{r.Num()}
	den := BigInt{r.Denom()}
	return fmt.Sprintf("%s/%s", num, den)
}

func (r BigRat) Eval() Value {
	return r
}

func (r BigRat) toType(which valueType) Value {
	switch which {
	case intType:
		panic("big rat to int")
	case bigIntType:
		panic("big rat to big int")
	case bigRatType:
		return r
	case vectorType:
		return NewVector([]Value{r})
	case matrixType:
		return newMatrix([]Value{one, one}, []Value{r})
	}
	panic("BigRat.toType")
}

// shrink pulls, if possible, a BigRat down to a BigInt or Int.
func (r BigRat) shrink() Value {
	if !r.IsInt() {
		return r
	}
	return BigInt{r.Num()}.shrink()
}
