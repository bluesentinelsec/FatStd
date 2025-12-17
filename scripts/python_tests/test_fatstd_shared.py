#!/usr/bin/env python3

from __future__ import annotations

import argparse
import ctypes
import os
import platform
import re
import sys
import unittest
from pathlib import Path

from fatstd_test_support import FatStdTestContext, set_context


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


def main(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description="unittest smoke tests for FatStd shared lib")
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
    parser.add_argument("-v", "--verbose", action="store_true", help="Verbose test output")
    args, unittest_args = parser.parse_known_args(argv)

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

    set_context(FatStdTestContext(lib=lib, expected_version=expected_version))

    python_tests_dir = Path(__file__).resolve().parent
    start_dir = python_tests_dir / "fatstd_tests"
    suite = unittest.defaultTestLoader.discover(
        start_dir=str(start_dir),
        pattern="test_*.py",
        top_level_dir=str(python_tests_dir),
    )
    runner = unittest.TextTestRunner(verbosity=2 if args.verbose else 1)
    result = runner.run(suite)
    return 0 if result.wasSuccessful() else 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
