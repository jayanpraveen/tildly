package datastore

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	LOCK   = "ON"
	UNLOCK = "OFF"
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

	rangeKey        string
	lock            string
	rangeLockSwitch string
}

func (e *EtcdStore) init() {
	e.rangeKey = "/range"
	e.RangeCountKey = e.rangeKey + "/counter"
	e.rangeLockSwitch = e.rangeKey + "/lock"
}

func NewEtcd() *EtcdStore {

	etcdHost1 := flag.String("infra1", "http://0.0.0.0:2379", "V3 host-1")
	flag.Parse()
	fmt.Println("connecting to V3 - " + *etcdHost1)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdHost1},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	es := EtcdStore{
		KV:  clientv3.NewKV(cli),
		CTX: context.Background(),
		V3:  cli,
	}

	es.init()
	es.setUpKeys()

	return &es
}

func (e *EtcdStore) setUpKeys() {

	e.KV.Put(e.CTX, e.RangeCountKey, "0")
	// e.KV.Put(e.CTX, "/range/range0", "0")

	log.Println("V3 initzed required keys")

}

func (e *EtcdStore) GetNextRange() (start int, end int) {

	key := "/range/range"

	for i := 0; true; i++ {

		r, _ := e.KV.Get(e.CTX, fmt.Sprintf("%s%d", key, i))
		log.Println(r.Count)

		iKey := fmt.Sprintf("%s%d", key, i)

		if r.Count == 0 {
			e.KV.Put(e.CTX, iKey, "locked")
			e.RangeCountKey = iKey + "/counter"
			return next100Range(i)
		}

	}

	return

}

func next100Range(r int) (int, int) {

	s := ((r * 100) + 1)
	e := ((s + 100) - 1)

	return s, e
}

func (e *EtcdStore) Commit(head int) {
	e.KV.Put(e.CTX, e.RangeCountKey, fmt.Sprint(head))
}
