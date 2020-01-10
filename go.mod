module redis-story

go 1.13

require (
	github.com/shylinux/icebergs v0.1.8
	github.com/shylinux/redis-story v0.0.0-00010101000000-000000000000
	github.com/shylinux/toolkits v0.1.0
)

replace github.com/shylinux/redis-story => ./

replace github.com/shylinux/icebergs => ../../icebergs
