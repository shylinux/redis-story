module github.com/shylinux/redis-story

go 1.13

require (
	github.com/shylinux/icebergs v0.1.18
	github.com/shylinux/linux-story v0.0.0-00010101000000-000000000000
	github.com/shylinux/toolkits v0.1.6
)

replace github.com/shylinux/toolkits => ../../toolkits

replace github.com/shylinux/icebergs => ../../icebergs

replace github.com/shylinux/linux-story => ../20200108-linux_story
