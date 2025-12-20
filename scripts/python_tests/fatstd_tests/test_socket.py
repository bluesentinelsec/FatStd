from __future__ import annotations

import ctypes
import threading
import unittest

from fatstd_test_support import bind, fat_string_handle_type


class TestSocket(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        fat_handle = fat_string_handle_type()
        fat_error = fat_string_handle_type()

        cls.fat_StringLenBytes = bind(
            "fat_StringLenBytes", argtypes=[fat_handle], restype=ctypes.c_size_t
        )
        cls.fat_StringCopyOut = bind(
            "fat_StringCopyOut", argtypes=[fat_handle, ctypes.c_void_p, ctypes.c_size_t], restype=ctypes.c_size_t
        )
        cls.fat_StringFree = bind("fat_StringFree", argtypes=[fat_handle], restype=None)
        cls.fat_ErrorFree = bind("fat_ErrorFree", argtypes=[fat_error], restype=None)

        cls.fat_TcpDialUTF8 = bind(
            "fat_TcpDialUTF8", argtypes=[ctypes.c_char_p, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)], restype=ctypes.c_int
        )
        cls.fat_TcpListenerListenUTF8 = bind(
            "fat_TcpListenerListenUTF8",
            argtypes=[ctypes.c_char_p, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_TcpListenerAccept = bind(
            "fat_TcpListenerAccept",
            argtypes=[fat_handle, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)],
            restype=ctypes.c_int,
        )
        cls.fat_TcpListenerAddr = bind(
            "fat_TcpListenerAddr", argtypes=[fat_handle], restype=fat_handle
        )
        cls.fat_TcpListenerClose = bind(
            "fat_TcpListenerClose", argtypes=[fat_handle, ctypes.POINTER(fat_error)], restype=ctypes.c_int
        )
        cls.fat_TcpConnRead = bind(
            "fat_TcpConnRead",
            argtypes=[
                fat_handle,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.POINTER(ctypes.c_size_t),
                ctypes.POINTER(ctypes.c_bool),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_TcpConnWrite = bind(
            "fat_TcpConnWrite",
            argtypes=[
                fat_handle,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.POINTER(ctypes.c_size_t),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_TcpConnClose = bind(
            "fat_TcpConnClose", argtypes=[fat_handle, ctypes.POINTER(fat_error)], restype=ctypes.c_int
        )

        cls.fat_UdpListenUTF8 = bind(
            "fat_UdpListenUTF8", argtypes=[ctypes.c_char_p, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)], restype=ctypes.c_int
        )
        cls.fat_UdpDialUTF8 = bind(
            "fat_UdpDialUTF8", argtypes=[ctypes.c_char_p, ctypes.POINTER(fat_handle), ctypes.POINTER(fat_error)], restype=ctypes.c_int
        )
        cls.fat_UdpConnReadFrom = bind(
            "fat_UdpConnReadFrom",
            argtypes=[
                fat_handle,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.POINTER(ctypes.c_size_t),
                ctypes.POINTER(fat_handle),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_UdpConnWriteToUTF8 = bind(
            "fat_UdpConnWriteToUTF8",
            argtypes=[
                fat_handle,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.c_char_p,
                ctypes.POINTER(ctypes.c_size_t),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_UdpConnWrite = bind(
            "fat_UdpConnWrite",
            argtypes=[
                fat_handle,
                ctypes.c_void_p,
                ctypes.c_size_t,
                ctypes.POINTER(ctypes.c_size_t),
                ctypes.POINTER(fat_error),
            ],
            restype=ctypes.c_int,
        )
        cls.fat_UdpConnLocalAddr = bind(
            "fat_UdpConnLocalAddr", argtypes=[fat_handle], restype=fat_handle
        )
        cls.fat_UdpConnClose = bind(
            "fat_UdpConnClose", argtypes=[fat_handle, ctypes.POINTER(fat_error)], restype=ctypes.c_int
        )

        cls.FAT_OK = 0

    def _string_to_py(self, h: int) -> str:
        n = self.fat_StringLenBytes(h)
        if n == 0:
            self.fat_StringFree(h)
            return ""
        buf = ctypes.create_string_buffer(n, n)
        copied = self.fat_StringCopyOut(h, ctypes.addressof(buf), len(buf.raw))
        self.assertEqual(n, copied)
        self.fat_StringFree(h)
        return buf.raw.decode("utf-8")

    def _handle_error(self, err: ctypes.c_size_t) -> None:
        if err.value != 0:
            self.fat_ErrorFree(err.value)

    def test_tcp_roundtrip(self) -> None:
        listener = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = self.fat_TcpListenerListenUTF8(b"127.0.0.1:0", ctypes.byref(listener), ctypes.byref(err))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err.value)
        self.assertNotEqual(0, listener.value)

        addr = self._string_to_py(self.fat_TcpListenerAddr(listener.value))

        server_result: dict[str, bytes] = {}

        def server() -> None:
            conn = fat_string_handle_type()()
            err_local = fat_string_handle_type()()
            status_local = self.fat_TcpListenerAccept(listener.value, ctypes.byref(conn), ctypes.byref(err_local))
            self.assertEqual(self.FAT_OK, status_local)
            self.assertEqual(0, err_local.value)

            buf = ctypes.create_string_buffer(256, 256)
            out_n = ctypes.c_size_t(0)
            out_eof = ctypes.c_bool(False)
            err_read = fat_string_handle_type()()
            status_read = self.fat_TcpConnRead(
                conn.value, ctypes.addressof(buf), len(buf.raw), ctypes.byref(out_n), ctypes.byref(out_eof), ctypes.byref(err_read)
            )
            self.assertEqual(self.FAT_OK, status_read)
            self.assertFalse(out_eof.value)
            self.assertEqual(0, err_read.value)
            server_result["received"] = buf.raw[: out_n.value]

            reply = b"pong"
            raw = ctypes.create_string_buffer(reply, len(reply))
            out_written = ctypes.c_size_t(0)
            err_write = fat_string_handle_type()()
            status_write = self.fat_TcpConnWrite(
                conn.value, ctypes.addressof(raw), len(raw.raw), ctypes.byref(out_written), ctypes.byref(err_write)
            )
            self.assertEqual(self.FAT_OK, status_write)
            self.assertEqual(len(reply), out_written.value)
            self.assertEqual(0, err_write.value)

            err_close = fat_string_handle_type()()
            status_close = self.fat_TcpConnClose(conn.value, ctypes.byref(err_close))
            self.assertEqual(self.FAT_OK, status_close)
            self.assertEqual(0, err_close.value)

        thread = threading.Thread(target=server, daemon=True)
        thread.start()

        client = fat_string_handle_type()()
        err_client = fat_string_handle_type()()
        status_client = self.fat_TcpDialUTF8(addr.encode("utf-8"), ctypes.byref(client), ctypes.byref(err_client))
        self.assertEqual(self.FAT_OK, status_client)
        self.assertEqual(0, err_client.value)

        payload = b"ping"
        raw_payload = ctypes.create_string_buffer(payload, len(payload))
        out_written = ctypes.c_size_t(0)
        err_write = fat_string_handle_type()()
        status_write = self.fat_TcpConnWrite(
            client.value, ctypes.addressof(raw_payload), len(raw_payload.raw), ctypes.byref(out_written), ctypes.byref(err_write)
        )
        self.assertEqual(self.FAT_OK, status_write)
        self.assertEqual(len(payload), out_written.value)
        self.assertEqual(0, err_write.value)

        buf = ctypes.create_string_buffer(256, 256)
        out_n = ctypes.c_size_t(0)
        out_eof = ctypes.c_bool(False)
        err_read = fat_string_handle_type()()
        status_read = self.fat_TcpConnRead(
            client.value, ctypes.addressof(buf), len(buf.raw), ctypes.byref(out_n), ctypes.byref(out_eof), ctypes.byref(err_read)
        )
        self.assertEqual(self.FAT_OK, status_read)
        self.assertFalse(out_eof.value)
        self.assertEqual(0, err_read.value)
        self.assertEqual(b"pong", buf.raw[: out_n.value])

        err_close = fat_string_handle_type()()
        status_close = self.fat_TcpConnClose(client.value, ctypes.byref(err_close))
        self.assertEqual(self.FAT_OK, status_close)
        self.assertEqual(0, err_close.value)

        thread.join(timeout=2)
        self.assertEqual(b"ping", server_result.get("received"))

        err_listener_close = fat_string_handle_type()()
        status_listener_close = self.fat_TcpListenerClose(listener.value, ctypes.byref(err_listener_close))
        self.assertEqual(self.FAT_OK, status_listener_close)
        self.assertEqual(0, err_listener_close.value)

    def test_udp_roundtrip(self) -> None:
        server_conn = fat_string_handle_type()()
        err = fat_string_handle_type()()
        status = self.fat_UdpListenUTF8(b"127.0.0.1:0", ctypes.byref(server_conn), ctypes.byref(err))
        self.assertEqual(self.FAT_OK, status)
        self.assertEqual(0, err.value)

        server_addr = self._string_to_py(self.fat_UdpConnLocalAddr(server_conn.value))

        server_result: dict[str, bytes] = {}

        def server() -> None:
            buf = ctypes.create_string_buffer(256, 256)
            out_n = ctypes.c_size_t(0)
            out_addr = fat_string_handle_type()()
            err_read = fat_string_handle_type()()
            status_read = self.fat_UdpConnReadFrom(
                server_conn.value, ctypes.addressof(buf), len(buf.raw), ctypes.byref(out_n), ctypes.byref(out_addr), ctypes.byref(err_read)
            )
            self.assertEqual(self.FAT_OK, status_read)
            self.assertEqual(0, err_read.value)
            server_result["received"] = buf.raw[: out_n.value]

            sender_addr = self._string_to_py(out_addr.value)
            reply = b"pong"
            raw_reply = ctypes.create_string_buffer(reply, len(reply))
            out_written = ctypes.c_size_t(0)
            err_write = fat_string_handle_type()()
            status_write = self.fat_UdpConnWriteToUTF8(
                server_conn.value,
                ctypes.addressof(raw_reply),
                len(raw_reply.raw),
                sender_addr.encode("utf-8"),
                ctypes.byref(out_written),
                ctypes.byref(err_write),
            )
            self.assertEqual(self.FAT_OK, status_write)
            self.assertEqual(len(reply), out_written.value)
            self.assertEqual(0, err_write.value)

        thread = threading.Thread(target=server, daemon=True)
        thread.start()

        client_conn = fat_string_handle_type()()
        err_client = fat_string_handle_type()()
        status_client = self.fat_UdpDialUTF8(server_addr.encode("utf-8"), ctypes.byref(client_conn), ctypes.byref(err_client))
        self.assertEqual(self.FAT_OK, status_client)
        self.assertEqual(0, err_client.value)

        payload = b"ping"
        raw_payload = ctypes.create_string_buffer(payload, len(payload))
        out_written = ctypes.c_size_t(0)
        err_write = fat_string_handle_type()()
        status_write = self.fat_UdpConnWrite(
            client_conn.value, ctypes.addressof(raw_payload), len(raw_payload.raw), ctypes.byref(out_written), ctypes.byref(err_write)
        )
        self.assertEqual(self.FAT_OK, status_write)
        self.assertEqual(len(payload), out_written.value)
        self.assertEqual(0, err_write.value)

        buf = ctypes.create_string_buffer(256, 256)
        out_n = ctypes.c_size_t(0)
        out_addr = fat_string_handle_type()()
        err_read = fat_string_handle_type()()
        status_read = self.fat_UdpConnReadFrom(
            client_conn.value, ctypes.addressof(buf), len(buf.raw), ctypes.byref(out_n), ctypes.byref(out_addr), ctypes.byref(err_read)
        )
        self.assertEqual(self.FAT_OK, status_read)
        self.assertEqual(0, err_read.value)
        self.assertEqual(b"pong", buf.raw[: out_n.value])
        self.fat_StringFree(out_addr.value)

        thread.join(timeout=2)
        self.assertEqual(b"ping", server_result.get("received"))

        err_close = fat_string_handle_type()()
        status_close = self.fat_UdpConnClose(client_conn.value, ctypes.byref(err_close))
        self.assertEqual(self.FAT_OK, status_close)
        self.assertEqual(0, err_close.value)

        err_server_close = fat_string_handle_type()()
        status_server_close = self.fat_UdpConnClose(server_conn.value, ctypes.byref(err_server_close))
        self.assertEqual(self.FAT_OK, status_server_close)
        self.assertEqual(0, err_server_close.value)


if __name__ == "__main__":
    unittest.main()
