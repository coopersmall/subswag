package utils

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

var (
	entropy = rand.New(rand.NewSource(rand.Int63()))
	digits  = 1000000
)

type ID int64

func (id ID) Validate() error {
	if id == 0 {
		return errors.New("ID cannot be empty")
	}
	return nil
}

func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func newId() ID {
	unix := time.Now().UnixMilli()
	ent := int64(rand.Intn(digits))
	return ID(unix*int64(digits) + ent)
}

var NewID = newId

func ParseID(s string) (ID, error) {
	parsedID, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return ID(parsedID), nil
}
