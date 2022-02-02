package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
	//add  a new ErrInvaliCredentials erros. We'll user this lates is a user
	//tires to login with an incorrect email address or password.

	ErrInvalidCredentials = errors.New("models:invalid credentials")

	//add a new ErroDuplicateEmail error. We'll use this lates if a user
	//treies a to signup with an email address that's already in use

	ErrDuplicateEmail = errors.New("midels: duplicate email")
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Create  time.Time
	Expires time.Time
}

//Definea enw user type. Notice how the field names and types align
// with the colimns in the database 'users' table?
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Create         time.Time
	Active         bool
}
