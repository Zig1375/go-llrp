package binutil

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

var (
	//hexRunes to hold hex chars
	hexRunes = []rune("abcdef0123456789")
)

// GenerateNLengthHexString returns random hex rune for n length
func GenerateNLengthHexString(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := range b {
		b[i] = hexRunes[rand.Intn(len(hexRunes))]
	}
	return string(b)
}

// GenerateNLengthRandomBinRuneSlice returns n-length random binary string
// max == 0 for no cap limit
func GenerateNLengthRandomBinRuneSlice(n int, max uint) ([]rune, uint) {
	binstr := make([]rune, n)
	sum := uint(0)
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < n; i++ {
		var b rune
		if max != uint(0) && max < uint(math.Pow(float64(2), float64(n-i))) {
			b = '0'
		} else if rand.Intn(2) == 0 {
			b = '0'
		} else {
			b = '1'
		}
		binstr[i] = b
		if b == '1' {
			sum += uint(math.Pow(float64(2), float64(n-i-1)))
		}
	}

	if max != uint(0) && max < sum {
		binstr, sum = GenerateNLengthRandomBinRuneSlice(n, max)
	}

	return binstr, sum
}

// GenerateNLengthZeroPaddingRuneSlice returns n-length zero padding string
func GenerateNLengthZeroPaddingRuneSlice(n int) []rune {
	binstr := make([]rune, n)

	for i := 0; i < n; i++ {
		binstr[i] = '0'
	}

	return binstr
}

// GenerateRandomInt return random int value with min-max
func GenerateRandomInt(min, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(max-min) + min
}

// Pack the data into (partial) LLRP packet payload.
func Pack(data []interface{}) []byte {
	buf := new(bytes.Buffer)
	for _, v := range data {
		binary.Write(buf, binary.BigEndian, v)
	}
	return buf.Bytes()
}

// Parse6BinRuneSliceToRune translate 6 rune slices into a 6-bit encoded rune
func Parse6BinRuneSliceToRune(r []rune) (rune, error) {
	if len(r) != 6 {
		return 0, errors.New("Given rune should be exactley len=6")
	}
	i, err := strconv.ParseInt(string(r), 2, 64)
	if err != nil {
		return 0, err
	}
	if (48&i)>>4 < 2 {
		i |= 64
	}
	return rune(i), nil
}

// ParseBinRuneSliceToUint8Slice returns uint8 slice from binary string
// Precondition: len(bs) % 8 == 0
func ParseBinRuneSliceToUint8Slice(bs []rune) ([]uint8, error) {
	if len(bs)%8 != 0 {
		return nil, errors.New("non-8 bit length binary string passed to ParseBinRuneSliceToUint8Slice")
	} else if len(bs) < 8 {
		return nil, errors.New("binary string length less than 8 given to ParseBinRuneSliceToUint8Slice")
	}

	bsSize := len(bs) / 8
	uints := make([]uint8, bsSize)

	for j := 0; j < bsSize; j++ {
		uintRep := uint8(0)
		for i := 0; i < 8; i++ {
			if bs[j*8-i+7] == '1' {
				uintRep += uint8(math.Pow(float64(2), float64(i)))
			}
		}
		uints[j] = uintRep
	}

	return uints, nil
}

// ParseDecimalStringToBinRuneSlice convert serial to binary rune slice
func ParseDecimalStringToBinRuneSlice(s string) []rune {
	n, _ := strconv.ParseInt(s, 10, 64)
	return []rune(fmt.Sprintf("%b", big.NewInt(n)))
}

// ParseHexStringToBinString converts hex string to binary string
func ParseHexStringToBinString(s string) (string, error) {
	re := regexp.MustCompile("[^0-9a-fA-F]")
	if re.FindStringIndex(s) != nil {
		return "", errors.New("Input to ParseHexStringToBinString is not a hex string!")
	}

	var bs string
	for _, c := range s {
		n, _ := strconv.ParseInt(string(c), 16, 32)
		bs = fmt.Sprintf("%s%.4b", bs, n)
	}
	return bs, nil
}

// ParseRuneTo6BinRuneSlice coverts rune into 6-bit encoding rune slice
func ParseRuneTo6BinRuneSlice(r rune) []rune {
	if r >= 64 { // if the rune is after '@' in ASCII table
		r -= 64
	}
	binString := fmt.Sprintf("%.6b", r)
	return []rune(binString)
}
