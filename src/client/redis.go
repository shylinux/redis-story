package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"

	kit "github.com/shylinux/toolkits"
)

type redis struct {
	bio *bufio.Scanner
	net.Conn
}

func (r *redis) Do(arg ...string) (interface{}, error) {
	if len(arg) > 1 {
		fmt.Fprintf(r.Conn, "*%d\r\n", len(arg))
	}
	for _, v := range arg {
		fmt.Fprintf(r.Conn, "$%d\r\n%s\r\n", len(v), v)
	}

	r.bio.Scan()
	line := r.bio.Text()
	if len(line) == 0 {
		return nil, nil
	}
	switch line[0] {
	case '-':
		return nil, errors.New(line[1:])
	case '+':
		return line[1:], nil
	case '$':
		list := []string{}
		for rest := kit.Int(line[1:]); rest > 0; {
			r.bio.Scan()
			list = append(list, r.bio.Text())
			rest -= len(r.bio.Text()) + 1
		}
		return strings.Join(list, "\n"), nil
	case ':':
		return kit.Int(line[1:]), nil
	case '*':
		list := []string{}
		data := []int{}
		for i := 0; i < kit.Int(line[1:]); i++ {
			r.bio.Scan()
			line := r.bio.Text()
			switch line[0] {
			case '$':
				r.bio.Scan()
				list = append(list, r.bio.Text())
			case ':':
				list = append(list, line[1:])
			case '*':
			}
		}
		if len(data) > 0 {
			return data, nil
		}
		return list, nil
	}
	return nil, nil
}
func (r *redis) Close() {
	r.Conn.Close()
}

func NewClient(addr string) (*redis, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &redis{Conn: conn, bio: bufio.NewScanner(conn)}, nil
}

type RedisPool struct {
	addr string
	sync.Pool
}

func (rp *RedisPool) Get() *redis {
	rc, _ := rp.Pool.Get().(*redis)
	return rc
}
func (rp *RedisPool) Put(r *redis) {
	rp.Pool.Put(r)
}

func NewRedisPool(addr string) *RedisPool {
	return &RedisPool{addr: addr, Pool: sync.Pool{New: func() interface{} {
		if c, e := NewClient(addr); e == nil {
			return c
		}
		return nil
	}}}
}
