package server

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gomodule/redigo/redis"
	log "github.com/shylinux/toolkits/logs"
	"github.com/shylinux/toolkits/task"
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
	mu   sync.Mutex

	Up   float64
	Down float64
	QPS  float64
	AVG  time.Duration
}

func _GET(i int64) []interface{} { return []interface{}{fmt.Sprintf("hi%d", i)} }
func _SET(i int64) []interface{} { return []interface{}{fmt.Sprintf("hi%d", i), "hello"} }

var trans = map[string]func(i int64) []interface{}{"GET": _GET, "SET": _SET}

func Bench(nconn, nreq int64, hosts []string, cmds []string, check func(cmd string, arg []interface{}, res interface{})) (*Stat, error) {
	// 请求统计
	s := &Stat{BeginTime: time.Now()}
	defer func() {
		if s.EndTime = time.Now(); s.BeginTime != s.EndTime {
			d := float64(s.EndTime.Sub(s.BeginTime)) / float64(time.Second)

			s.QPS = float64(s.NReq) / d
			s.AVG = s.EndTime.Sub(s.BeginTime) / time.Duration(nreq)
			s.Down = float64(s.NRead) / d
			s.Up = float64(s.NWrite) / d
		}
	}()

	// 连接池
	rp := redis.NewPool(func() (redis.Conn, error) { return redis.Dial("tcp", hosts[0]) }, 10)

	// 协程池
	list := []interface{}{}
	for i := int64(0); i < nconn; i++ {
		list = append(list, i)
	}

	task.Wait(list, func(task *task.Task, lock *task.Lock) error {
		// 请求汇总
		var nerr, nok int64
		defer func() {
			atomic.AddInt64(&s.NReq, nreq)
			atomic.AddInt64(&s.NErr, nerr)
			atomic.AddInt64(&s.NOK, nok)
		}()

		conn := rp.Get()
		defer conn.Close()

		cmd := strings.ToUpper(cmds[0])
		method := trans[cmd]
		if method == nil {
			log.Warn("method %v not found", cmd)
			return errors.New("not found")
		}

		for i := int64(0); i < nreq; i++ {
			func() {
				// 请求耗时
				begin := time.Now()
				defer func() {
					d := time.Now().Sub(begin)

					s.mu.Lock()
					defer s.mu.Unlock()
					s.Cost += d
				}()

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
