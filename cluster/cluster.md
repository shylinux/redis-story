{{title "redis"}}
{{brief `redis是一种简单的数据存储`}}
{{refer "官方网站" `
官网 https://redis.io
开源 https://github.com/antirez/redis
`}}

{{shell "下载源码" "usr" "install" `wget http://download.redis.io/releases/redis-5.0.7.tar.gz`}}

{{shell "解压源码" "usr" "install" `tar xvf redis-5.0.7.tar.gz`}}
{{shell "解压源码" "usr" "install" `sudo yum install ruby`}}

{{shell "编译源码" "usr" "install" `cd redis-5.0.7 && make && PREFIX=../../cluster make install`}}

{{shell "启动服务" "usr/redis-5.0.7" "install" `src/redis-server &`}}

{{shell "启动终端" "usr/redis-5.0.7" "install" `src/redis-cli`}}
