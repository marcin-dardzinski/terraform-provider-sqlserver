package main

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

func (s *ConnectionString) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Server=tcp:%s,%d;", s.ServerAddress, s.Port)
	fmt.Fprintf(&sb, "Database=%s;", s.Database)
	fmt.Fprintf(&sb, "User ID=%s;Password=%s;", s.Username, s.Password)
	fmt.Fprintf(&sb, "Connection Timeout=%d;", s.ConnectionTimeout)
	fmt.Fprintf(&sb, "Persist Security Info=%s;", formatBool(s.PersistSecurityInfo))
	fmt.Fprintf(&sb, "Trust Server Certificate=%s;", formatBool(s.TrustServerCertificate))
	fmt.Fprintf(&sb, "Max Pool Size=%d;", s.MaxPoolSize)

	return sb.String()
}

func formatBool(x bool) string {
	if x {
		return "True"
	}
	return "False"
}
