#!/bin/bash

[ -f ~/.ish/plug.sh ] || [ -f $PWD/.ish/plug.sh ] || git clone ${ISH_CONF_HUB_PROXY:="https://"}github.com/shylinux/intshell $PWD/.ish
[ "$ISH_CONF_PRE" != "" ] || source $PWD/.ish/plug.sh || source ~/.ish/plug.sh
require miss.sh

ish_miss_prepare_compile
ish_miss_prepare_install

# ish_miss_prepare_icebergs
# ish_miss_prepare_toolkits

make
