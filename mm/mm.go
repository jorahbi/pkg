package mm

import (
	"unsafe"

	"github.com/heiyeluren/xmm"
)

var menory xmm.XMemory

func init() {
	var err error
	menory, err = new(xmm.Factory).CreateMemory(0.7)
	if err != nil {
		panic("CreateMemory fail ")
	}
}

func MustAlloc[T any]() *T {
	size := unsafe.Sizeof(new(T))
	p, err := menory.Alloc(size)
	if err != nil {
		return new(T)
	}
	return (*T)(p)
}

func Free(obj ...any) {
	if len(obj) == 0 {
		return
	}
	for _, o := range obj {
		menory.Free(uintptr(unsafe.Pointer(&o)))
	}

}
