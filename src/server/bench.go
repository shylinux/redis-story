package server

import (
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"shylinux.com/x/ice"
	"shylinux.com/x/redis-story/src/client"
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

var trans = map[string]func(i int64) []string{
	"GET": func(i int64) []string { return []string{fmt.Sprintf("hi%d", i)} },
	"SET": func(i int64) []string { return []string{fmt.Sprintf("hi%d", i), "hello"} },
}
var ErrNotFoundCmd = errors.New("not found cmd")

func Bench(nconn, nreq int64, hosts []string, cmds []string, check func(cmd string, arg []string, res ice.Any)) (*Stat, error) {
	s := &Stat{BeginTime: time.Now()}
	defer func() { // 请求统计
		if s.EndTime = time.Now(); s.BeginTime != s.EndTime {
			d := float64(s.EndTime.Sub(s.BeginTime)) / float64(time.Second)

			s.Up = float64(s.NWrite) / d
			s.Down = float64(s.NRead) / d
			s.QPS = float64(s.NReq) / d
			s.AVG = s.EndTime.Sub(s.BeginTime) / time.Duration(nreq)
		}
	}()

	// 连接池
	rp := client.NewRedisPool(hosts[0], "")

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

		cmd := strings.ToUpper(cmds[0])
		method := trans[cmd]
		if method == nil {
			return ErrNotFoundCmd
		}

		for i := int64(0); i < nreq; i++ {
			func() {
				defer logs.CostTime(func(d time.Duration) {
					defer lock.Lock()()
					s.Cost += d // 请求耗时
				})()

				arg := method(i)
				if reply, err := conn.Do(cmd, arg...); err != nil {
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
