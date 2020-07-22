package sql

import (
	"testing"
)

func TestParseDatabaseId(t *testing.T) {
	connString := "Server=localhost;Database=Db1;User Id=sa;Password=Passwd1!;"
	id := parseDatabaseId(connString)
	expected := "localhost/Db1"

	if id != expected {
		t.Errorf("Expected %s, got %s", expected, id)
	}
}
