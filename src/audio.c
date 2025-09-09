#include "raylib.h"
#include "audio.h"
#include <stddef.h>

Sound audios[10];  // remove const

void LoadAudios(void) {
    audios[0] = LoadSound("assets/audio/mouse-click.mp3");
}

void UnloadAudios(void) {
    for (size_t i = 0; i < sizeof(audios) / sizeof(audios[0]); i++) {
        UnloadSound(audios[0]);
    }
}