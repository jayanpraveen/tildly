package service

import (
	"log"
	"strconv"

	"github.com/jayanpraveen/tildly/datastore"
	cli "go.etcd.io/etcd/client/v3"
)

type atomicCounter struct {
	min    int
	max    int
	etcd   *datastore.EtcdStore
	modRev PrevModRevision
}

func NewAtomicCounter(etcd *datastore.EtcdStore) *atomicCounter {

	min, max := etcd.GetNextRange()

	return &atomicCounter{
		min:  min,
		max:  max,
		etcd: etcd,
	}
}

type PrevModRevision struct {
	ModRevision int64
}

func (a *atomicCounter) next() int {

	kv := cli.NewKV(a.etcd.V3)
	key := a.etcd.RangeCountKey

	txn := kv.Txn(a.etcd.V3.Ctx())

	a.setModRevision(a.modRev.ModRevision)

	res, err := txn.If(cli.Compare(cli.ModRevision(key), "=", a.modRev.ModRevision)).
		Then(
			cli.OpPut(key, a.incCount()),
			cli.OpGet(key),
		).
		Commit()

	if err != nil {
		panic(err)
	}
	if !res.Succeeded {
		log.Println("A newer ModRevision exists.")
		return 0
	}

	a.setModRevision(res.Responses[1].GetResponseRange().Kvs[0].ModRevision)

	return a.min
}

func (a *atomicCounter) incCount() string {
	a.min++
	return strconv.Itoa(a.min)
}

func (a *atomicCounter) setModRevision(mr int64) {

	// this sets the lastest ModRevision,
	// reason: useful when working locally where you restart the etcd server muliple times
	// with old "/range/counter" still existing.
	if a.modRev.ModRevision == 0 {
		res, _ := a.etcd.KV.Get(a.etcd.CTX, a.etcd.RangeCountKey)
		mr = res.Kvs[0].ModRevision
	}

	a.modRev.ModRevision = mr
}
