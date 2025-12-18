#include <stdint.h>

// Internal linkage anchor to ensure SDL3 objects are pulled into libfatstd when
// SDL3 is linked statically (otherwise linkers may drop unused objects).

extern int SDL_Init(uint32_t flags);
extern void SDL_Quit(void);

static void *fatstd_sdl3_link_anchors[] = {(void *)&SDL_Init, (void *)&SDL_Quit};

void fatstd__sdl3_link_anchor(void) {
  (void)fatstd_sdl3_link_anchors;
}

