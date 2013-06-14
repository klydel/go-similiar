// Copyright 2013 Adam Peck

package util

import (
	"reflect"
	"testing"
)

func TestBitSet64_Cardinality(t *testing.T) {
	b := BitSet64(0x6969696969696969)

	if a, e := b.Cardinality(), uint(32); a != e {
		t.Error(a, "!=", e)
	}
}

func TestBitSet64_Set(t *testing.T) {
	var b BitSet64
	b.Set(31)
	b.Set(32)
	b.Set(32)
	b.Set(31)

	if a, e := b, BitSet64(0x180000000); a != e {
		t.Error(a, "!=", e)
	}
}

func TestBitSet64_Set_Panic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("Expected panic")
		}
		if _, ok := r.(error); !ok {
			t.Error(reflect.TypeOf(r).Elem(), "!=", reflect.TypeOf((*error)(nil)).Elem());
		}
	}()

	var b BitSet64
	b.Set(64)
}

func TestBitSet64_XOR(t *testing.T) {
	b := BitSet64(0x0555555555555550)
	b.XOR(BitSet64(0x0AAAAAAAAAAAAAA0))

	if a, e := b, BitSet64(0x0FFFFFFFFFFFFFF0); a != e {
		t.Error(a, "!=", e)
	}
}
