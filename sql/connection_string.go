package sql

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type ConnectionString struct {
	ServerAddress          string
	Port                   int
	Database               string
	Username               string
	Password               string
	ConnectionTimeout      int
	MaxPoolSize            int
	TrustServerCertificate bool
	PersistSecurityInfo    bool
	Encrypt                bool
}

func (s *ConnectionString) String() (string, error) {
	var sb strings.Builder

	if s.ServerAddress == "" {
		return "", fmt.Errorf("missing ServerAddress")
	}

	if s.Database == "" {
		return "", fmt.Errorf("missing Database")
	}

	fmt.Fprintf(&sb, "Server=%s;", s.ServerAddress)
	if s.Port != 0 {
		fmt.Fprintf(&sb, "Port=%d;", s.Port)
	}

	fmt.Fprintf(&sb, "Database=%s;", s.Database)

	if s.Username != "" {
		fmt.Fprintf(&sb, "User ID=%s;", s.Username)
	}
	if s.Password != "" {
		fmt.Fprintf(&sb, "Password=%s;", s.Password)
	}

	fmt.Fprintf(&sb, "Connection Timeout=%d;", s.ConnectionTimeout)
	fmt.Fprintf(&sb, "Persist Security Info=%s;", formatBool(s.PersistSecurityInfo))
	fmt.Fprintf(&sb, "Trust Server Certificate=%s;", formatBool(s.TrustServerCertificate))
	fmt.Fprintf(&sb, "Max Pool Size=%d;", s.MaxPoolSize)

	return sb.String(), nil
}

func ParseConnectionString(connectionString string) (*ConnectionString, error) {
	// TODO: this is not perfect

	log.Printf("[ERROR] %s", connectionString)

	result := &ConnectionString{}

	split := strings.Split(connectionString, ";")

	for _, pairRaw := range split {
		pair := strings.SplitN(pairRaw, "=", 2)
		if len(pair) != 2 {
			return nil, fmt.Errorf("unexpected value %s", pairRaw)
		}

		k, v := strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1])

		var err error

		if k == "Server" {
			var server string
			var port int
			server, port, err = getServerAndPort(v)
			result.ServerAddress = server
			result.Port = port
		} else if k == "Database" {
			result.Database = v
		} else if k == "Username" {
			result.Username = v
		} else if k == "Password" {
			result.Password = v
		} else if k == "ConnectionTimeout" {
			var timeout int
			timeout, err = strconv.Atoi(v)
			result.ConnectionTimeout = timeout
		} else if k == "MaxPoolSize" {
			var pool int
			pool, err = strconv.Atoi(v)
			result.MaxPoolSize = pool
		} else if k == "TrustServerCertificate" {
			var trust bool
			trust, err = parseBool(v)
			result.TrustServerCertificate = trust
		} else if k == "PersistSecurityInfo" {
			var persist bool
			persist, err = parseBool(v)
			result.PersistSecurityInfo = persist
		} else if k == "Encrypt" {
			var encrypt bool
			encrypt, err = parseBool(v)
			result.Encrypt = encrypt
		}

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func getServerAndPort(address string) (string, int, error) {
	split := strings.SplitN(address, ",", 2)
	if len(split) == 1 {
		return split[0], 0, nil
	} else if len(split) == 2 {
		port, err := strconv.Atoi(split[1])
		return split[0], port, err
	} else {
		return "", 0, fmt.Errorf("invalid server address: %s", address)
	}
}

func parseBool(x string) (bool, error) {
	if x == "True" {
		return true, nil
	} else if x == "False" {
		return false, nil
	} else {
		return false, fmt.Errorf("unexpected boolean value %s", x)
	}
}

func formatBool(x bool) string {
	if x {
		return "True"
	}
	return "False"
}
