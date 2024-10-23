package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type UserEmail struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainSuffix := "." + domain

	decoder := json.NewDecoder(r)
	for {
		var user UserEmail
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domainSuffix) {
			atIndex := strings.Index(user.Email, "@")
			if atIndex != -1 {
				domainPart := strings.ToLower(user.Email[atIndex+1:])
				result[domainPart]++
			}
		}
	}

	return result, nil
}
