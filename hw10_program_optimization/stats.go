package hw10programoptimization

import (
	"errors"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func getUsers(r io.Reader) (chan User, chan error) {
	dataChan := make(chan User, 100_0)
	errChan := make(chan error)

	dec := jsoniter.NewDecoder(r)

	go func() {
		for {
			var user User
			err := dec.Decode(&user)
			if err != nil && !errors.Is(err, io.EOF) {
				errChan <- err
				break
			}

			dataChan <- user

			if errors.Is(err, io.EOF) {
				break
			}
		}
		close(dataChan)
	}()

	return dataChan, errChan
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	dataChan, errChan := getUsers(r)

	result := make(DomainStat)

	for {
		select {
		case user, ok := <-dataChan:

			if !ok {
				return result, nil
			}

			matched := strings.HasSuffix(user.Email, "."+domain)

			if matched {
				baseRaw := strings.Split(user.Email, "@")
				base := strings.ToLower(baseRaw[1])
				result[base]++
			}

		case err := <-errChan:
			return nil, err
		}
	}
}
