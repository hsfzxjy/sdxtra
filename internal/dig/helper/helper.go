package helper

import (
	"encoding/binary"
	"hash"
	"math"
	"reflect"
)

type Helper interface {
	WriteUint64(x uint64)
	WriteInt64(x int64)
	WriteFloat64(x float64)
	WriteString(x string)
	WriteBool(x bool)
	TryWritePrimitive(x reflect.Value) (ok bool)
	SumReset() []byte
}

type H struct {
	hash.Hash
	buffer [32]byte
}

// Encode int64, uint64, float64 with a universal format.
// sign (1 byte): 00 for 0, 01 for positive, 10 for negative, 11 for Inf, 100 for -Inf, 101 for NaN
// integeral (8 byte): if x < 2^64-1, encode x directly; otherwise, encode int(x) as float64
// fractional (8 byte): if x is integer, encode 0; otherwise, encode x - int(x)

const (
	zero     = 0
	positive = 1
	negative = 2
	inf      = 3
	neginf   = 4
	nan      = 5
)

func (h *H) SumReset() []byte {
	hsh := h.Sum(h.buffer[:0])
	h.Reset()
	return hsh
}

func (h *H) WriteUint64(x uint64) {
	buf := h.buffer[:17]
	clear(buf)
	if x > 0 {
		buf[0] = positive
		binary.LittleEndian.PutUint64(buf[1:], x)
	}
	h.Write(buf)
}

func (h *H) WriteInt64(x int64) {
	buf := h.buffer[:17]
	clear(buf)
	if x > 0 {
		buf[0] = positive
		binary.LittleEndian.PutUint64(buf[1:], uint64(x))
	} else if x < 0 {
		buf[0] = negative
		binary.LittleEndian.PutUint64(buf[1:], uint64(-x))
	}
	h.Write(buf)
}

func (h *H) WriteFloat64(x float64) {
	buf := h.buffer[:17]
	clear(buf)
	switch {
	case x == 0:
		buf[0] = zero
	case math.IsInf(x, 1):
		buf[0] = inf
	case math.IsInf(x, -1):
		buf[0] = neginf
	case math.IsNaN(x):
		buf[0] = nan
	default:
		goto NORMAL
	}
	h.Write(buf)
	return
NORMAL:
	if x < 0 {
		buf[0] = negative
		x = -x
	} else {
		buf[0] = positive
	}
	i, f := math.Modf(x)
	if i < 1<<64 {
		binary.LittleEndian.PutUint64(buf[1:], uint64(i))
	} else {
		binary.LittleEndian.PutUint64(buf[1:], math.Float64bits(x))
	}
	binary.LittleEndian.PutUint64(buf[9:], math.Float64bits(f))
	h.Write(buf)
}

func (h *H) WriteString(x string) {
	h.Write([]byte(x))
}

func (h *H) WriteBool(x bool) {
	buf := h.buffer[:1]
	if x {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	h.Write(buf)
}

func (h *H) TryWritePrimitive(x reflect.Value) (ok bool) {
	switch x.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		h.WriteUint64(x.Uint())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		h.WriteInt64(x.Int())
	case reflect.Float32, reflect.Float64:
		h.WriteFloat64(x.Float())
	case reflect.String:
		h.WriteString(x.String())
	case reflect.Bool:
		h.WriteBool(x.Bool())
	default:
		return false
	}
	return true
}
