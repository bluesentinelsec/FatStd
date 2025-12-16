#!/usr/bin/env python3

from __future__ import annotations

import argparse
import ctypes
import os
import platform
import re
import sys
from pathlib import Path


def _fatal(message: str, *, exit_code: int = 2) -> "None":
    print(f"fatal: {message}", file=sys.stderr)
    raise SystemExit(exit_code)


def _find_repo_root(start: Path) -> Path:
    current = start
    for _ in range(10):
        if (current / "CMakeLists.txt").is_file():
            return current
        if current.parent == current:
            break
        current = current.parent
    _fatal(f"could not find repo root (no CMakeLists.txt found above {start})")


def _project_version_from_cmakelists(cmake_lists_path: Path) -> str:
    text = cmake_lists_path.read_text(encoding="utf-8", errors="replace")
    match = re.search(
        r"project\(\s*FatStd\s+VERSION\s+([0-9]+\.[0-9]+\.[0-9]+)\b", text
    )
    if not match:
        _fatal(f"could not parse project version from {cmake_lists_path}")
    return match.group(1)


def _expected_library_filenames() -> list[str]:
    system = platform.system().lower()
    if system == "windows":
        return ["fatstd.dll"]
    if system == "darwin":
        return ["libfatstd.dylib"]
    return ["libfatstd.so"]


def _discover_library_path(repo_root: Path, build_dir: Path) -> Path:
    candidate_names = _expected_library_filenames()
    search_roots = [build_dir]

    # Helpful fallback if the user has extra build dirs (e.g. build_shared).
    for extra_dir in sorted(repo_root.glob("build_*")):
        if extra_dir.is_dir() and extra_dir not in search_roots:
            search_roots.append(extra_dir)

    found: list[Path] = []
    for root in search_roots:
        if not root.is_dir():
            continue
        for name in candidate_names:
            found.extend(root.rglob(name))

    found = [p for p in found if p.is_file()]
    if not found:
        attempted = "\n".join(f"  - {p}" for p in search_roots)
        expected = ", ".join(candidate_names)
        _fatal(
            "shared library not found.\n"
            f"Expected filename(s): {expected}\n"
            "Searched directories:\n"
            f"{attempted}\n\n"
            "Build it with one of:\n"
            "  - make shared\n"
            "  - cmake -S . -B build -DFATSTD_BUILD_SHARED=ON && cmake --build build"
        )

    def sort_key(path: Path) -> tuple[int, int, str]:
        return (len(path.parts), len(str(path)), str(path))

    return sorted(found, key=sort_key)[0]


def _load_library(lib_path: Path) -> ctypes.CDLL:
    try:
        if platform.system().lower() == "windows":
            # Ensure local DLL directory is searchable for transitive deps.
            os.add_dll_directory(str(lib_path.parent))
        return ctypes.CDLL(str(lib_path))
    except OSError as exc:
        _fatal(f"failed to load shared library {lib_path}: {exc}")


def test_fat_version_string(lib: ctypes.CDLL, expected_version: str) -> None:
    lib.fat_VersionString.argtypes = []
    lib.fat_VersionString.restype = ctypes.c_char_p

    raw = lib.fat_VersionString()
    assert raw is not None, "fat_VersionString returned NULL"
    version = raw.decode("utf-8", errors="strict")

    assert version == expected_version, f"expected {expected_version!r}, got {version!r}"


def test_fat_go_add(lib: ctypes.CDLL) -> None:
    lib.fat_GoAdd.argtypes = [ctypes.c_int, ctypes.c_int]
    lib.fat_GoAdd.restype = ctypes.c_int

    got = lib.fat_GoAdd(2, 3)
    assert got == 5, f"expected 5, got {got}"


def test_fat_string_create_free(lib: ctypes.CDLL) -> None:
    fat_string = ctypes.c_size_t

    lib.fat_StringNewUTF8.argtypes = [ctypes.c_char_p]
    lib.fat_StringNewUTF8.restype = fat_string

    lib.fat_StringNewUTF8N.argtypes = [ctypes.c_void_p, ctypes.c_size_t]
    lib.fat_StringNewUTF8N.restype = fat_string

    lib.fat_StringFree.argtypes = [fat_string]
    lib.fat_StringFree.restype = None

    s1 = lib.fat_StringNewUTF8(b"lorem ipsum")
    assert s1 != 0, "fat_StringNewUTF8 returned 0 handle"
    lib.fat_StringFree(s1)

    raw = ctypes.create_string_buffer(b"abc\x00def")
    s2 = lib.fat_StringNewUTF8N(ctypes.addressof(raw), len(raw.raw))
    assert s2 != 0, "fat_StringNewUTF8N returned 0 handle"
    lib.fat_StringFree(s2)


def _run_test(name: str, fn) -> None:
    print(f"test: {name} ... ", end="", flush=True)
    try:
        fn()
    except Exception:
        print("FAIL", flush=True)
        raise
    print("ok", flush=True)


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(
        description="ctypes-based smoke tests for the FatStd shared library"
    )
    parser.add_argument(
        "--lib",
        type=Path,
        default=None,
        help="Path to shared library (overrides auto-discovery)",
    )
    parser.add_argument(
        "--build-dir",
        type=Path,
        default=None,
        help="Build directory to search (default: <repo>/build)",
    )
    args = parser.parse_args(argv)

    repo_root = _find_repo_root(Path(__file__).resolve())
    expected_version = _project_version_from_cmakelists(repo_root / "CMakeLists.txt")

    if args.lib is not None:
        lib_path = args.lib
        if not lib_path.is_file():
            _fatal(f"--lib path does not exist or is not a file: {lib_path}")
    else:
        build_dir = args.build_dir if args.build_dir is not None else repo_root / "build"
        lib_path = _discover_library_path(repo_root, build_dir)

    lib = _load_library(lib_path)

    _run_test("fat_VersionString", lambda: test_fat_version_string(lib, expected_version))
    _run_test("fat_GoAdd", lambda: test_fat_go_add(lib))
    _run_test(
        "fat_StringNewUTF8 / fat_StringNewUTF8N / fat_StringFree",
        lambda: test_fat_string_create_free(lib),
    )
    print("ok")
    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
