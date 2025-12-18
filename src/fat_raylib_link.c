#include <stdint.h>

// Internal linkage anchor to ensure raylib objects are pulled into libfatstd when
// raylib is linked statically (otherwise linkers may drop unused objects).

extern int GetRandomValue(int min, int max);

static void *fatstd_raylib_link_anchors[] = {(void *)&GetRandomValue};

void fatstd__raylib_link_anchor(void) {
  (void)fatstd_raylib_link_anchors;
}

