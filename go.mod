module shylinux.com/x/redis-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

require (
	shylinux.com/x/ice v1.3.11
	shylinux.com/x/icebergs v1.5.16
	shylinux.com/x/toolkits v0.7.8
)

require (
	github.com/99designs/keyring v1.2.1
	github.com/apache/pulsar-client-go v0.10.0
	github.com/sirupsen/logrus v1.6.0
)
