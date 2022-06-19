#!/bin/bash

[ -f $PWD/.ish/plug.sh ] || [ -f $HOME/.ish/plug.sh ] || git clone ${ISH_CONF_HUB_PROXY:="https://"}shylinux.com/x/intshell $PWD/.ish
if [ "$ISH_CONF_PRE" = "" ]; then
    source $PWD/.ish/plug.sh || source $HOME/.ish/plug.sh
fi

require miss.sh
ish_miss_prepare_compile
ish_miss_prepare_develop
ish_miss_prepare_install

ish_miss_make; if [ -n "$*" ]; then ./bin/ice.bin forever serve "$@"; fi
