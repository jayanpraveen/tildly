package datastore

import (
	"log"

	"github.com/gocql/gocql"
)

func DialCassandra() *gocql.Session {
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Keyspace = "tildly"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to cassandra")
	return session
}
