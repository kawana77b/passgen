package uid

import (
	"fmt"

	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Kind string

const (
	KindV4     Kind = "v4"
	KindV7     Kind = "v7"
	KindNanoID Kind = "nanoid"
)

func New(kind Kind) (string, error) {
	switch kind {
	case KindV4:
		id, err := uuid.NewRandom()
		return id.String(), err
	case KindV7:
		id, err := uuid.NewV7()
		return id.String(), err
	case KindNanoID:
		return gonanoid.New()
	default:
		return "", fmt.Errorf("unknown type %q: choose from v4, v7, nanoid", kind)
	}
}

func NewN(kind Kind, count int) ([]string, error) {
	results := make([]string, count)
	for i := range results {
		id, err := New(kind)
		if err != nil {
			return nil, err
		}
		results[i] = id
	}
	return results, nil
}
