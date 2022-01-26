package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/jayanpraveen/tildly/datastore"
	cli "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
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

	if a.min >= a.max {
		a.min, a.max = a.etcd.GetNextRange()
	}

	kv := cli.NewKV(a.etcd.V3)
	key := a.etcd.RangeCountKey
	txn := kv.Txn(a.etcd.V3.Ctx())

	s, _ := concurrency.NewSession(a.etcd.V3)
	defer s.Close()

	l := concurrency.NewMutex(s, "/key-lock")

	{
		// Locking
		if err := l.Lock(a.etcd.CTX); err != nil {
			panic(err)
		}

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

		// Unlocking
		if err := l.Unlock(a.etcd.CTX); err != nil {
			panic(err)
		}
	}

	return a.min
}

func (a *atomicCounter) incCount() string {
	a.min++
	return strconv.Itoa(a.min)
}

func (a *atomicCounter) setModRevision(mr int64) {

	if a.modRev.ModRevision == 0 {
		res, _ := a.etcd.KV.Get(a.etcd.CTX, a.etcd.RangeCountKey)
		mr = res.Kvs[0].ModRevision
	}

	a.modRev.ModRevision = mr
}

func (a *atomicCounter) DisplayCurrentRange() string {
	return fmt.Sprintf("%d-%d", a.min, a.max)
}
