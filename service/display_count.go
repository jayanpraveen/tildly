package service

import (
	"fmt"
	"net/http"

	"github.com/jayanpraveen/tildly/datastore"
)

func DisplayCount(etcd *datastore.EtcdStore) http.HandlerFunc {

	ac := NewAtomicCounter(etcd)

	return func(w http.ResponseWriter, r *http.Request) {

		c := ac.next()

		fmt.Fprint(w, c)

	}
}
