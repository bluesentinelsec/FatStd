from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestHttp(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_handle = fat_string_handle_type()
        fat_error = fat_string_handle_type()
        fat_bytes = fat_string_handle_type()

        cls.fat_BytesNewN = bind(
            "fat_BytesNewN", argtypes=[ctypes.c_void_p, ctypes.c_size_t], restype=fat_bytes
        )
        cls.fat_BytesLen = bind("fat_BytesLen", argtypes=[fat_bytes], restype=ctypes.c_size_t)
        cls.fat_BytesCopyOut = bind(
            "fat_BytesCopyOut", argtypes=[fat_bytes, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_BytesFree = bind("fat_BytesFree", argtypes=[fat_bytes], restype=None)

        cls.fat_StringLenBytes = bind(
            "fat_StringLenBytes", argtypes=[fat_handle], restype=ctypes.c_size_t
        )
        cls.fat_StringCopyOut = bind(
            "fat_StringCopyOut", argtypes=[fat_handle, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_handle], restype=None)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_HttpClientNew = bind("fat_HttpClientNew", argtypes=[], restype=fat_handle)
        cls.fat_HttpClientFree = bind("fat_HttpClientFree", argtypes=[fat_handle], restype=None)
        cls.fat_HttpClientGetUTF8 = bind(
            "fat_HttpClientGetUTF8",
            argtypes=[fat_handle, ctypes.c_char_p, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_HttpClientPostBytesUTF8 = bind(
            "fat_HttpClientPostBytesUTF8",
            argtypes=[
                fat_handle,
                ctypes.c_char_p,
                ctypes.c_char_p,
                fat_bytes,
                ctypes.POINTER(fat_handle),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )

        cls.fat_HttpResponseStatus = bind(
            "fat_HttpResponseStatus", argtypes=[fat_handle], restype=ctypes.c_int
        )
        cls.fat_HttpResponseBody = bind(
            "fat_HttpResponseBody", argtypes=[fat_handle], restype=fat_bytes
        )
        cls.fat_HttpResponseHeaderGetUTF8 = bind(
            "fat_HttpResponseHeaderGetUTF8", argtypes=[fat_handle, ctypes.c_char_p], restype=fat_handle
        )
        cls.fat_HttpResponseFree = bind("fat_HttpResponseFree", argtypes=[fat_handle], restype=None)

        cls.fat_HttpServerNewUTF8 = bind(
            "fat_HttpServerNewUTF8",
            argtypes=[ctypes.c_char_p, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_HttpServerAddr = bind("fat_HttpServerAddr", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_HttpServerSetStaticResponse = bind(
            "fat_HttpServerSetStaticResponse",
            argtypes=[fat_handle, ctypes.c_int, fat_bytes, ctypes.c_char_p],
            restype=None,
        )
        cls.fat_HttpServerNextRequest = bind(
            "fat_HttpServerNextRequest",
            argtypes=[fat_handle, ctypes.c_int64, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_HttpRequestMethod = bind("fat_HttpRequestMethod", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_HttpRequestPath = bind("fat_HttpRequestPath", argtypes=[fat_handle], restype=fat_handle)
        cls.fat_HttpRequestBody = bind("fat_HttpRequestBody", argtypes=[fat_handle], restype=fat_bytes)
        cls.fat_HttpRequestHeaderGetUTF8 = bind(
            "fat_HttpRequestHeaderGetUTF8", argtypes=[fat_handle, ctypes.c_char_p], restype=fat_handle
        )
        cls.fat_HttpRequestFree = bind("fat_HttpRequestFree", argtypes=[fat_handle], restype=None)
        cls.fat_HttpServerClose = bind(
            "fat_HttpServerClose", argtypes=[fat_handle, ctypes.POINTER(fat_error)], restype=ctypes.c_int
        )

        cls.FAT_OK = 0

    def _bytes_new(self, data: bytes) -> int:
        raw = ctypes.create_string_buffer(data, len(data))
        h = self.fat_BytesNewN(ctypes.addressof(raw), len(raw.raw))
        self.assertNotEqual(0, h)
        return h

    def _bytes_to_py(self, h: int) -> bytes:
        n = self.fat_BytesLen(h)
        if n == 0:
            dst = ctypes.create_string_buffer(1, 1)
            copied = self.fat_BytesCopyOut(h, ctypes.addressof(dst), 0)
            self.assertEqual(0, copied)
            return b""
        dst = ctypes.create_string_buffer(n, n)
        copied = self.fat_BytesCopyOut(h, ctypes.addressof(dst), len(dst.raw))
        self.assertEqual(n, copied)
        return dst.raw

    def _string_to_py(self, h: int) -> str:
        n = self.fat_StringLenBytes(h)
        if n == 0:
            self.fat_StringFree(h)
            return ""
        dst = ctypes.create_string_buffer(n, n)
        copied = self.fat_StringCopyOut(h, ctypes.addressof(dst), len(dst.raw))
        self.assertEqual(n, copied)
        s = dst.raw.decode("utf-8")
        self.fat_StringFree(h)
        return s

    def test_http_get_and_post(self) -> None:
        server = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = self.fat_HttpServerNewUTF8(b"127.0.0.1:0", ctypes.byref(server), ctypes.byref(err))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, server.value)

        addr = self._string_to_py(self.fat_HttpServerAddr(server.value))
        base_url = f"http://{addr}"

        resp_body = self._bytes_new(b"hello")
        self.fat_HttpServerSetStaticResponse(server.value, 200, resp_body, b"text/plain")
        self.fat_BytesFree(resp_body)

        client = self.fat_HttpClientNew()
        self.assertNotEqual(0, client)

        resp = fat_string_handle_type()()
        err_resp = fat_string_handle_type()()
        status = self.fat_HttpClientGetUTF8(
            client, f"{base_url}/hello".encode("utf-8"), ctypes.byref(resp), ctypes.byref(err_resp)
        )
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err_resp.value)
        self.assertNotEqual(0, resp.value)
        self.assertEqual(200, self.fat_HttpResponseStatus(resp.value))

        body_handle = self.fat_HttpResponseBody(resp.value)
        self.assertEqual(b"hello", self._bytes_to_py(body_handle))
        self.fat_BytesFree(body_handle)

        header = self._string_to_py(self.fat_HttpResponseHeaderGetUTF8(resp.value, b"Content-Type"))
        self.assertEqual("text/plain", header)

        self.fat_HttpResponseFree(resp.value)

        req = fat_string_handle_type()()
        err_req = fat_string_handle_type()()
        status = self.fat_HttpServerNextRequest(server.value, 2000, ctypes.byref(req), ctypes.byref(err_req))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err_req.value)
        self.assertNotEqual(0, req.value)

        method = self._string_to_py(self.fat_HttpRequestMethod(req.value))
        path = self._string_to_py(self.fat_HttpRequestPath(req.value))
        body_handle = self.fat_HttpRequestBody(req.value)
        body = self._bytes_to_py(body_handle)
        self.fat_BytesFree(body_handle)
        self.assertEqual("GET", method)
        self.assertEqual("/hello", path)
        self.assertEqual(b"", body)

        self.fat_HttpRequestFree(req.value)

        post_body = self._bytes_new(b"{\"hello\":\"world\"}")
        resp_body = self._bytes_new(b"posted")
        self.fat_HttpServerSetStaticResponse(server.value, 201, resp_body, b"application/json")
        self.fat_BytesFree(resp_body)

        resp2 = fat_string_handle_type()()
        err_resp2 = fat_string_handle_type()()
        status = self.fat_HttpClientPostBytesUTF8(
            client,
            f"{base_url}/submit".encode("utf-8"),
            b"application/json",
            post_body,
            ctypes.byref(resp2),
            ctypes.byref(err_resp2),
        )
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err_resp2.value)
        self.assertEqual(201, self.fat_HttpResponseStatus(resp2.value))

        body_handle = self.fat_HttpResponseBody(resp2.value)
        self.assertEqual(b"posted", self._bytes_to_py(body_handle))
        self.fat_BytesFree(body_handle)
        self.fat_HttpResponseFree(resp2.value)

        req2 = fat_string_handle_type()()
        err_req2 = fat_string_handle_type()()
        status = self.fat_HttpServerNextRequest(server.value, 2000, ctypes.byref(req2), ctypes.byref(err_req2))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err_req2.value)

        method2 = self._string_to_py(self.fat_HttpRequestMethod(req2.value))
        path2 = self._string_to_py(self.fat_HttpRequestPath(req2.value))
        body_handle2 = self.fat_HttpRequestBody(req2.value)
        body2 = self._bytes_to_py(body_handle2)
        self.fat_BytesFree(body_handle2)
        self.assertEqual("POST", method2)
        self.assertEqual("/submit", path2)
        self.assertEqual(b"{\"hello\":\"world\"}", body2)

        header2 = self._string_to_py(self.fat_HttpRequestHeaderGetUTF8(req2.value, b"Content-Type"))
        self.assertEqual("application/json", header2)

        self.fat_BytesFree(post_body)
        self.fat_HttpRequestFree(req2.value)

        self.fat_HttpClientFree(client)

        err_close = fat_string_handle_type()()
        status = self.fat_HttpServerClose(server.value, ctypes.byref(err_close))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err_close.value)


if __name__ == "__main__":
    unittest.main()
