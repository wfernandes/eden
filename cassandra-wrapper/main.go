package main

import (
	"github.com/wfernandes/eden/cassandra-wrapper/datastores"
)

func main() {

	keyspace := "mykeyspace"
	hosts := []string{"127.0.0.1"}

	cassandra := datastores.NewCassandra(keyspace, hosts)
	defer cassandra.Close()
	metric := datastores.Metric{
		Name:     "mymetric",
		ReqCount: 111,
		ErrCount: 10,
	}

	err := cassandra.InsertMetric(&metric)

	if err != nil {
		panic(err)
	}
}
