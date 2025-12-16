#pragma once

/*
  Public symbol visibility for FatStd.

  - Static builds: FATSTD_API expands to nothing.
  - Shared builds:
      - Windows: uses dllimport/dllexport.
      - Others: uses default visibility.
*/

#if defined(FATSTD_SHARED) && FATSTD_SHARED
  #if defined(_WIN32) || defined(__CYGWIN__)
    #if defined(FATSTD_BUILDING)
      #define FATSTD_API __declspec(dllexport)
    #else
      #define FATSTD_API __declspec(dllimport)
    #endif
  #elif defined(__GNUC__) || defined(__clang__)
    #define FATSTD_API __attribute__((visibility("default")))
  #else
    #define FATSTD_API
  #endif
#else
  #define FATSTD_API
#endif

