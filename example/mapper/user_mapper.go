package mapper

import "time"

// +mapgen:mapper impl:userMapper
type UserMapper interface {
	// +mapgen:mapping from:UserName to:Name
	// +mapgen:mapping from:CreatedAt to:CreatedAt using:TimeToUnix
	// +mapgen:mapping ignore:PasswordHash
	ToDTO(*User) *UserDTO

	// +mapgen:mapping from:Name to:UserName
	FromDTO(*UserDTO) *User
}

type User struct {
	Name         string
	CreatedAt    time.Time
	PasswordHash string
}
type UserDTO struct {
	UserName  string
	CreatedAt time.Time
}
