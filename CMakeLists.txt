cmake_minimum_required(VERSION 3.20)

project(FatStd VERSION 0.1.0 LANGUAGES C)

option(FATSTD_BUILD_SHARED "Build FatStd as a shared library" OFF)

set(fatstd_generated_include_dir "${CMAKE_CURRENT_BINARY_DIR}/generated/include")
file(MAKE_DIRECTORY "${fatstd_generated_include_dir}/fat")

configure_file(
  "${CMAKE_CURRENT_SOURCE_DIR}/include/fat/version_generated.h.in"
  "${fatstd_generated_include_dir}/fat/version_generated.h"
  @ONLY
)

if(FATSTD_BUILD_SHARED)
  add_library(fatstd SHARED)
else()
  add_library(fatstd STATIC)
endif()

add_library(fatstd::fatstd ALIAS fatstd)

target_sources(fatstd PRIVATE src/fat_version.c)

target_include_directories(
  fatstd
  PUBLIC
    $<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}/include>
    $<BUILD_INTERFACE:${fatstd_generated_include_dir}>
)

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
