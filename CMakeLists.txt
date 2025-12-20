cmake_minimum_required(VERSION 3.20)

project(FatStd VERSION 0.1.0 LANGUAGES C)

option(FATSTD_BUILD_SHARED "Build FatStd as a shared library" OFF)

find_program(FATSTD_GO_EXECUTABLE go REQUIRED)

include(FetchContent)

set(fatstd_generated_include_dir "${CMAKE_CURRENT_BINARY_DIR}/generated/include")
file(MAKE_DIRECTORY "${fatstd_generated_include_dir}/fat")

configure_file(
  "${CMAKE_CURRENT_SOURCE_DIR}/include/fat/version_generated.h.in"
  "${fatstd_generated_include_dir}/fat/version_generated.h"
  @ONLY
)

set(fatstd_go_out_dir "${CMAKE_CURRENT_BINARY_DIR}/generated/go")
file(MAKE_DIRECTORY "${fatstd_go_out_dir}")

set(fatstd_go_archive "${fatstd_go_out_dir}/fatstd_go${CMAKE_STATIC_LIBRARY_SUFFIX}")
set(fatstd_go_header "${fatstd_go_out_dir}/fatstd_go.h")

add_custom_command(
  OUTPUT "${fatstd_go_archive}" "${fatstd_go_header}"
  COMMAND "${CMAKE_COMMAND}" -E make_directory "${fatstd_go_out_dir}"
  COMMAND
    "${CMAKE_COMMAND}" -E env
      "GOCACHE=${fatstd_go_out_dir}/gocache"
      "GOMODCACHE=${fatstd_go_out_dir}/gomodcache"
      "GOPATH=${fatstd_go_out_dir}/gopath"
      "${FATSTD_GO_EXECUTABLE}" build
        -buildmode=c-archive
        -o "${fatstd_go_archive}"
        "./pkg/fatstd_go"
  WORKING_DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}"
  DEPENDS
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/fatstd_go.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/handles.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/error_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/string_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/bytes_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/bytes_buffer_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/bytes_reader_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/conv_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/zip_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/tar_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/base64_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/csv_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstrings/fatstrings.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatbytes/fatbytes.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatbytes/buffer.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatbytes/reader.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatconv/fatconv.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/go.mod"
  VERBATIM
)

add_custom_target(fatstd_go_build DEPENDS "${fatstd_go_archive}" "${fatstd_go_header}")

add_library(fatstd_go STATIC IMPORTED GLOBAL)
set_target_properties(
  fatstd_go
  PROPERTIES
    IMPORTED_LOCATION "${fatstd_go_archive}"
    INTERFACE_INCLUDE_DIRECTORIES "${fatstd_go_out_dir}"
)
add_dependencies(fatstd_go fatstd_go_build)

if(FATSTD_BUILD_SHARED)
  add_library(fatstd SHARED)
else()
  add_library(fatstd STATIC)
endif()

add_library(fatstd::fatstd ALIAS fatstd)

target_sources(
  fatstd
  PRIVATE
    src/fat_version.c
    src/fat_go.c
    src/fat_string.c
    src/fat_string_builder.c
    src/fat_string_reader.c
    src/fat_error.c
    src/fat_bytes.c
    src/fat_bytes_buffer.c
    src/fat_bytes_reader.c
    src/fat_conv.c
    src/fat_base64.c
    src/fat_csv.c
    src/fat_zip.c
    src/fat_tar.c
    src/fat_sdl3_link.c
    src/fat_raylib_link.c
    src/fat_raygui.c
    src/getopt.c
)

target_include_directories(
  fatstd
  PUBLIC
    $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
    $<BUILD_INTERFACE:${fatstd_generated_include_dir}>
)

target_link_libraries(fatstd PRIVATE fatstd_go)

# --- SDL3 (FetchContent) ---
set(SDL_SHARED OFF CACHE BOOL "Build SDL as a shared library" FORCE)
set(SDL_STATIC ON CACHE BOOL "Build SDL as a static library" FORCE)
set(SDL_TESTS OFF CACHE BOOL "Build SDL tests" FORCE)
set(SDL_TEST_LIBRARY OFF CACHE BOOL "Build SDL test library" FORCE)

FetchContent_Declare(
  SDL3
  GIT_REPOSITORY https://github.com/libsdl-org/SDL.git
  GIT_TAG release-3.2.28
)

FetchContent_MakeAvailable(SDL3)

set(fatstd_sdl3_target "")
if(TARGET SDL3::SDL3-static)
  set(fatstd_sdl3_target SDL3::SDL3-static)
elseif(TARGET SDL3::SDL3)
  set(fatstd_sdl3_target SDL3::SDL3)
elseif(TARGET SDL3-static)
  set(fatstd_sdl3_target SDL3-static)
elseif(TARGET SDL3)
  set(fatstd_sdl3_target SDL3)
else()
  message(FATAL_ERROR "SDL3 target not found after FetchContent_MakeAvailable(SDL3)")
endif()

# Expose SDL3 to consumers linking fatstd, and prefer static SDL linkage.
target_link_libraries(fatstd PUBLIC ${fatstd_sdl3_target})

# When SDL is built as a static archive and fatstd is a shared library, linkers may
# drop SDL object files if no symbols are referenced directly by fatstd consumers.
# Force SDL to be pulled into libfatstd.
if(FATSTD_BUILD_SHARED)
  if(APPLE)
    target_link_options(fatstd PRIVATE "LINKER:-force_load,$<TARGET_FILE:${fatstd_sdl3_target}>")
  elseif(MSVC)
    target_link_options(fatstd PRIVATE "/WHOLEARCHIVE:$<TARGET_FILE:${fatstd_sdl3_target}>")
  else()
    target_link_options(fatstd PRIVATE "LINKER:--whole-archive,$<TARGET_FILE:${fatstd_sdl3_target}>,--no-whole-archive")
  endif()
endif()

# --- raylib (FetchContent) ---
set(BUILD_SHARED_LIBS OFF CACHE BOOL "Build libraries as shared" FORCE)
set(BUILD_EXAMPLES OFF CACHE BOOL "Build raylib examples" FORCE)
set(BUILD_GAMES OFF CACHE BOOL "Build raylib games" FORCE)

FetchContent_Declare(
  raylib
  GIT_REPOSITORY https://github.com/raysan5/raylib.git
  GIT_TAG 5.5
)

if(APPLE)
  enable_language(OBJC)
endif()

FetchContent_MakeAvailable(raylib)

set(fatstd_raylib_target "")
if(TARGET raylib)
  set(fatstd_raylib_target raylib)
elseif(TARGET raylib_static)
  set(fatstd_raylib_target raylib_static)
else()
  message(FATAL_ERROR "raylib target not found after FetchContent_MakeAvailable(raylib)")
endif()

# raylib is statically linked into libfatstd (symbols are available to consumers),
# but we keep the CMake dependency PRIVATE so exporting fatstd doesn't require exporting
# raylib and its internal dependency targets (e.g. glfw).
target_link_libraries(fatstd PRIVATE ${fatstd_raylib_target})

if(FATSTD_BUILD_SHARED)
  if(APPLE)
    target_link_options(fatstd PRIVATE "LINKER:-force_load,$<TARGET_FILE:${fatstd_raylib_target}>")
  elseif(MSVC)
    target_link_options(fatstd PRIVATE "/WHOLEARCHIVE:$<TARGET_FILE:${fatstd_raylib_target}>")
  else()
    target_link_options(fatstd PRIVATE "LINKER:--whole-archive,$<TARGET_FILE:${fatstd_raylib_target}>,--no-whole-archive")
  endif()
endif()

set_target_properties(
  fatstd
  PROPERTIES
    C_STANDARD 17
    C_STANDARD_REQUIRED YES
    C_EXTENSIONS NO
)

if(FATSTD_BUILD_SHARED)
  target_compile_definitions(fatstd PUBLIC FATSTD_SHARED=1)
  target_compile_definitions(fatstd PRIVATE FATSTD_BUILDING=1)
endif()

if(MSVC)
  target_compile_options(fatstd PRIVATE /W4)
else()
  target_compile_options(fatstd PRIVATE -Wall -Wextra -Wpedantic)
endif()

include(CMakePackageConfigHelpers)

set(fatstd_package_dir "${CMAKE_CURRENT_BINARY_DIR}/cmake")
file(MAKE_DIRECTORY "${fatstd_package_dir}")

configure_package_config_file(
  "${CMAKE_CURRENT_SOURCE_DIR}/cmake/FatStdConfig.cmake.in"
  "${fatstd_package_dir}/FatStdConfig.cmake"
  INSTALL_DESTINATION "${fatstd_package_dir}"
  NO_CHECK_REQUIRED_COMPONENTS_MACRO
)

write_basic_package_version_file(
  "${fatstd_package_dir}/FatStdConfigVersion.cmake"
  VERSION "${PROJECT_VERSION}"
  COMPATIBILITY SameMajorVersion
)

export(
  TARGETS fatstd
  NAMESPACE fatstd::
  FILE "${fatstd_package_dir}/FatStdTargets.cmake"
)
