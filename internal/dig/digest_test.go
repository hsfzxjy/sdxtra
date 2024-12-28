package dig_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/hsfzxjy/sdxtra/internal/dig"
	"github.com/hsfzxjy/sdxtra/internal/dig/helper"
)

func SumStream(args ...any) []byte {
	h := helper.H{Hash: sha256.New()}
	for _, arg := range args {
		val := reflect.ValueOf(arg)
		if h.TryWritePrimitive(val) {
			continue
		}
		switch arg := arg.(type) {
		case []any:
			h.Write(SumStream(arg...))
		case []byte:
			h.Write(arg)
		}
	}
	return h.Sum(nil)
}

func TestDigest(t *testing.T) {
	type Foo struct {
		B string
		A int `dig:"AA"`
		C any `dig:"-"`
	}
	type Case struct {
		values []any
		stream []any
	}
	cases := []Case{
		{[]any{
			Foo{A: 1, B: "2"},
			&Foo{A: 1, B: "2"},
		}, []any{
			dig.KindStruct,
			"AA",
			[]any{dig.KindNumber,
				1},
			"B",
			[]any{dig.KindString,
				"2"}}},
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {

			target := SumStream(c.stream...)
			for j, v := range c.values {
				t.Run(strconv.Itoa(j), func(t *testing.T) {
					got, err := dig.Sum(context.Background(), v)
					if err != nil {
						t.Fatal(err)
					}
					if !bytes.Equal(got, target) {
						t.Errorf("got %x, want %x", got, target)
					}
				})
			}
		})
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type FilePath string

func (p FilePath) DigestTo(w dig.Writer) error {
	return w.CopyFromFile(string(p))
}

func (p FilePath) DigestCoreKind() dig.CoreKind {
	return dig.KindString
}

func TestDigestFile(t *testing.T) {
	tmpdir := t.TempDir()
	var content [1024]byte
	rand.Read(content[:])
	filepath := filepath.Join(tmpdir, "test")
	check(os.WriteFile(filepath, content[:], 0644))
	target := SumStream(dig.KindString, []any{content[:]})
	got,err := dig.Sum(context.Background(), FilePath(filepath))
	check(err)
	if !bytes.Equal(got, target) {
		t.Errorf("got %x, want %x", got, target)
	}
}
