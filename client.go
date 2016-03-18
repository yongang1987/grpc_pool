package grpc_pool

import (
	log "code.google.com/p/log4go"
	"google.golang.org/grpc"
	"sync"
	"errors"
)

var (
	ClientClosedError = errors.New("client closed")
)

type Client struct {
	conn   *grpc.ClientConn
	C      interface{}
	mutex  sync.Mutex
	check  func(c interface{})
	closed bool
}

func newClient(conn *grpc.ClientConn, c interface{}, check func(c interface{})) *Client {
	client := &Client{
		conn:  conn,
		C:     c,
		check: check,
		mutex: sync.Mutex{},
	}
	go client.checkClose()
	return client
}

func (this *Client) checkClose() {
	//	stream, err := this.oc.Notify(context.Background(), &pbc.Request{})
	//	if err != nil {
	//		log.Error("client.fc.Notify() error(%v)", err)
	//		return
	//	}
	//	_, err = stream.Recv()
	//	this.Close()
	//	log.Info("notify close")
	//	return
	this.check(this.C)
	this.close()
}

func (this *Client) close() (err error) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if this.conn != nil {
		err = this.conn.Close()
		if err != nil {
			log.Error("client.conn.Close() error(%v)", err)
		}
	}
	this.closed = true
	return
}
