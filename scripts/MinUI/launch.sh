#!/bin/sh
CUR_DIR="$(dirname "$0")"
cd "$CUR_DIR" || exit 1

# Apply pending update
if [ -d "../.update" ]; then
    cp -rf ../.update/* ..
    rm -rf ../.update
fi

export CFW=MinUI

ARCH=$(uname -m)
case "$ARCH" in
    aarch64|arm64)
        export LD_LIBRARY_PATH=$CUR_DIR/lib64:$LD_LIBRARY_PATH
        ./grout64
        ;;
    armv7*|armhf)
        export IS_MIYOO=1
        export SDL_VIDEODRIVER=mmiyoo
        export SDL_AUDIODRIVER=mmiyoo
        export EGL_VIDEODRIVER=mmiyoo
        export SDL_MMIYOO_DOUBLE_BUFFER=1
        export LD_LIBRARY_PATH=$CUR_DIR/lib32:$LD_LIBRARY_PATH
        ./grout32
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac
