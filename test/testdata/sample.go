package testdata

import (
	"github.com/nduyhai/mapgen/test/bootiful"
	"github.com/nduyhai/mapgen/test/fancy"
	"time"
)

type User struct {
	ID           int
	UserName     string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func TimeToUnix(src time.Time) int64 {
	return src.Unix()
}

// +mapgen:mapper impl:userMapper target:user_mapper.go
type UserMapper interface {
	// +mapgen:mapping from:UserName to:Name
	// +mapgen:mapping from:CreatedAt to:CreatedAt using:TimeToUnix
	// +mapgen:mapping ignore:PasswordHash
	// +mapgen:mapping to:RetrievedAt using:bootiful.GetCurrentTime
	ToDTO(*User) *bootiful.UserDTO
}

// +mapgen:mapper impl:addressDtoMapper target:address_dto_mapper.go
type AddressDtoMapper interface {
	// +mapgen:mapping from:Street to:Street
	// +mapgen:mapping from:CreatedAt to:CreatedAt using:TimeToUnix
	// +mapgen:mapping ignore:ZipCode
	ToDTO(address *fancy.Address) *AddressDto
}
type AddressDto struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Created time.Time
}
