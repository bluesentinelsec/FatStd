package main

/*
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
*/
import "C"

import (
	"bytes"
	"errors"
	"io"
	"net/url"
	"time"

	"github.com/bluesentinelsec/FatStd/pkg/fatbytes"
	"github.com/bluesentinelsec/FatStd/pkg/net/httpx"
	"net/http"
)

const (
	fatHttpErrCode = 400
)

func fatstdHttpStatusFromError(err error) C.int {
	if err == nil {
		return fatStatusOK
	}
	if errors.Is(err, io.EOF) {
		return fatStatusEOF
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return fatStatusSyntax
	}
	return fatStatusOther
}

type fatHttpClient struct {
	client *httpx.Client
}

type fatHttpResponse struct {
	resp *httpx.Response
}

type fatHttpServer struct {
	server *httpx.Server
}

type fatHttpRequest struct {
	req *httpx.Request
}

func fatstdHttpClientFromHandle(handle uintptr) *fatHttpClient {
	if handle == 0 {
		panic("fatstdHttpClientFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdHttpClientFromHandle: invalid handle")
	}
	c, ok := value.(*fatHttpClient)
	if !ok {
		panic("fatstdHttpClientFromHandle: handle is not http client")
	}
	return c
}

func fatstdHttpResponseFromHandle(handle uintptr) *fatHttpResponse {
	if handle == 0 {
		panic("fatstdHttpResponseFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdHttpResponseFromHandle: invalid handle")
	}
	r, ok := value.(*fatHttpResponse)
	if !ok {
		panic("fatstdHttpResponseFromHandle: handle is not http response")
	}
	return r
}

func fatstdHttpServerFromHandle(handle uintptr) *fatHttpServer {
	if handle == 0 {
		panic("fatstdHttpServerFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdHttpServerFromHandle: invalid handle")
	}
	s, ok := value.(*fatHttpServer)
	if !ok {
		panic("fatstdHttpServerFromHandle: handle is not http server")
	}
	return s
}

func fatstdHttpRequestFromHandle(handle uintptr) *fatHttpRequest {
	if handle == 0 {
		panic("fatstdHttpRequestFromHandle: handle is 0")
	}
	value, ok := fatstdHandles.get(handle)
	if !ok {
		panic("fatstdHttpRequestFromHandle: invalid handle")
	}
	r, ok := value.(*fatHttpRequest)
	if !ok {
		panic("fatstdHttpRequestFromHandle: handle is not http request")
	}
	return r
}

//export fatstd_go_http_client_new
func fatstd_go_http_client_new() C.uintptr_t {
	return C.uintptr_t(fatstdHandles.register(&fatHttpClient{client: httpx.NewClient()}))
}

//export fatstd_go_http_client_free
func fatstd_go_http_client_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_http_client_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_http_client_free: invalid handle")
	}
	if _, ok := value.(*fatHttpClient); !ok {
		panic("fatstd_go_http_client_free: handle is not http client")
	}
}

//export fatstd_go_http_client_get_utf8
func fatstd_go_http_client_get_utf8(handle C.uintptr_t, urlStr *C.char, outResp *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outResp == nil {
		panic("fatstd_go_http_client_get_utf8: outResp is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_http_client_get_utf8: outErr is NULL")
	}
	if urlStr == nil {
		panic("fatstd_go_http_client_get_utf8: url is NULL")
	}

	if _, err := url.ParseRequestURI(C.GoString(urlStr)); err != nil {
		*outResp = 0
		*outErr = C.uintptr_t(fatstdNewError(fatHttpErrCode, err.Error()))
		return fatStatusSyntax
	}

	req, err := http.NewRequest(http.MethodGet, C.GoString(urlStr), nil)
	if err != nil {
		*outResp = 0
		*outErr = C.uintptr_t(fatstdNewError(fatHttpErrCode, err.Error()))
		return fatstdHttpStatusFromError(err)
	}

	client := fatstdHttpClientFromHandle(uintptr(handle))
	resp, err := client.client.Do(req)
	if err != nil {
		*outResp = 0
		*outErr = C.uintptr_t(fatstdNewError(fatHttpErrCode, err.Error()))
		return fatstdHttpStatusFromError(err)
	}

	*outResp = C.uintptr_t(fatstdHandles.register(&fatHttpResponse{resp: resp}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_http_client_post_bytes_utf8
func fatstd_go_http_client_post_bytes_utf8(
	handle C.uintptr_t,
	urlStr *C.char,
	contentType *C.char,
	bodyHandle C.uintptr_t,
	outResp *C.uintptr_t,
	outErr *C.uintptr_t,
) C.int {
	if outResp == nil {
		panic("fatstd_go_http_client_post_bytes_utf8: outResp is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_http_client_post_bytes_utf8: outErr is NULL")
	}
	if urlStr == nil {
		panic("fatstd_go_http_client_post_bytes_utf8: url is NULL")
	}

	if _, err := url.ParseRequestURI(C.GoString(urlStr)); err != nil {
		*outResp = 0
		*outErr = C.uintptr_t(fatstdNewError(fatHttpErrCode, err.Error()))
		return fatStatusSyntax
	}

	b := fatstdBytesFromHandle(uintptr(bodyHandle))
	payload := append([]byte(nil), b.Value()...)
	req, err := http.NewRequest(http.MethodPost, C.GoString(urlStr), bytes.NewReader(payload))
	if err != nil {
		*outResp = 0
		*outErr = C.uintptr_t(fatstdNewError(fatHttpErrCode, err.Error()))
		return fatstdHttpStatusFromError(err)
	}
	if contentType != nil {
		req.Header.Set("Content-Type", C.GoString(contentType))
	}

	client := fatstdHttpClientFromHandle(uintptr(handle))
	resp, err := client.client.Do(req)
	if err != nil {
		*outResp = 0
		*outErr = C.uintptr_t(fatstdNewError(fatHttpErrCode, err.Error()))
		return fatstdHttpStatusFromError(err)
	}

	*outResp = C.uintptr_t(fatstdHandles.register(&fatHttpResponse{resp: resp}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_http_response_status
func fatstd_go_http_response_status(handle C.uintptr_t) C.int {
	resp := fatstdHttpResponseFromHandle(uintptr(handle))
	return C.int(resp.resp.StatusCode)
}

//export fatstd_go_http_response_body
func fatstd_go_http_response_body(handle C.uintptr_t) C.uintptr_t {
	resp := fatstdHttpResponseFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Clone(resp.resp.Body)))
}

//export fatstd_go_http_response_header_get_utf8
func fatstd_go_http_response_header_get_utf8(handle C.uintptr_t, name *C.char) C.uintptr_t {
	if name == nil {
		panic("fatstd_go_http_response_header_get_utf8: name is NULL")
	}
	resp := fatstdHttpResponseFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(resp.resp.Headers.Get(C.GoString(name))))
}

//export fatstd_go_http_response_free
func fatstd_go_http_response_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_http_response_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_http_response_free: invalid handle")
	}
	if _, ok := value.(*fatHttpResponse); !ok {
		panic("fatstd_go_http_response_free: handle is not http response")
	}
}

//export fatstd_go_http_server_new_utf8
func fatstd_go_http_server_new_utf8(addr *C.char, outServer *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outServer == nil {
		panic("fatstd_go_http_server_new_utf8: outServer is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_http_server_new_utf8: outErr is NULL")
	}
	if addr == nil {
		panic("fatstd_go_http_server_new_utf8: addr is NULL")
	}

	server, err := httpx.NewServer(C.GoString(addr))
	if err != nil {
		*outServer = 0
		*outErr = C.uintptr_t(fatstdNewError(fatHttpErrCode, err.Error()))
		return fatstdHttpStatusFromError(err)
	}

	*outServer = C.uintptr_t(fatstdHandles.register(&fatHttpServer{server: server}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_http_server_addr
func fatstd_go_http_server_addr(handle C.uintptr_t) C.uintptr_t {
	server := fatstdHttpServerFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(server.server.AddrString()))
}

//export fatstd_go_http_server_set_static_response
func fatstd_go_http_server_set_static_response(handle C.uintptr_t, status C.int, bodyHandle C.uintptr_t, contentType *C.char) {
	server := fatstdHttpServerFromHandle(uintptr(handle))
	body := fatstdBytesFromHandle(uintptr(bodyHandle))
	ctype := ""
	if contentType != nil {
		ctype = C.GoString(contentType)
	}
	server.server.SetStaticResponse(int(status), body.Value(), ctype)
}

//export fatstd_go_http_server_next_request
func fatstd_go_http_server_next_request(handle C.uintptr_t, timeoutMs C.int64_t, outReq *C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outReq == nil {
		panic("fatstd_go_http_server_next_request: outReq is NULL")
	}
	if outErr == nil {
		panic("fatstd_go_http_server_next_request: outErr is NULL")
	}

	server := fatstdHttpServerFromHandle(uintptr(handle))
	var timeout time.Duration
	if timeoutMs < 0 {
		timeout = -1
	} else {
		timeout = time.Duration(timeoutMs) * time.Millisecond
	}

	request, ok := server.server.NextRequestTimeout(timeout)
	if !ok {
		*outReq = 0
		*outErr = 0
		return fatStatusEOF
	}

	*outReq = C.uintptr_t(fatstdHandles.register(&fatHttpRequest{req: request}))
	*outErr = 0
	return fatStatusOK
}

//export fatstd_go_http_request_method
func fatstd_go_http_request_method(handle C.uintptr_t) C.uintptr_t {
	req := fatstdHttpRequestFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(req.req.Method))
}

//export fatstd_go_http_request_path
func fatstd_go_http_request_path(handle C.uintptr_t) C.uintptr_t {
	req := fatstdHttpRequestFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(req.req.Path))
}

//export fatstd_go_http_request_body
func fatstd_go_http_request_body(handle C.uintptr_t) C.uintptr_t {
	req := fatstdHttpRequestFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdBytesNewFromGoBytes(fatbytes.Clone(req.req.Body)))
}

//export fatstd_go_http_request_header_get_utf8
func fatstd_go_http_request_header_get_utf8(handle C.uintptr_t, name *C.char) C.uintptr_t {
	if name == nil {
		panic("fatstd_go_http_request_header_get_utf8: name is NULL")
	}
	req := fatstdHttpRequestFromHandle(uintptr(handle))
	return C.uintptr_t(fatstdStringNewFromGoString(req.req.Headers.Get(C.GoString(name))))
}

//export fatstd_go_http_request_free
func fatstd_go_http_request_free(handle C.uintptr_t) {
	if handle == 0 {
		panic("fatstd_go_http_request_free: handle is 0")
	}
	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_http_request_free: invalid handle")
	}
	if _, ok := value.(*fatHttpRequest); !ok {
		panic("fatstd_go_http_request_free: handle is not http request")
	}
}

//export fatstd_go_http_server_close
func fatstd_go_http_server_close(handle C.uintptr_t, outErr *C.uintptr_t) C.int {
	if outErr == nil {
		panic("fatstd_go_http_server_close: outErr is NULL")
	}
	if handle == 0 {
		panic("fatstd_go_http_server_close: handle is 0")
	}

	value, ok := fatstdHandles.take(uintptr(handle))
	if !ok {
		panic("fatstd_go_http_server_close: invalid handle")
	}
	s, ok := value.(*fatHttpServer)
	if !ok {
		panic("fatstd_go_http_server_close: handle is not http server")
	}

	if err := s.server.Close(); err != nil {
		*outErr = C.uintptr_t(fatstdNewError(fatHttpErrCode, err.Error()))
		return fatstdHttpStatusFromError(err)
	}

	*outErr = 0
	return fatStatusOK
}
