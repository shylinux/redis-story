package server

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/redis-story/src/client"
	kit "shylinux.com/x/toolkits"
	"shylinux.com/x/toolkits/conf"
	"shylinux.com/x/toolkits/logs"
	"shylinux.com/x/toolkits/task"
)

type Stat struct {
	NReq int64
	NErr int64
	NOK  int64

	NRead  int64
	NWrite int64

	BeginTime time.Time
	EndTime   time.Time

	Cost time.Duration

	Up   float64
	Down float64
	QPS  float64
	AVG  time.Duration
}

var trans = map[string]func(i int64) ice.List{
	"GET": func(i int64) ice.List { return ice.List{fmt.Sprintf("hi%d", i)} },
	"SET": func(i int64) ice.List { return ice.List{fmt.Sprintf("hi%d", i), "hello"} },
}
var ErrNotFound = errors.New("not found cmd")

func Bench(nconn, nreq int64, hosts []string, cmds []string, check func(cmd string, arg ice.List, res ice.Any)) (*Stat, error) {
	s := &Stat{BeginTime: time.Now()}
	defer func() { // 请求统计
		if s.EndTime = time.Now(); s.BeginTime != s.EndTime {
			d := float64(s.EndTime.Sub(s.BeginTime)) / float64(time.Second)

			s.QPS = float64(s.NReq) / d
			s.AVG = s.EndTime.Sub(s.BeginTime) / time.Duration(nreq)
			s.Down = float64(s.NRead) / d
			s.Up = float64(s.NWrite) / d
		}
	}()

	// 连接池
	rp := client.NewRedisPool(hosts[0], "")
	// rp := redis.NewPool(func() (redis.Conn, error) { return redis.Dial("tcp", hosts[0]) }, 10)

	// 协程池
	tp := task.New(conf.Sub(task.TASK))
	defer tp.Close()

	tp.WaitN(int(nconn), func(task *task.Task, lock *task.Lock) error {
		var nerr, nok int64
		defer func() { // 请求汇总
			atomic.AddInt64(&s.NReq, nreq)
			atomic.AddInt64(&s.NErr, nerr)
			atomic.AddInt64(&s.NOK, nok)
		}()

		conn := rp.Get()
		defer rp.Put(conn)
		// defer conn.Close()

		cmd := strings.ToUpper(cmds[0])
		method := trans[cmd]
		if method == nil {
			return ErrNotFound
		}

		for i := int64(0); i < nreq; i++ {
			func() {
				defer logs.CostTime(func(d time.Duration) {
					defer lock.Lock()()
					s.Cost += d // 请求耗时
				})()

				arg := method(i)
				if reply, err := conn.Do(cmd, kit.Simple(arg)...); err != nil {
					// 请求失败
					nerr++
				} else {
					// 请求成功
					if nok++; check != nil {
						check(cmd, arg, reply)
					}
				}
			}()
		}
		return nil
	})
	return s, nil
}
