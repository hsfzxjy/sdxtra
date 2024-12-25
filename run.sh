#!/bin/bash

ROOT_DIR=$(dirname $(readlink -f $0))

export LD_PRELOAD="${ROOT_DIR}/sdcpp/build/shared/bin/libstable-diffusion.so"
# stat $LD_PRELOAD
exec $@
