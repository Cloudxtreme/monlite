// Copyright 2012 Felipe Alves Cavani. All rights reserved.
// Start date:      2012-04-26
// Last modification:   2012-

// Rand package have auxiliary function to generate random numbers and random strings.
// It's give you the choice of use the random and pseudo-random number
// generator of go or you can use it direct (not portable way).
package rand

import (
	crypto "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/rand"
	"os"
	"time"

	"github.com/fcavani/e"
)

const ErrInvalidLength = "invalid length"

var Number []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
var Letters []string = []string{"a", "b", "c", "d", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "x", "w", "y", "z", "A", "B", "C", "D", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "W", "Y", "Z"}
var LowerCase []string = Letters[:25]
var UpperCase []string = Letters[25:]
var NumberLetters []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "x", "w", "y", "z", "A", "B", "C", "D", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "W", "Y", "Z"}
var NumberLettersSimbols []string = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "x", "w", "y", "z", "A", "B", "C", "D", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "X", "W", "Y", "Z", "!", "@", "#", "$", "%", "&", "*", "_", "-", "+", "=", "|", "\\", "/", "?", ":", ";", "<", ">", ",", "."}

//MaxInt64 = 1<<63 - 1 = 2^64 - 1
var MaxInt64 = 1<<63 - 1

func init() {
	var timeSeed int64
	i, err := Int64R()
	if err != nil {
		//println(err.Error())
		timeSeed = time.Now().Unix()
	} else {
		timeSeed = time.Now().Unix() ^ i
	}
	rand.Seed(timeSeed)
}

// Int63n returns, as an int64, a non-negative random number in [0,n).
// Where dev is the path to the random number device or dev is go for the go
// pseudo random number generator.
func Int63n(n int64, dev string) (int64, error) {
	if n <= 0 {
		return 0, e.New("invalid interval")
	}
	switch dev {
	case "go":
		return rand.Int63n(n), nil
	case "go-crypto":
		var val int64
		err := binary.Read(crypto.Reader, binary.LittleEndian, &val)
		if err != nil {
			return 0, e.Forward(err)
		}
		//return int64(math.Abs(float64(val))) % n, nil
		return int64(math.Abs(float64(val % n))), nil
	default:
		f, err := os.Open(dev)
		if err != nil {
			return 0, e.Forward(err)
		}
		defer f.Close()
		var val int64
		err = binary.Read(f, binary.LittleEndian, &val)
		if err != nil {
			return 0, e.Forward(err)
		}
		return int64(math.Abs(float64(val % n))), nil
	}
	panic("don't get here")
}

// Bytes generate random bytes
func Bytes(l int, dev string) ([]byte, error) {
	if l <= 0 {
		return nil, e.New(ErrInvalidLength)
	}
	switch dev {
	case "go":
		b := make([]byte, l)
		for i := 0; i < l; i++ {
			b[i] = byte(rand.Intn(255))
		}
		return b, nil
	case "go-crypto":
		b := make([]byte, l)
		n, err := io.ReadFull(crypto.Reader, b)
		if n != len(b) || err != nil {
			return nil, e.Forward(err)
		}
		return b, nil
	default:
		f, err := os.Open(dev)
		if err != nil {
			return nil, e.Forward(err)
		}
		b := make([]byte, l)
		f.Read(b)
		f.Close()
		return b, nil
	}
	panic("don't get here")
}

// Int64U generate a int64 using /dev/urandom.
func Int64U() (int64, error) {
	i, err := randDevInt64("/dev/urandom")
	return i, e.Forward(err)
}

// Int64R generate a int64 using /dev/random.
func Int64R() (int64, error) {
	i, err := randDevInt64("/dev/random")
	return i, e.Forward(err)
}

// Int64Go generate a int64 using go random number generator
func Int64Go(l int) (int64, error) {
	i, err := randDevInt64("go")
	return i, e.Forward(err)
}

// Int64Crypto generate a int64 using go crypto random number generator
func Int64Crypto(l int) (int64, error) {
	i, err := randDevInt64("go-crypto")
	return i, e.Forward(err)
}

func randDevInt64(dev string) (int64, error) {
	b, err := Bytes(64, dev)
	if err != nil {
		return 0, e.Forward(err)
	}
	i := big.NewInt(0)
	i.SetBytes(b)
	return i.Int64(), nil
}

// Chars generate a random string with selected characters.
// l is the number of characters. Where chars is array with selected characters
// and dev is the path to the random number device or dev is go for the go
// pseudo random number generator or go-crypto for the random number generator.
func Chars(l uint64, chars []string, dev string) (string, error) {
	if l <= 0 {
		return "", e.New(ErrInvalidLength)
	}
	str := ""
	charslen := int64(len(chars))
	for i := uint64(0); i < l; i++ {
		i, err := Int63n(charslen, dev)
		if err != nil {
			return "", e.Forward(err)
		}
		str += chars[i]
	}
	//fmt.Println("seed:", seed, "rand:", str)
	return str, nil
}

// StringU is like String but uses /dev/urandom.
// Where l is the number of characters.
func StringU(l int) (string, error) {
	s, err := String(l, "/dev/urandom")
	return s, e.Forward(err)
}

// StringR is like String but uses /dev/random.
// Where l is the number of characters.
func StringR(l int) (string, error) {
	s, err := String(l, "/dev/random")
	return s, e.Forward(err)
}

// Stringgo is like String but uses go random number generator.
// Where l is the number of characters.
func StringGo(l int) (string, error) {
	s, err := String(l, "go")
	return s, e.Forward(err)
}

// StringCrypto is like String but uses go crypto random number generator.
// Where l is the number of characters.
func StringCrypto(l int) (string, error) {
	s, err := String(l, "go-crypto")
	return s, e.Forward(err)
}

// String generate a string of random characters.
// Where l is the number of characters and dev is the path
// to the random number device or dev is go for the go
// pseudo random number generator.
func String(chars int, dev string) (string, error) {
	b, err := Bytes(chars, dev)
	if err != nil {
		return "", e.Forward(err)
	}
	return string(b), nil
}

// Uuid generate the "Universally unique identifier"
func Uuid() (string, error) {
	b, err := Bytes(17, "go-crypto")
	if err != nil {
		return "", e.Forward(err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	// 0x80 0x90 0xA0 0xB0
	var b8 byte
	switch {
	case b[16] < 0x40:
		b8 = 0x80
	case b[16] >= 0x40 && b[16] < 0x80:
		b8 = 0x90
	case b[16] >= 0x80 && b[16] < 0xC0:
		b8 = 0xA0
	case b[16] >= 0xC0 && b[16] < 0xFF:
		b8 = 0xB0
	}
	b[8] = (b[8] & 0x0f) | b8
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

type LenAtSetter interface {
	// Len return the length of slice.
	Len() int
	// At return the element in i index.
	At(i int) interface{}
	// Set sets the i element with val.
	Set(i int, val interface{})
}

// StringPermutation is an adaptor for permutate a string slice.
type StringPermutation []string

func (s StringPermutation) Len() int {
	return len(s)
}

func (s StringPermutation) At(i int) interface{} {
	return s[i]
}

func (s StringPermutation) Set(i int, val interface{}) {
	str := val.(string)
	s[i] = str
}

func permutation(length int, dev string) ([]int64, error) {
	a := make([]int64, length)
	m := make(map[int64]bool, length)
	i := 0
	count := 0
	for {
		r, err := Int63n(int64(length), dev)
		if err != nil {
			return nil, e.Forward(err)
		}
		if _, ok := m[r]; ok {
			if count > 100*length {
				return nil, e.New("can't generate the vector with permutations")
			}
			count++
			continue
		}
		count = 0
		m[r] = true
		a[i] = r
		i++
		if i >= length {
			break
		}
	}
	return a, nil
}

// RandomPermutation scramble the in array and put the result in the
// out array. in and outs arrays must implement the
// LenAtSetter interface. dev is the random number generator device
// or go for the go pseudo-random ou go-crypto for the go random number
// generator.
func RandomPermutation(in, out LenAtSetter, dev string) error {
	if in.Len() > out.Len() {
		return e.New("in length is great than out length")
	}
	perm, err := permutation(in.Len(), dev)
	if err != nil {
		return e.Forward(err)
	}
	for i, p := range perm {
		val := in.At(int(p))
		out.Set(i, val)
	}
	return nil
}

// FileName return a filename with letters number of letter, prefix and ext extension.
func FileName(prefix, ext string, letters int) (string, error) {
	name, err := Chars(uint64(letters), NumberLetters, "go-crypto")
	if err != nil {
		return "", e.Forward(err)
	}
	if ext != "" {
		return prefix + name + "." + ext, nil
	}
	return prefix + name, nil
}
