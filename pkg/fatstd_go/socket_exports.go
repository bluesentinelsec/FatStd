package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"errors"
	"io"
	"math"
	"net"
	"unsafe"

	fatsocket "github.com/bluesentinelsec/FatStd/pkg/net/socket"
)

const (
	fatSocketErrCode = 300
)

func fatstdSocketStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusOK
	}
	if errors.Is(err, io.EOF) {
		return fatStatusEOF
	}
	var addrErr *net.AddrError
	if errors.As(err, &addrErr) {
		return fatStatusSyntax
	}
	return fatStatusOther
}

func fatstdTcpListenerFromHandle(handle uintptr) *fatsocket.TCPListener {
	if handle == 0 {
		panic("fatstdTcpListenerFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTcpListenerFromHandle: invalid handle")
	}
	l, ok := value.(*fatsocket.TCPListener)
	if !ok {
		panic("fatstdTcpListenerFromHandle: handle is not tcp listener")
	}
	return l
}

func fatstdTcpConnFromHandle(handle uintptr) *fatsocket.TCPConn {
	if handle == 0 {
		panic("fatstdTcpConnFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdTcpConnFromHandle: invalid handle")
	}
	c, ok := value.(*fatsocket.TCPConn)
	if !ok {
		panic("fatstdTcpConnFromHandle: handle is not tcp conn")
	}
	return c
}

func fatstdUdpConnFromHandle(handle uintptr) *fatsocket.UDPConn {
	if handle == 0 {
		panic("fatstdUdpConnFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdUdpConnFromHandle: invalid handle")
	}
	c, ok := value.(*fatsocket.UDPConn)
	if !ok {
		panic("fatstdUdpConnFromHandle: handle is not udp conn")
	}
	return c
}

func fatstdGoBytesFromCPtr(ptr *C.char, length C.size_t, label string) []byte {
	if ptr == nil {
		if length == 0 {
			return []byte{}
		}
		panic("fatstd: " + label + " is NULL but len > 0")
	}
	if length > C.size_t(math.MaxInt) {
		panic("fatstd: " + label + " too large")
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(ptr)), int(length))
}

//export fatstd_go_tcp_dial_utf8
func fatstd_go_tcp_dial_utf8(addr *C.char, outConn *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outConn == nil {
		panic("fatstd_go_tcp_dial_utf8: outConn is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tcp_dial_utf8: outErr is NULL")
	}
	if addr == nil {
		panic("fatstd_go_tcp_dial_utf8: addr is NULL")
	}

	conn, err := fatsocket.DialTCP(C.GoString(addr))
	if err != nil {
		*outConn = 0
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outConn = C.uintptr_t(fatstdHandles.register(conn))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tcp_listener_listen_utf8
func fatstd_go_tcp_listener_listen_utf8(addr *C.char, outListener *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outListener == nil {
		panic("fatstd_go_tcp_listener_listen_utf8: outListener is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tcp_listener_listen_utf8: outErr is NULL")
	}
	if addr == nil {
		panic("fatstd_go_tcp_listener_listen_utf8: addr is NULL")
	}

	listener, err := fatsocket.ListenTCP(C.GoString(addr))
	if err != nil {
		*outListener = 0
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outListener = C.uintptr_t(fatstdHandles.register(listener))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tcp_listener_accept
func fatstd_go_tcp_listener_accept(handle C.uintptr_t, outConn *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outConn == nil {
		panic("fatstd_go_tcp_listener_accept: outConn is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tcp_listener_accept: outErr is NULL")
	}

	listener := fatstdTcpListenerFromHandle(uintptr(handle))
	conn, err := listener.Accept()
	if err != nil {
		*outConn = 0
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outConn = C.uintptr_t(fatstdHandles.register(conn))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tcp_listener_addr
func fatstd_go_tcp_listener_addr(handle C.uintptr_t) C.uintptr_t {
	listener := fatstdTcpListenerFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(listener.AddrString()))
}

//export fatstd_go_tcp_listener_close
func fatstd_go_tcp_listener_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_tcp_listener_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_tcp_listener_close: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tcp_listener_close: invalid handle")
	}
	listener, ok := value.(*fatsocket.TCPListener)
	if !ok {
		panic("fatstd_go_tcp_listener_close: handle is not tcp listener")
	}

	if err := listener.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tcp_conn_read
func fatstd_go_tcp_conn_read(handle C.uintptr_t, dst *C.char, dstLen C.size_t, outN *C.size_t, outEOF *C.bool, outErr *C.uintptr_t) C.int {
	if outN == nil {
		panic("fatstd_go_tcp_conn_read: outN is NULL")
	}
	if outEOF == nil {
		panic("fatstd_go_tcp_conn_read: outEOF is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tcp_conn_read: outErr is NULL")
	}

	buf := fatstdGoBytesFromCPtr(dst, dstLen, "dst")
	conn := fatstdTcpConnFromHandle(uintptr(handle))
	n, err := conn.Read(buf)
	*outN = C.size_t(n)

	if err != nil {
		if errors.Is(err, io.EOF) {
			*outEOF = C.bool(true)
			*outErr = 0
			return fatStatusEOF
		}
		*outEOF = C.bool(false)
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outEOF = C.bool(false)
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tcp_conn_write
func fatstd_go_tcp_conn_write(handle C.uintptr_t, src *C.char, srcLen C.size_t, outN *C.size_t, outErr *C.uintptr_t) C.int {
	if outN == nil {
		panic("fatstd_go_tcp_conn_write: outN is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_tcp_conn_write: outErr is NULL")
	}

	buf := fatstdGoBytesFromCPtr(src, srcLen, "src")
	conn := fatstdTcpConnFromHandle(uintptr(handle))
	n, err := conn.Write(buf)
	*outN = C.size_t(n)
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_tcp_conn_local_addr
func fatstd_go_tcp_conn_local_addr(handle C.uintptr_t) C.uintptr_t {
	conn := fatstdTcpConnFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(conn.LocalAddrString()))
}

//export fatstd_go_tcp_conn_remote_addr
func fatstd_go_tcp_conn_remote_addr(handle C.uintptr_t) C.uintptr_t {
	conn := fatstdTcpConnFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(conn.RemoteAddrString()))
}

//export fatstd_go_tcp_conn_close
func fatstd_go_tcp_conn_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_tcp_conn_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_tcp_conn_close: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_tcp_conn_close: invalid handle")
	}
	conn, ok := value.(*fatsocket.TCPConn)
	if !ok {
		panic("fatstd_go_tcp_conn_close: handle is not tcp conn")
	}

	if err := conn.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_udp_listen_utf8
func fatstd_go_udp_listen_utf8(addr *C.char, outConn *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outConn == nil {
		panic("fatstd_go_udp_listen_utf8: outConn is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_udp_listen_utf8: outErr is NULL")
	}
	if addr == nil {
		panic("fatstd_go_udp_listen_utf8: addr is NULL")
	}

	conn, err := fatsocket.ListenUDP(C.GoString(addr))
	if err != nil {
		*outConn = 0
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outConn = C.uintptr_t(fatstdHandles.register(conn))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_udp_dial_utf8
func fatstd_go_udp_dial_utf8(addr *C.char, outConn *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outConn == nil {
		panic("fatstd_go_udp_dial_utf8: outConn is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_udp_dial_utf8: outErr is NULL")
	}
	if addr == nil {
		panic("fatstd_go_udp_dial_utf8: addr is NULL")
	}

	conn, err := fatsocket.DialUDP(C.GoString(addr))
	if err != nil {
		*outConn = 0
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outConn = C.uintptr_t(fatstdHandles.register(conn))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_udp_conn_read_from
func fatstd_go_udp_conn_read_from(handle C.uintptr_t, dst *C.char, dstLen C.size_t, outN *C.size_t, outAddr *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outN == nil {
		panic("fatstd_go_udp_conn_read_from: outN is NULL")
	}
	if outAddr == nil {
		panic("fatstd_go_udp_conn_read_from: outAddr is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_udp_conn_read_from: outErr is NULL")
	}

	buf := fatstdGoBytesFromCPtr(dst, dstLen, "dst")
	conn := fatstdUdpConnFromHandle(uintptr(handle))
	n, addr, err := conn.ReadFrom(buf)
	*outN = C.size_t(n)

	if err != nil {
		*outAddr = 0
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	if addr != nil {
		*outAddr = C.uintptr_t(fatstdStringNewFromGoString(addr.String()))
	} else {
		*outAddr = 0
	}
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_udp_conn_write_to_utf8
func fatstd_go_udp_conn_write_to_utf8(handle C.uintptr_t, src *C.char, srcLen C.size_t, addr *C.char, outN *C.size_t, outErr *C.uintptr_t) C.int {
	if outN == nil {
		panic("fatstd_go_udp_conn_write_to_utf8: outN is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_udp_conn_write_to_utf8: outErr is NULL")
	}
	if addr == nil {
		panic("fatstd_go_udp_conn_write_to_utf8: addr is NULL")
	}

	buf := fatstdGoBytesFromCPtr(src, srcLen, "src")
	raddr, err := net.ResolveUDPAddr("udp", C.GoString(addr))
	if err != nil {
		*outN = 0
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	conn := fatstdUdpConnFromHandle(uintptr(handle))
	n, err := conn.WriteTo(buf, raddr)
	*outN = C.size_t(n)
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_udp_conn_write
func fatstd_go_udp_conn_write(handle C.uintptr_t, src *C.char, srcLen C.size_t, outN *C.size_t, outErr *C.uintptr_t) C.int {
	if outN == nil {
		panic("fatstd_go_udp_conn_write: outN is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_udp_conn_write: outErr is NULL")
	}

	buf := fatstdGoBytesFromCPtr(src, srcLen, "src")
	conn := fatstdUdpConnFromHandle(uintptr(handle))
	n, err := conn.Write(buf)
	*outN = C.size_t(n)
	if err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_udp_conn_local_addr
func fatstd_go_udp_conn_local_addr(handle C.uintptr_t) C.uintptr_t {
	conn := fatstdUdpConnFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(conn.LocalAddrString()))
}

//export fatstd_go_udp_conn_remote_addr
func fatstd_go_udp_conn_remote_addr(handle C.uintptr_t) C.uintptr_t {
	conn := fatstdUdpConnFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(conn.RemoteAddrString()))
}

//export fatstd_go_udp_conn_close
func fatstd_go_udp_conn_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_udp_conn_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_udp_conn_close: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_udp_conn_close: invalid handle")
	}
	conn, ok := value.(*fatsocket.UDPConn)
	if !ok {
		panic("fatstd_go_udp_conn_close: handle is not udp conn")
	}

	if err := conn.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatSocketErrCode, err.Error()))
		return fatstdSocketStatusFromError(err)
	}

	*outErr = 0
	return fatStatusOK
}
