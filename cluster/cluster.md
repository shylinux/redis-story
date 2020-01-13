{{title "redis"}}
{{brief `redis是一种简单的数据存储`}}
{{refer "官方网站" `
官网 https://redis.io
开源 https://github.com/antirez/redis
`}}

{{chapter "下载安装"}}
{{shell "下载源码" "usr" "install" `wget http://download.redis.io/releases/redis-5.0.7.tar.gz`}}
{{shell "解压源码" "usr" "install" `tar xvf redis-5.0.7.tar.gz`}}
{{shell "编译源码" "usr" "install" `cd redis-5.0.7 && make && PREFIX=../../cluster make install`}}

{{shell "启动服务" "usr/redis-5.0.7" "install" `src/redis-server`}}
{{shell "启动终端" "usr/redis-5.0.7" "install" `src/redis-cli`}}

{{chapter "项目结构"}}
{{shell "项目目录" "usr" `dir redis-5.0.7`}}
{{shell "源码目录" "usr" `dir redis-5.0.7/src`}}

{{chapter "事件模型"}}
{{order "事件模型" `
ae.c
ae.h
`}}

{{chapter "网络连接"}}
{{order "网络连接" `
anet.c
anet.h
networking.c
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

{{chapter "配置化"}}
{{order "配置化" `
config.c
config.h
`}}

{{chapter "集群化"}}
{{order "集群化" `
cluster.c
cluster.h
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
`}}

{{chapter "启动流程"}}
