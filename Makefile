.PHONY: all help configure build clean distclean shared static reconfigure test

BUILD_DIR ?= build
BUILD_TYPE ?= Release
SHARED ?= 0
PYTHON ?= python3

# Normalize (strip trailing slash) so comparisons work consistently.
BUILD_DIR_NORM := $(patsubst %/,%,$(BUILD_DIR))

# Build directories commonly created in this repo.
DETECTED_BUILD_DIRS := $(sort $(patsubst %/,%,$(wildcard build/ build_*/)))
ALL_BUILD_DIRS := $(sort $(DETECTED_BUILD_DIRS) $(BUILD_DIR_NORM))

# Extra flags passed to CMake at configure time, e.g.:
#   make configure CMAKE_FLAGS="-DCMAKE_C_COMPILER=clang"
CMAKE_FLAGS ?=

ifeq ($(SHARED),1)
  FATSTD_BUILD_SHARED := ON
else ifeq ($(SHARED),ON)
  FATSTD_BUILD_SHARED := ON
else ifeq ($(SHARED),on)
  FATSTD_BUILD_SHARED := ON
else
  FATSTD_BUILD_SHARED := OFF
endif

all: build

help:
	@printf "%s\n" \
	  "Targets:" \
	  "  make [all|build]        Configure (if needed) and build" \
	  "  make configure          Configure into \`$(BUILD_DIR)\`" \
	  "  make reconfigure        Re-run configure (clears CMake cache)" \
	  "  make clean              Remove all build directories" \
	  "  make distclean          Alias for \`make clean\`" \
	  "  make test               Build shared library and run Python tests" \
	  "" \
	  "Options:" \
	  "  BUILD_DIR=<dir>         Build directory (default: build)" \
	  "  BUILD_TYPE=<type>       CMake build type (default: Release)" \
	  "  SHARED=0|1              Build shared library (default: 0)" \
	  "  PYTHON=python3          Python interpreter for tests" \
	  "  CMAKE_FLAGS='...'       Extra CMake configure flags"

configure:
	cmake -S . -B "$(BUILD_DIR)" -DCMAKE_BUILD_TYPE="$(BUILD_TYPE)" -DFATSTD_BUILD_SHARED=$(FATSTD_BUILD_SHARED) $(CMAKE_FLAGS)

build: configure
	cmake --build "$(BUILD_DIR)"

reconfigure:
	@rm -f "$(BUILD_DIR)/CMakeCache.txt"
	@rm -rf "$(BUILD_DIR)/CMakeFiles"
	@$(MAKE) configure

clean:
	@test -z "$(ALL_BUILD_DIRS)" || rm -rf $(ALL_BUILD_DIRS)

distclean:
	@$(MAKE) clean

shared:
	@$(MAKE) build SHARED=1

static:
	@$(MAKE) build SHARED=0

test: shared
	$(PYTHON) scripts/python_tests/test_fatstd_shared.py
