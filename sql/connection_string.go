package sql

import (
	"fmt"
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

func formatBool(x bool) string {
	if x {
		return "True"
	}
	return "False"
}
