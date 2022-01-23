package datastore

import (
	"flag"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdCoordinator interface {
	getNextRange() (start int, end int, err error)
}

func etcdConn() {

	etcdHost1 := flag.String("infra1", "http://127.0.0.1:2379", "etcd host-1")
	etcdHost2 := flag.String("infra2", "http://127.0.0.1:22379", "etcd host-2")
	etcdHost3 := flag.String("infra3", "http://127.0.0.1:32379", "etcd host-3")

	flag.Parse()

	fmt.Println("connecting to etcd - "+*etcdHost1, *etcdHost2, *etcdHost3)

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdHost1, *etcdHost2, *etcdHost3},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer etcd.Close()

}
