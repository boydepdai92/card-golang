package common

var CurrentUser = "user"

type Requester interface {
	GetId() int
	GetUsername() string
}
