#include "raylib.h"
#include "icons.h"

Texture2D iconSprite;
Texture2D iconMap;
Texture2D iconMusic;
Texture2D iconPlay;
Texture2D iconCode;
Texture2D iconUndefined;

void LoadIcons(void)
{
    iconSprite = LoadTexture("assets/icons/brush.png");
    iconMap = LoadTexture("assets/icons/tile.png");
    iconMusic = LoadTexture("assets/icons/music.png");
    iconPlay = LoadTexture("assets/icons/play.png");
    iconCode = LoadTexture("assets/icons/code.png");
    iconUndefined = LoadTexture("assets/icons/undefined.png");
}

void UnloadIcons(void)
{
    UnloadTexture(iconSprite);
    UnloadTexture(iconMap);
    UnloadTexture(iconMusic);
    UnloadTexture(iconPlay);
    UnloadTexture(iconCode);
    UnloadTexture(iconUndefined);
}