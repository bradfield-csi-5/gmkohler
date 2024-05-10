package pkg

// Role is a first pass at making the main functions dumber, might be better located somewhere else, and might break
// down once there are many replicas etc.
type Role int

const (
	RolePrimary = iota + 1
	RoleReplica
)
