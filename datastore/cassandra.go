package datastore

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

func DialCassandra() *gocql.Session {
	cluster := gocql.NewCluster("cassandra0", "cassandra1", "cassandra2")
	keyspace := "tildly"
	table := `CREATE TABLE IF NOT EXISTS tildly.url (
		hash text,
		LongUrl text,
		createdat bigint,
		exipireat bigint,
		PRIMARY KEY ((hash), createdat)
	  )
	  WITH CLUSTERING ORDER BY (createdat DESC);`

	s := createKeyspace(cluster, keyspace)
	createTable(s, table)
	return s
}

func createSession(cluster *gocql.ClusterConfig) *gocql.Session {
	for {
		session, err := cluster.CreateSession()
		if err == nil {
			log.Println("Connected to cassandra")
			return session
		}
		log.Print("csdra err: reconnecting in 3 seconds")
		time.Sleep(3 * time.Second)
	}
}

func createKeyspace(cluster *gocql.ClusterConfig, keyspace string) *gocql.Session {
	flagRF := flag.Int("rf", 3, "replication factor for keyspace")
	flag.Parse()

	c := *cluster
	c.Keyspace = "system"
	c.Timeout = 30 * time.Second
	c.Consistency = gocql.Quorum

	session := createSession(cluster)

	err := createTable(session, fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, keyspace, *flagRF))

	if err != nil {
		panic(fmt.Sprintf("unable to create keyspace: %v", err))
	}
	return session
}

func createTable(s *gocql.Session, table string) error {
	if err := s.Query(table).RetryPolicy(&gocql.SimpleRetryPolicy{}).Exec(); err != nil {
		log.Printf("error creating table table=%q err=%v\n", table, err)
		return err
	}
	return nil
}
