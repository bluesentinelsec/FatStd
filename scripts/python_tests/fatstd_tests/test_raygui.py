from __future__ import annotations

import ctypes
import unittest

from fatstd_test_support import get_context


class TestRaygui(unittest.TestCase):
    def test_raygui_symbols(self) -> None:
        lib = get_context().lib

        GuiSetState = lib.GuiSetState
        GuiSetState.argtypes = [ctypes.c_int]
        GuiSetState.restype = None

        GuiGetState = lib.GuiGetState
        GuiGetState.argtypes = []
        GuiGetState.restype = ctypes.c_int

        prev = GuiGetState()
        try:
            GuiSetState(3)
            self.assertEqual(3, GuiGetState())
        finally:
            GuiSetState(prev)

    def test_raylib_and_raygui_together(self) -> None:
        lib = get_context().lib

        GetRandomValue = lib.GetRandomValue
        GetRandomValue.argtypes = [ctypes.c_int, ctypes.c_int]
        GetRandomValue.restype = ctypes.c_int

        GuiGetState = lib.GuiGetState
        GuiGetState.argtypes = []
        GuiGetState.restype = ctypes.c_int

        self.assertEqual(7, GetRandomValue(7, 7))
        self.assertIsInstance(GuiGetState(), int)

