#!/bin/bash

require &>/dev/null || if [ -f $PWD/.ish/plug.sh ]; then source $PWD/.ish/plug.sh; elif [ -f $HOME/.ish/plug.sh ]; then source $HOME/.ish/plug.sh; else
	ctx_temp=$(mktemp); if curl -h &>/dev/null; then curl -o $ctx_temp -fsSL https://shylinux.com; else wget -O $ctx_temp -q http://shylinux.com; fi; source $ctx_temp intshell
fi

require miss.sh
ish_miss_prepare_compile
ish_miss_prepare_develop
ish_miss_prepare_operate

ish_miss_prepare release
ish_miss_prepare_icebergs
ish_miss_prepare_toolkits

ish_miss_make; if [ -n "$*" ]; then ish_miss_serve "$@"; fi
