package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
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

func IsPassword(s, specialSigns string) bool {
	pattern := fmt.Sprintf("^[a-zA-Z0-9%s]+$", specialSigns)
	match, _ := regexp.Match(pattern, []byte(s))
	return match
}

func IsOneWord(s string) bool {
	return len(strings.Split(s, " ")) == 1
}

func IsLetter(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return len(s) > 0
}

// IsOneWordLetter проверяет, состоит ли строка из одного слова, состоящего только из букв.
// Поддерживает не-ASCII символы. !!!Для пустой строки вернет true.
func IsOneWordLetter(s string) bool {
	if len(s) == 0 {
		return true
	}

	if len(strings.Split(s, " ")) > 1 {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
