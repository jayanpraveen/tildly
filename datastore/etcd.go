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

	e.KV.Put(e.CTX, e.rangeLockSwitch, UNLOCK)

	go e.watchOnLock()

	e.KV.Put(e.CTX, e.RangeCountKey, "0")

	log.Println("V3 initzed required keys")

}

func (e *EtcdStore) watchOnLock() {

	watchChan := e.V3.Watch(e.CTX, e.rangeLockSwitch)

	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			fmt.Printf("[Event received]: %s executed on %q with value %q\n", event.Type, event.Kv.Key, event.Kv.Value)
			e.lock = string(event.Kv.Value)
		}
	}
}

func (e *EtcdStore) GetNextRange() (int, int) {
	return 0, 100
}

func (e *EtcdStore) Commit(head int) {
	e.KV.Put(e.CTX, e.RangeCountKey, fmt.Sprint(head))
}
