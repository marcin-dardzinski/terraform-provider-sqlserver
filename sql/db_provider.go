package sql

import "sync"

var pool = map[string]*SqlClient{}
var lock sync.Mutex

func CreatePooledSqlClient(config SqlClientConfig) (*SqlClient, error) {
	lock.Lock()
	defer lock.Unlock()

	if conn, ok := pool[config.ConnectionString]; ok {
		return conn, nil
	}

	conn, err := CreateSqlClient(config)
	if err != nil {
		return nil, err
	}

	pool[config.ConnectionString] = conn
	return conn, nil
}

func DisposeConnections() error {
	lock.Lock()
	defer lock.Unlock()

	for connString, conn := range pool {
		delete(pool, connString)
		if err := conn.Db.Close(); err != nil {
			return err
		}
	}

	return nil
}
