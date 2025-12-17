from __future__ import annotations

import ctypes
from dataclasses import dataclass


@dataclass(frozen=True)
class FatStdTestContext:
    lib: ctypes.CDLL
    expected_version: str


_CTX: FatStdTestContext | None = None


def set_context(ctx: FatStdTestContext) -> None:
    global _CTX
    _CTX = ctx


def get_context() -> FatStdTestContext:
    if _CTX is None:
        raise RuntimeError("FatStd test context not initialized (did you run test_fatstd_shared.py?)")
    return _CTX


def fat_string_handle_type():
    return ctypes.c_size_t


def bind(name: str, *, argtypes: list, restype):
    lib = get_context().lib
    fn = getattr(lib, name)
    fn.argtypes = argtypes
    fn.restype = restype
    return fn

