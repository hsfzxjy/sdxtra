package dig

import (
	"context"
	"errors"
	"reflect"
)

type Digester interface {
	DigestTo(w Writer) error
	DigestCoreKind() CoreKind
}

var digesterContentType = reflect.TypeFor[Digester]()

func sum(w Writer, val reflect.Value) error {
	for val.Kind() == reflect.Ptr &&
		!val.Type().Implements(digesterContentType) {
		val = val.Elem()
	}
	if val.Type().Implements(digesterContentType) {
		iface := val.Interface().(Digester)
		coreKind := iface.DigestCoreKind()
		w.WriteUint64(uint64(coreKind))
		return iface.DigestTo(w)
	}

	rk := val.Kind()
	coreKind := reflectKindToCoreKind[rk]
	if coreKind == KindInvalid {
		return errors.New("unsupported kind: " + rk.String())
	}
	w.WriteUint64(uint64(coreKind))
	if coreKind != KindStruct {
		if !w.TryWritePrimitive(val) {
			panic("unreachable: " + rk.String() + " should have been handled by TryWritePrimitive")
		}
		return nil
	} else {
		return hashStruct(w, val)
	}
}

func Sum(ctx context.Context, val any) ([]byte, error) {
	w := NewWriter(ctx)
	err := sum(w, reflect.ValueOf(val))
	if err != nil {
		return nil, err
	}
	return w.SumReset(), nil
}
