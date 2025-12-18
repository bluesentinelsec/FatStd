cmake_minimum_required(VERSION 3.20)

project(FatStd VERSION 0.1.0 LANGUAGES C)

option(FATSTD_BUILD_SHARED "Build FatStd as a shared library" OFF)

find_program(FATSTD_GO_EXECUTABLE go REQUIRED)

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
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstd_go/string_exports.go"
    "${CMAKE_CURRENT_SOURCE_DIR}/pkg/fatstrings/fatstrings.go"
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
)

target_include_directories(
  fatstd
  PUBLIC
    $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
    $<BUILD_INTERFACE:${fatstd_generated_include_dir}>
)

target_link_libraries(fatstd PRIVATE fatstd_go)

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
