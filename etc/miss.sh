#!/bin/bash
if [ "$ISH_CONF_PRE" = "" ]; then
    [ -f $PWD/.ish/plug.sh ] || [ -f $HOME/.ish/plug.sh ] || git clone ${ISH_CONF_HUB_PROXY:="https://"}shylinux.com/x/intshell $PWD/.ish
    source $PWD/.ish/plug.sh || source $HOME/.ish/plug.sh
fi

require miss.sh
ish_miss_prepare_compile
ish_miss_prepare_develop
ish_miss_prepare_install

ish_miss_prepare release
ish_miss_prepare_icebergs
ish_miss_prepare_toolkits

make
