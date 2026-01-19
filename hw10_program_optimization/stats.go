package hw10programoptimization

import (
	"encoding/json"
	"fmt"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainLower := strings.ToLower(domain)

	for _, user := range u {
		emailLower := strings.ToLower(user.Email)
		atIdx := strings.LastIndex(emailLower, "@")
		if atIdx == -1 {
			continue
		}
		domainPart := emailLower[atIdx+1:]
		if strings.HasSuffix(domainPart, domainLower) {
			if len(domainPart) == len(domainLower) || domainPart[len(domainPart)-len(domainLower)-1] == '.' {
				result[domainPart]++
			}
		}
	}
	return result, nil
}
