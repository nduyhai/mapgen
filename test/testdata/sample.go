package testdata

import (
	"github.com/nduyhai/mapgen/test/fancy"
	"time"
)

type User struct {
	ID        int
	UserName  string
	FirstName string
	LastName  string
	Email     string
	CreatedAt time.Time
}

type UserDTO struct {
	ID        int
	Name      string
	FirstName string
	LastName  string
	Email     string
	CreatedAt int64
}

func TimeToUnix(src time.Time) int64 {
	return src.Unix()
}

// +mapgen:mapper impl:userMapper target:user_mapper.go
type UserMapper interface {
	// +mapgen:mapping from:UserName to:Name
	// +mapgen:mapping from:CreatedAt to:CreatedAt using:TimeToUnix
	// +mapgen:mapping ignore:PasswordHash
	ToDTO(*User) *UserDTO
}

// +mapgen:mapper impl:addressDtoMapper target:address_dto_mapper.go
type AddressDtoMapper interface {
	// +mapgen:mapping from:UserName to:Name
	// +mapgen:mapping from:CreatedAt to:CreatedAt using:TimeToUnix
	// +mapgen:mapping ignore:PasswordHash
	ToDTO(address *fancy.Address) *AddressDto
}
type AddressDto struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Created time.Time
}
