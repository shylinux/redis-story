{{title "redis"}}
{{brief `redis是一种简单的数据存储`}}
{{refer "官方网站" `
官网 https://redis.io
开源 https://github.com/antirez/redis
`}}

{{chapter "下载安装"}}
{{shell "下载源码" "usr" "install" ` wget http://download.redis.io/releases/redis-5.0.7.tar.gz`}}
{{shell "解压源码" "usr" "install" `tar xvf redis-5.0.7.tar.gz`}}
{{shell "编译源码" "usr" "install" `cd redis-5.0.7 && make && PREFIX=../../cluster make install`}}

{{shell "启动服务" "usr/redis-5.0.7" "install" `src/redis-server`}}
{{shell "启动终端" "usr/redis-5.0.7" "install" `src/redis-cli`}}

{{chapter "项目结构"}}
{{shell "项目目录" "usr" `dir redis-5.0.7`}}
{{shell "源码目录" "usr" `dir redis-5.0.7/src`}}
{{shell "生成索引" "usr" "install" `ctags -a tags *`}}

{{chapter "启动流程"}}
{{order "启动流程" `
server.h
server.c
redis-cli.c
`}}

{{stack "启动流程" `
main() bg red
    initServerConfig()
        server:redisServer
    initSentinelConfig()
    initSentinel()
        sentinel:sentinelState
    loadServerConfig()
        loadServerConfigFromString()
            lines[i]=for_sdssplitlen()
            argv=sdssplitargs(lines[i])
            server.port=atoi(argv[1])
    daemonize()
    initServer()
        setupSignalHandlers()
        openlog()
        aeCreateEventLoop()
            aeApiCreate()
                epoll_create()
        listenToPort()
            anetTcpServer()
        anetUnixServer()
        aeCreateTimeEvent()
        aeCreateFileEvent(acceptTcpHandler)
        aeCreateFileEvent(acceptUnixHandler)
        clusterInit()
            clusterLoadConfig()
        replicationScriptCacheInit()
        scriptingInit(1)
            lua_open()
        slowlogInit()
        latencyMonitorInit()
    moduleLoadFromQueue()
    loadDataFromDisk()
    aeSetBeforeSleepProc()
    aeSetAfterSleepProc()
    aeMain(eventLoop) bg red
        while(!eventLoop->stop)
            eventLoop.beforesleep(eventLoop)/beforeSleep(eventLoop)
                clusterBeforeSleep()
                activeExpireCycle()
                handleClientsWithPendingWrites() bg green
                    sendReplyToClient(c)
                        writeToClient(c)
                            write(c->buf)
            aeProcessEvents()
                aeGetTime()
                aeApiPoll()
                processTimeEvents()
                fe.rfileProc()/acceptTcpHandler() bg green
                    anetTcpAccept()
                    acceptCommonHandler()
                        createClient()
                            aeCreateFileEvent(readQueryFromClient)
                fe.rfileProc()/readQueryFromClient() bg green
                    read()
                    processInputBufferAndReplicate(c)
                        processInputBuffer(c)
                            processCommand()
                                lookupCommand(c.argv[0])
                                    dictFetchValue(server.commands,name)/redisCommand{}
                                call(c)
                                    c.cmd.proc(c)
                                        getComand(c)/getGenericCommand(c)
                                        lookupKeyReadOrReply(c,c.argv[1],shared.nullbulk)
                                            lookupKeyRead(c.db,key)
                                            addReply(c,reply)
                fe.wfileProc()
`}}


{{chapter "网络连接"}}
{{order "网络连接" `
anet.c
anet.h
networking.c
`}}

{{chapter "事件模型"}}
{{order "事件模型" `
ae.c
ae.h
`}}

{{chapter "消息订阅"}}
{{order "消息订阅" `
pubsub.c
`}}

{{chapter "数据库"}}
{{order "数据库" `
db.c
expire.c
`}}

{{chapter "模块化"}}
{{order "模块化" `
module.c
`}}

{{chapter "配置化"}}
{{order "配置化" `
config.c
config.h
`}}

{{chapter "脚本化"}}
{{order "脚本化" `
scripting.c
`}}

{{chapter "集群化"}}
{{order "集群化" `
cluster.c
cluster.h
sentinel.c
replication.c
`}}

{{chapter "序列化"}}
{{order "序列化" `
aof.c
rdb.c
`}}

{{chapter "数据类型"}}
{{shell "数据类型" "usr" `ls redis-5.0.7/src/t_*.c`}}

{{chapter "数据结构"}}
{{order "数据结构" `
sds.c
sds.h
object.c
adlist.c
adlist.h
rax.c
rax.h
dict.c
dict.h
intset.c
intset.h
listpack.c
listpack.h
quicklist.c
quicklist.h
ziplist.c
ziplist.h
zipmap.c
zipmap.h
`}}


