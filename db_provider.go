package main

import "sync"

var pool = map[string]SqlUserClient{}
var lock sync.Mutex

func GetSqlClient(connString, dbId string) (SqlUserClient, error) {
	lock.Lock()
	defer lock.Unlock()

	if conn, ok := pool[connString]; ok {
		return conn, nil
	}

	conn, err := CreateSqlClient(connString, dbId)
	if err != nil {
		return nil, err
	}

	pool[connString] = conn
	return conn, nil
}

func DisposeConnections() error {
	lock.Lock()
	defer lock.Unlock()

	for connString, conn := range pool {
		delete(pool, connString)
		if err := conn.Close(); err != nil {
			return err
		}
	}

	return nil
}
