#!/bin/sh
CUR_DIR="$(dirname "$0")"
cd "$CUR_DIR"/grout || exit 1

export CFW=ONION
export LD_LIBRARY_PATH="$CUR_DIR"/grout/lib:$LD_LIBRARY_PATH

export SDL_VIDEODRIVER=mmiyoo
export SDL_AUDIODRIVER=mmiyoo
export EGL_VIDEODRIVER=mmiyoo
export SDL_MMIYOO_DOUBLE_BUFFER=1

./grout
