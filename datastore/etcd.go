package datastore

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Coordinator interface {
	GetNextRange() (start int, end int)
	Commit(head int)
}

type EtcdStore struct {
	KV  clientv3.KV
	CTX context.Context
	V3  *clientv3.Client
	COD Coordinator

	RangeCountKey string
}

func NewEtcd() *EtcdStore {

	etcdHost1 := flag.String("infra1", "http://etcd0:2379", "V3 host-1")
	etcdHost2 := flag.String("infra2", "http://etcd1:22379", "V3 host-1")
	etcdHost3 := flag.String("infra3", "http://etcd2:32379", "V3 host-1")
	flag.Parse()
	log.Println("connecting to V3 - ", *etcdHost1, *etcdHost2, *etcdHost3)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdHost1, *etcdHost2, *etcdHost3},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	log.Println("connected to V3 - ", *etcdHost1, *etcdHost2, *etcdHost3)

	es := EtcdStore{
		KV:  clientv3.NewKV(cli),
		CTX: context.Background(),
		V3:  cli,
	}

	return &es
}

func (e *EtcdStore) GetNextRange() (start int, end int) {

	for i := 0; true; i++ {
		key := "/range/range"
		iKey := fmt.Sprintf("%s-%d", key, i)
		r, _ := e.KV.Get(e.CTX, iKey)
		if r.Count == 0 {
			e.KV.Put(e.CTX, iKey, "locked")
			e.RangeCountKey = iKey + "/counter"
			return nextMillionRange(i)
		}
	}
	return
}

func nextMillionRange(r int) (s int, e int) {
	s = ((r * 1000000) + 1)
	e = ((s + 1000000) - 1)
	return s, e
}
