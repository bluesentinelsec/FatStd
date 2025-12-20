package main

/*
#include <stdint.h>
*/
import "C"

type fatstdError struct {
	code    int32
	message string
}

func fatstdNewError(code int32, message string) uintptr {
	return fatstdHandles.register(&fatstdError{code: code, message: message})
}

func fatstdErrorFromHandle(handle uintptr) *fatstdError {
	if handle == 0 {
		panic("fatstdErrorFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdErrorFromHandle: invalid handle")
	}
	e, ok := value.(*fatstdError)
	if !ok {
		panic("fatstdErrorFromHandle: handle is not fat error")
	}
	return e
}

//export fatstd_go_error_message
func fatstd_go_error_message(handle C.uintptr_t) C.uintptr_t {
	e := fatstdErrorFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(e.message))
}

//export fatstd_go_error_free
func fatstd_go_error_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_error_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_error_free: invalid handle")
	}
	if _, ok := value.(*fatstdError); !ok {
		panic("fatstd_go_error_free: handle is not fat error")
	}
}

