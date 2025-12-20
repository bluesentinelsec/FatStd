#include "fat/error.h"

#include "fatstd_go.h"

void fat_ErrorFree(fat_Error e) {
  fatstd_go_error_free((uintptr_t)e);
}

fat_String fat_ErrorMessage(fat_Error e) {
  return (fat_String)fatstd_go_error_message((uintptr_t)e);
}

