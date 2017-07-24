package util

import (
	"reflect"
	"unsafe"
)

func HashBytesInt(data []byte) int {
	var v int = 2166136261
	for _, c := range data {
		v = (v ^ int(c)) * 16777619
	}
	return v
}

func HashIntInt(i int) int {
	ptr := (*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&i)),
		Len:  4,
		Cap:  4,
	}))
	return HashBytesInt(*ptr)
}
