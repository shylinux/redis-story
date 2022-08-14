package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"

	"shylinux.com/x/ice"
	kit "shylinux.com/x/toolkits"
)

type redis struct {
	bio *bufio.Reader
	net.Conn
}

var ErrReadLine = errors.New("read redis line error")

func (r *redis) readLine() (line []byte, err error) {
	for {
		buf, err := r.bio.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		line = append(line, buf...)
		if len(line) > 1 && line[len(line)-2] == '\r' {
			return line[:len(line)-2], nil
		}
	}
	return nil, ErrReadLine
}
func (r *redis) readList(count int) ([]byte, error) {
	buf := make([]byte, 0, count+2)
	for begin := 0; count+2-begin > 0; {
		n, e := r.bio.Read(buf[begin : count+2])
		if begin += n; e != nil {
			return buf[0:begin], e
		}
	}
	return buf[:count], nil
}
func (r *redis) readItem(line string) (ice.Any, error) {
	switch line[0] {
	case '-': // error
		return nil, errors.New(line[1:])
	case '+': // string
		return line[1:], nil
	case '$': // bulk
		list, err := r.readList(kit.Int(line[1:]))
		return string(list), err
	case ':': // int
		return kit.Int(line[1:]), nil
	case '*': // list
		list := ice.List{}
		for i := 0; i < kit.Int(line[1:]); i++ {
			line, err := r.readLine()
			if err != nil {
				return nil, err
			}
			item, err := r.readItem(string(line))
			if err != nil {
				return nil, err
			}
			list = append(list, item)
		}
		return list, nil
	}
	return nil, nil
}

func (r *redis) Do(cmd string, arg ...string) (ice.Any, error) {
	fmt.Fprintf(r.Conn, "*%d\r\n", len(arg)+1)
	fmt.Fprintf(r.Conn, "$%d\r\n%s\r\n", len(cmd), cmd)
	for _, v := range arg {
		fmt.Fprintf(r.Conn, "$%d\r\n%s\r\n", len(v), v)
	}

	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	return r.readItem(string(line))
}
func (r *redis) Done(cmd string, arg ...string) ice.Any {
	res, _ := r.Do(cmd, arg...)
	return res
}
func (r *redis) Close() {
	r.Conn.Close()
}

func NewClient(addr string) (*redis, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &redis{Conn: conn, bio: bufio.NewReader(conn)}, nil
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

func NewRedisPool(addr string, password string) *RedisPool {
	return &RedisPool{addr: addr, Pool: sync.Pool{New: func() ice.Any {
		if c, e := NewClient(addr); e == nil {
			if password != "" {
				c.Do("auth", password)
			}
			return c
		}
		return nil
	}}}
}
