package helper_test

import (
	"bytes"
	"crypto/sha256"
	"math"
	"reflect"
	"strconv"
	"testing"

	"github.com/hsfzxjy/sdxtra/internal/dig/helper"
)

func sum(x any) []byte {
	h := helper.H{Hash: sha256.New()}
	val := reflect.ValueOf(x)
	if h.TryWritePrimitive(val) {
		return h.SumReset()
	}
	panic("unreachable")
}

func TestHashNumber(t *testing.T) {
	type Float32 float32
	type Float64 float64

	type Case []any
	const _1p53 = 1 << 53
	const _1p63 = 1 << 63
	const _1_5p63 = 1<<63 + 1<<62
	_1p64_1 := uint64(math.Nextafter(float64(1<<64-1), 0))

	cases := []Case{
		{int(0), int8(0), int16(0), int32(0), int64(0), uint(0), uint8(0), uint16(0), uint32(0), uint64(0), uintptr(0), Float32(0), Float64(0), -Float32(0), -Float64(0)},

		{math.NaN(), float32(math.NaN()), float64(math.NaN())},

		{int64(_1p53), uint64(_1p53),
			Float32(_1p53), Float64(_1p53),
			Float32(_1p53 + 1), Float64(_1p53 + 1)},

		{int64(-_1p53), Float32(-_1p53), Float64(-_1p53)},

		{uint64(_1p63), Float32(_1p63), Float64(_1p63)},

		{int64(-_1p63), Float32(-_1p63), Float64(-_1p63)},

		{uint64(_1_5p63), Float32(_1_5p63), Float64(_1_5p63)},

		{uint64(_1p64_1), Float64(_1p64_1)},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			prev := sum(c[0])
			for j, x := range c[1:] {
				t.Run(strconv.Itoa(j), func(t *testing.T) {
					cur := sum(x)
					if !bytes.Equal(prev, cur) {
						t.Fatalf("expect %#V:%x, got %#V:%x", c[0], prev, x, cur)
					}
				})
			}
		})
	}
}
