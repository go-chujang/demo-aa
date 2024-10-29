package id

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

func Version(uid string) (int, error) {
	u, err := uuid.Parse(uid)
	if err != nil {
		return 0, err
	}
	return int(u[6] >> 4), nil
}

// Versions UUID
const (
	Fastest       = 1
	FastestSorted = 6
	Randomly      = 4
	Sha1          = 5
	Md5           = 3
	TsRandom      = 7
)

func Uuid(version ...int) string {
	if version == nil || version[0] == 4 {
		return uuid.NewString() // version 4
	}

	var stringer fmt.Stringer
	switch version[0] {
	case 1:
		if v1, err := uuid.NewUUID(); err == nil {
			stringer = v1
		}
	case 3:
		v1, err := uuid.NewUUID()
		if err != nil {
			break
		}
		token := make([]byte, 4)
		rand.Read(token)
		stringer = uuid.NewMD5(v1, token)
	case 5:
		v1, err := uuid.NewUUID()
		if err != nil {
			break
		}
		token := make([]byte, 4)
		rand.Read(token)
		stringer = uuid.NewSHA1(v1, token)
	case 6:
		if v6, err := uuid.NewV6(); err == nil {
			stringer = v6
		}
	case 7:
		if v7, err := uuid.NewV7(); err == nil {
			stringer = v7
		}
	}
	if stringer == nil {
		return uuid.NewString()
	}
	return stringer.String()
}
