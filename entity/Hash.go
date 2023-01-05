package entity

import (
	"crypto/md5"
	"fmt"
)

type Hash struct {
	value string
}

// NewHash
func NewHash(word string) Hash {
	md5 := md5.Sum([]byte(word))
	return Hash{
		value: fmt.Sprintf("%x", md5),
	}
}

// StringToHash
func StringToHash(value string) Hash {
	return Hash{
		value: value,
	}
}

func (hash *Hash) String() string {
	return hash.value
}
