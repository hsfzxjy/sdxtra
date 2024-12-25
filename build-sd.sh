#!/bin/bash

ROOT=$(dirname $(readlink -f $0))

COMMON_ARGS=(
    "-DCMAKE_CUDA_ARCHITECTURES=86"
    "-DSD_CUBLAS=ON"
    "-DSD_BUILD_EXAMPLES=OFF"
)

function build-shared {
    cd $ROOT/sdcpp/
    mkdir build/shared/ -p
    cd build/shared/
    cmake ../../ -DSD_BUILD_SHARED_LIBS=ON "${COMMON_ARGS[@]}"
    cmake --build . --parallel 12 --config Release
}

function build-static {
    cd $ROOT/sdcpp/
    mkdir build/static/ -p
    cd build/static/
    cmake ../../ "${COMMON_ARGS[@]}"
    cmake --build . --parallel 12 --config Release
}

build-shared
# build-static
