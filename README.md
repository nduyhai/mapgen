# mapgen

[![Go](https://img.shields.io/badge/go-1.24+-blue)](https://go.dev/)
[![License](https://img.shields.io/github/license/nduyhai/mapgen)](LICENSE)

A GitHub template repository for bootstrapping a new Go project with a clean, idiomatic layout.

## Features


## Getting Started

###
```shell
go run cmd/mapgen/*.go
```

### 
```go

// +mapgen:mapper impl:userMapper
type UserMapper interface {
	
    // +mapgen:mapping from:UserName to:Name
    // +mapgen:mapping from:CreatedAt to:CreatedAt using:TimeToUnix
    // +mapgen:mapping ignore:PasswordHash
    ToDTO(*proto.User) *dto.UserDTO
}

```
