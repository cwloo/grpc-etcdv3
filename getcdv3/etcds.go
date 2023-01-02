package getcdv3

import (
	"strings"
	"time"

	"github.com/cwloo/gonet/core/base/sys"
	"github.com/cwloo/gonet/logs"
	"github.com/cwloo/uploader/src/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	etcds = newEtcds()
)

// <summary>
// Etcds
// <summary>
type Etcds interface {
	Get() (cli Clientv3, e error)
	Put(cli Clientv3)
	Close(reset func(Clientv3))
	Update(etcdAddr string, reset func(Clientv3))
}

// <summary>
// etcds_
// <summary>
type etcds_ struct {
	etcdAddr string
	pool     sys.FreeValues
}

func newEtcds() Etcds {
	s := &etcds_{}
	s.pool = *sys.NewFreeValuesWith(s.new)
	return s
}

func (s *etcds_) Update(etcdAddr string, reset func(Clientv3)) {
	switch s.etcdAddr == etcdAddr {
	case true:
	default:
		switch etcdAddr {
		case "":
			return
		default:
		}
		switch s.etcdAddr {
		case "":
		default:
		}
		s.etcdAddr = etcdAddr
		s.pool.Update(func(value any, cb func(error, ...any)) (e error) {
			client := *value.(*Clientv3)
			client.Cancel()
			reset(*value.(*Clientv3))
			client.Close()
			c, err := s.new(cb)
			e = err
			switch err {
			case nil:
				*value.(*Clientv3) = c.(Clientv3)
			}
			return
		})
	}
}

func (s *etcds_) assertAddr() {
	if s.etcdAddr == "" {
		logs.Fatalf("error")
	}
}

func (s *etcds_) new(cb func(error, ...any)) (cli any, e error) {
	s.assertAddr()
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(s.etcdAddr, ","),
		DialTimeout: time.Duration(5) * time.Second,
		Username:    config.Config.Etcd.UserName,
		Password:    config.Config.Etcd.Password,
	})
	e = err
	switch err {
	case nil:
		cli = newClientv3(c)
		cb(err, s.etcdAddr)
	default:
		cb(err, s.etcdAddr)
	}
	return
}

func (s *etcds_) Get() (cli Clientv3, e error) {
	v, err := s.pool.Get()
	e = err
	switch err {
	case nil:
		cli = v.(Clientv3)
	default:
		logs.Errorf(err.Error())
	}
	return
}

func (s *etcds_) Put(cli Clientv3) {
	s.pool.Put(cli)
}

func (s *etcds_) Close(reset func(Clientv3)) {
	s.pool.Reset(func(value any) {
		value.(Clientv3).Cancel()
		reset(value.(Clientv3))
		value.(Clientv3).Close()
	}, false)
}