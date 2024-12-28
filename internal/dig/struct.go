package dig

import (
	"cmp"
	"errors"
	"reflect"
	"slices"
)

func hashStruct(w Writer, val reflect.Value) error {
	type fieldMeta struct {
		v    reflect.Value
		name string
	}
	fields := make([]fieldMeta, 0, val.NumField())
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		ftyp := typ.Field(i)
		name := ftyp.Tag.Get("dig")
		if name == "" {
			name = ftyp.Name
		} else if name == "-" {
			continue
		}
		f := val.Field(i)
		fields = append(fields, fieldMeta{v: f, name: name})
	}
	slices.SortFunc(fields, func(a, b fieldMeta) int {
		return cmp.Compare(a.name, b.name)
	})
	subWriter := w.new()
	for _, f := range fields {
		err := sum(subWriter, f.v)
		if err != nil {
			return errors.New("hash field " + f.name + ": " + err.Error())
		}
		w.WriteString(f.name)
		w.Write(subWriter.SumReset())
	}
	return nil
}
