from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import get_context


class TestSDL3(unittest.TestCase):
    def test_init_quit(self) -> None:
        lib = get_context().lib

        SDL_Init = lib.SDL_Init
        SDL_Init.argtypes = [ctypes.c_uint32]
        SDL_Init.restype = ctypes.c_bool

        SDL_Quit = lib.SDL_Quit
        SDL_Quit.argtypes = []
        SDL_Quit.restype = None

        ok = SDL_Init(0)
        self.assertTrue(ok)
        SDL_Quit()
