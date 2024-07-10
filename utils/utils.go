package utils

import "github.com/jaevor/go-nanoid"

func GenerateID() string {
	id, err := nanoid.Standard(8)

	if err != nil {
		panic(err)
	}

	return id()
}