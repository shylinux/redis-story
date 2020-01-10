#! /bin/sh

export PATH=${PWD}:$PATH
export ctx_pid=${ctx_pid:=var/run/ice.pid}
export ctx_log=${ctx_log:=boot.log}

prepare() {
    [ -e ice.sh ] || curl $ctx_dev/publish/ice.sh -o ice.sh && chmod u+x ice.sh
    [ -e ice.bin ] && chmod u+x ice.bin && return

    bin="ice"
    case `uname -s` in
        Darwin) bin=${bin}.darwin ;;
        Linux) bin=${bin}.linux ;;
        *) bin=${bin}.windows ;;
    esac
    case `uname -m` in
        x86_64) bin=${bin}.amd64 ;;
        i686) bin=${bin}.386 ;;
        arm*) bin=${bin}.arm ;;
    esac
    curl $ctx_dev/publish/${bin} -o ice.bin && chmod u+x ice.bin
 }
start() {
    trap HUP hup && while true; do
        date && ice.bin $@ 2>$ctx_log && echo -e "\n\nrestarting..." || break
    done
}
serve() {
    prepare && shutdown && start $@
}
restart() {
    [ -e $ctx_pid ] && kill -2 `cat $ctx_pid` || echo
}
shutdown() {
    [ -e $ctx_pid ] && kill -3 `cat $ctx_pid` || echo
}

cmd=$1 && [ -n "$cmd" ] && shift || cmd=serve
$cmd $*
