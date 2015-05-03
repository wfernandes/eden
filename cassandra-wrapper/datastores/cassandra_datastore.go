package datastores

import (
	"github.com/gocql/gocql"
)

type CassandraDatastore struct {
	Session *gocql.Session
}

type Metric struct {
	Name     string
	ReqCount uint32
	ErrCount uint32
}

func NewCassandra(keyspace string, clusterHosts []string) *CassandraDatastore {
	cluster := gocql.NewCluster(clusterHosts)
	cluster.Keyspace = keyspace

	s, err := cluster.CreateSession()
	if err != nil {
		fmt.Printf("error creating cassandra session: %s", err.Error())
		panic(err)
	}

	return &CassandraDatastore{
		Session: s,
	}
}

func (c *CassandraDatastore) InsertMetric(metric *Metric) error {

	query := c.Session.Query(`INSERT INTO metrics (name, req_count, err_count) VALUES (?, ?, ?)`, metric.Name, metric.ReqCount, metric.ErrCount)

	err := query.Exec()
	if err != nil {
		fmt.Printf("error inserting metric: %s", err.Error())
		return err
	}
	return nil
}

func (c *CassandraDatastore) Close() {
	c.Session.Close()
}
