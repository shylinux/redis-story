module shylinux.com/x/redis-story

go 1.13

replace (
	shylinux.com/x/ice => ./usr/release
	shylinux.com/x/icebergs => ./usr/icebergs
	shylinux.com/x/toolkits => ./usr/toolkits
)

require (
	shylinux.com/x/ice v1.4.4
	shylinux.com/x/icebergs v1.8.4
	shylinux.com/x/toolkits v0.8.4
)

require (
	github.com/apache/pulsar-client-go v0.10.0
	github.com/segmentio/kafka-go v0.4.42
)
