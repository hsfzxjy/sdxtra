package dig

import "reflect"

type CoreKind uint64

const (
	KindInvalid CoreKind = iota
	KindBool
	KindNumber
	KindString
	KindStruct
)

var reflectKindToCoreKind = [reflect.UnsafePointer + 3]CoreKind{
	reflect.Bool:    KindBool,
	reflect.Int:     KindNumber,
	reflect.Int8:    KindNumber,
	reflect.Int16:   KindNumber,
	reflect.Int32:   KindNumber,
	reflect.Int64:   KindNumber,
	reflect.Uint:    KindNumber,
	reflect.Uint8:   KindNumber,
	reflect.Uint16:  KindNumber,
	reflect.Uint32:  KindNumber,
	reflect.Uint64:  KindNumber,
	reflect.Uintptr: KindNumber,
	reflect.Float32: KindNumber,
	reflect.Float64: KindNumber,
	reflect.String:  KindString,
	reflect.Struct:  KindStruct,
}
