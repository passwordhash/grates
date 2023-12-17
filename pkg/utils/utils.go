package utils

import (
	"fmt"
	"math/rand"
	"time"
)

type Date struct {
	time.Time
}

func (d Date) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("\"2006-01-02\"", string(b))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d Date) String() string {
	return d.Time.Format("2006-01-02")
}

func GenerateHash(length int) (string, error) {
	b := make([]byte, length)

	src := rand.NewSource(time.Now().Unix())
	r := rand.New(src)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func RandStringBytesRmndr(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
