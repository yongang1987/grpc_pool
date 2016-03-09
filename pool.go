package grpc_pool

import (
	log "code.google.com/p/log4go"
	"github.com/yongang1987/pool"
	"google.golang.org/grpc"
	"math/rand"
	"time"
)

type GrpcPool struct {
	p         *pool.Pool
	MaxActive int
	MaxIdle   int
	Timeout   time.Duration
	check     func(interface{})
	Addr      []string
}

func NewGrpcPool(addr []string, maxActive int, maxIdle int, timeout time.Duration, getClient func(conn *grpc.ClientConn) interface{}, check func(interface{})) *GrpcPool {
	return &GrpcPool{
		p: &pool.Pool{
			Dial: func() (interface{}, error) {
				conn, err := grpc.Dial(addr[rand.Intn(len(addr))], grpc.WithInsecure())
				if err != nil {
					log.Error("did not connect: %v", err)
					return nil, err
				}
				return newClient(conn, getClient(conn), check), nil
			},
			Close: func(v interface{}) error {
				return v.(*Client).close()
			},
			MaxActive:   maxActive,
			MaxIdle:     maxIdle,
			IdleTimeout: timeout,
		},
	}
}

func (gp *GrpcPool) Close() {
	if gp.p != nil {
		gp.p.Release()
	}
}

func (gp *GrpcPool) Get() (c *Client, err error) {
	v, err := gp.p.Get()
	if err != nil {
		return
	}
	c = v.(*Client)
	return
}

func (gp *GrpcPool) Put(c *Client, b bool) {
	gp.p.Put(c, b)
}
