// Copyright 2010 Felipe Alves Cavani. All rights reserved.
// Start date:		2013-05-12
// Last modification:	2013-

package rand

import (
	"testing"

	"github.com/fcavani/e"
)

const tests = 100000

func TestInt63n(t *testing.T) {
	for i := 0; i < tests; i++ {
		n, err := Int63n(10, "go")
		if err != nil {
			t.Fatal("Int63n go failed:", err)
		}
		if n >= 10 || n < 0 {
			t.Fatal("Int63n go failed", n)
		}
	}
	for i := 0; i < tests; i++ {
		n, err := Int63n(10, "go-crypto")
		if err != nil {
			t.Fatal("Int63n go-crypto failed:", err)
		}
		if n >= 10 || n < 0 {
			t.Fatal("Int63n go-crypto failed", n)
		}
	}
	for i := 0; i < tests; i++ {
		n, err := Int63n(10, "/dev/urandom")
		if err != nil {
			t.Fatal("Int63n /dev/urandom failed:", err)
		}
		if n >= 10 || n < 0 {
			t.Fatal("Int63n /dev/urandom failed", n)
		}
	}
}

func TestBytes(t *testing.T) {
	b, err := Bytes(10, "go")
	if err != nil {
		t.Fatal("Bytes go failed:", err)
	}
	if len(b) != 10 {
		t.Fatal("Bytes go failed")
	}

	b, err = Bytes(10, "go-crypto")
	if err != nil {
		t.Fatal("Bytes go-crypto failed:", err)
	}
	if len(b) != 10 {
		t.Fatal("Bytes go-crypto failed")
	}

	b, err = Bytes(10, "/dev/urandom")
	if err != nil {
		t.Fatal("Bytes /dev/urandom failed:", err)
	}
	if len(b) != 10 {
		t.Fatal("Bytes /dev/urandom failed")
	}
}

func TestrandDevInt64(t *testing.T) {
	_, err := randDevInt64("go")
	if err != nil {
		t.Fatal("randDevInt64 go failed:", err)
	}

	_, err = randDevInt64("go-crypto")
	if err != nil {
		t.Fatal("randDevInt64 go-crypto failed:", err)
	}

	_, err = randDevInt64("/dev/urandom")
	if err != nil {
		t.Fatal("randDevInt64 /dev/urandom failed:", err)
	}
}

func chkChars(str string, chars []string) bool {
F:
	for s := range str {
		for chr := range chars {
			if chr == s {
				continue F
			}
		}
		return false
	}
	return true
}

func TestChars(t *testing.T) {
	s, err := Chars(10, Letters, "go")
	if err != nil {
		t.Fatal("Chars go failed:", err)
	}
	if err != nil {
		t.Fatal("Chars go failed:", err)
	}
	if len(s) != 10 {
		t.Fatal("Chars go failed")
	}
	if chkChars(s, Letters) == false {
		t.Fatal("Chars go failed")
	}

	s, err = Chars(10, Letters, "go-crypto")
	if err != nil {
		t.Fatal("Chars go-crypto failed:", err)
	}
	if err != nil {
		t.Fatal("Chars go-crypto failed:", err)
	}
	if len(s) != 10 {
		t.Fatal("Chars go-crypto failed")
	}
	if chkChars(s, Letters) == false {
		t.Fatal("Chars go-crypto failed")
	}

	s, err = Chars(10, Letters, "/dev/urandom")
	if err != nil {
		t.Fatal("Chars /dev/urandom failed:", err)
	}
	if err != nil {
		t.Fatal("Chars /dev/urandom failed:", err)
	}
	if len(s) != 10 {
		t.Fatal("Chars /dev/urandom failed")
	}
	if chkChars(s, Letters) == false {
		t.Fatal("Chars /dev/urandom failed")
	}
}

func TestString(t *testing.T) {
	s, err := String(10, "go")
	if err != nil {
		t.Fatal("String go failed:", err)
	}
	if len(s) != 10 {
		t.Fatal("String go failed", len(s), s)
	}

	s, err = String(10, "go-crypto")
	if err != nil {
		t.Fatal("String go-crypto failed:", err)
	}
	if len(s) != 10 {
		t.Fatal("String go-crypto failed", len(s), s)
	}

	s, err = String(10, "/dev/urandom")
	if err != nil {
		t.Fatal("String /dev/urandom failed:", err)
	}
	if len(s) != 10 {
		t.Fatal("String /dev/urandom failed", len(s), s)
	}
}

func TestUuid(t *testing.T) {
	_, err := Uuid()
	if err != nil {
		t.Fatal("Uuid go failed:", err)
	}
}

var testsPerm = 1000

func TestPermutations(t *testing.T) {
	for i := 0; i < testsPerm; i++ {
		l, err := Int63n(10000, "go")
		if err != nil {
			t.Fatal("Int63n failed:", e.Trace(e.Forward(err)))
		}
		l++
		a, err := permutation(int(l), "go")
		if err != nil {
			t.Fatal("permutation failed:", e.Trace(e.Forward(err)))
		}
		for _, val := range a {
			if int64(val) >= l {
				t.Fatal("permutation failed", l, val)
			}
		}
	}
}

type Vector []int

func (v Vector) Len() int {
	return len(v)
}

func (v Vector) At(i int) interface{} {
	return v[i]
}

func (v Vector) Set(i int, val interface{}) {
	v[i] = val.(int)
}

func TestRandomPermutation(t *testing.T) {
	var in Vector = Vector{1, 2, 3, 4, 5, 6, 7}
	for i := 0; i < 1000000; i++ {
		out := make(Vector, len(in))
		err := RandomPermutation(in, out, "go-crypto")
		if err != nil {
			t.Fatal("RandomPermutation failed:", e.Trace(e.Forward(err)))
		}
		for _, valIn := range in {
			count := 0
			for _, valOut := range out {
				if valOut == valIn {
					count++
				}
				if count > 1 {
					t.Fatal("Duplicado!!!")
				}
			}
		}
	}
}

func TestStringPermutation(t *testing.T) {
	in := StringPermutation(LowerCase)
	for i := 0; i < 1000000; i++ {
		out := make(StringPermutation, len(in))
		err := RandomPermutation(in, out, "go-crypto")
		if err != nil {
			t.Fatal("RandomPermutation failed:", e.Trace(e.Forward(err)))
		}
		for _, valIn := range in {
			count := 0
			for _, valOut := range out {
				if valOut == valIn {
					count++
				}
				if count > 1 {
					t.Fatal("Duplicado!!!")
				}
			}
		}
	}
}
