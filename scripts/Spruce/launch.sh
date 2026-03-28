#!/bin/sh
CUR_DIR="$(dirname "$0")"
cd "$CUR_DIR"/grout || exit 1

# Apply pending update
if [ -d "../../.update" ]; then
    cp -rf ../../.update/* ../..
    rm -rf ../../.update
fi

export CFW=SPRUCE

case "$PLATFORM" in

############################################################
# A30
############################################################
    "A30" )
        echo "A30 detected, setting up environment variables for SDL2"
        export LD_LIBRARY_PATH="/mnt/SDCARD/spruce/a30/sdl2:$LD_LIBRARY_PATH"
        # TODO: try to keep miyoo, but add the subdirectory `sdl2` like in the spruce a30
        # so... even If i take the content of the spruce a30/sdl2 folder, it does not launch
        # I really don't understand why
#        export LD_LIBRARY_PATH="$CUR_DIR/grout/lib32/miyoo:$LD_LIBRARY_PATH"
        export LD_LIBRARY_PATH="$CUR_DIR/grout/lib32/a30:$LD_LIBRARY_PATH"
        export SPRUCE_DEVICE="A30"

        ./grout32
    ;;

############################################################
# Brick / SmartPro / SmartProS
############################################################
    "Brick" | "SmartPro" | "SmartProS" )
        echo "Brick/SmartPro/SmartProS detected, setting up environment variables for SDL2"
        export LD_LIBRARY_PATH="$CUR_DIR/grout/lib64:$LD_LIBRARY_PATH"
        export SPRUCE_DEVICE="TRIMUI"
        ./grout64
    ;;

############################################################
# Miyoo Flip
############################################################
    "Flip" )
        echo "Miyoo Flip detected, setting up environment variables for SDL2"
        export LD_LIBRARY_PATH="$CUR_DIR/grout/lib64:$LD_LIBRARY_PATH"
        export SPRUCE_DEVICE="MIYOOFLIP"
        ./grout64
    ;;
############################################################
# Miyoo Mini Flip
############################################################
    "MiyooMini" )
        echo "Miyoo Mini detected, setting up environment variables for SDL2"

        export SDL_VIDEODRIVER=mmiyoo
        export SDL_AUDIODRIVER=mmiyoo
        export EGL_VIDEODRIVER=mmiyoo
        export SDL_MMIYOO_DOUBLE_BUFFER=1

        export IS_MIYOO=1
        export SPRUCE_DEVICE="MIYOOMINI"
        export LD_LIBRARY_PATH="$CUR_DIR/grout/lib32/miyoo:$LD_LIBRARY_PATH"


        ./grout32
    ;;
esac
