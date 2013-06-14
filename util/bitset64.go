// Copyright 2013 Adam Peck

package util

import (
	"errors"
)

// Compactly stores 64 bits.
type BitSet64 int64

// Returns the number of bits set to true.
func (b BitSet64) Cardinality() uint {
	var c uint
	for ; b != 0; b &= b - 1 {
		c++
	}
	return c
}

// Sets the bit at index n to true.
// Runtime panic occurs if index n is out of range.
func (b *BitSet64) Set(n uint) {
	if n >= 64 {
		panic(errors.New("BitSet64::Set"))
	}
	*b |= 1 << n
}

// Performs a logical XOR with BitSet64 x.
func (b *BitSet64) XOR(x BitSet64) {
	*b ^= x
}
