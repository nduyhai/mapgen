package bootiful

import "time"

type UserDTO struct {
	ID          int
	Name        string
	FirstName   string
	LastName    string
	Email       string
	CreatedAt   int64
	RetrievedAt int64
}

func GetCurrentTime() int64 {
	return time.Now().Unix()
}
