package ffi

// #include <stdlib.h>
import "C"
import "unsafe"

type cgoArena struct {
	ptrs     []unsafe.Pointer
	children []freeable
}

func (a *cgoArena) NewArena() *cgoArena {
	arena := new(cgoArena)
	a.children = append(a.children, arena)
	return arena
}

func (a *cgoArena) CString(s string) *C.char {
	cs := C.CString(s)
	a.ptrs = append(a.ptrs, unsafe.Pointer(cs))
	return cs
}

func (a *cgoArena) Free() {
	ptrs := a.ptrs
	a.ptrs = nil
	for _, ptr := range ptrs {
		C.free(ptr)
	}
	children := a.children
	a.children = nil
	for _, child := range children {
		child.Free()
	}
}

func (a *cgoArena) FreeOnFailure(success *bool) {
	if !*success {
		a.Free()
	}
}
