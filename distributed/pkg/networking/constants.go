package networking

import (
	"distributed/pkg"
	"fmt"
)

const Unix = "unix"
const Port = 8907

const (
	primarySocketPath = "/tmp/distkv.sock"
	replicaSocketPath = "/tmp/distkv-replica.sock"
)

func SocketPath(role pkg.Role) (string, error) {
	switch role {
	case pkg.RolePrimary:
		return primarySocketPath, nil
	case pkg.RoleReplica:
		return replicaSocketPath, nil
	default:
		return "", fmt.Errorf("unrecognized Role %v", role)
	}
}
