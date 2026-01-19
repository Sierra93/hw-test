package hw10programoptimization

import (
	"bufio"
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

func getUsers(r io.Reader) (users, error) {
	var result users
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		line := scanner.Bytes() // Получаем []byte без копирования
		var user User
		if err := json.Unmarshal(line, &user); err != nil {
			return result, err
		}
		if i >= len(result) {
			break
		}
		result[i] = user
		i++
	}
	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainLower := strings.ToLower(domain)

	for _, user := range u {
		email := user.Email
		emailLower := strings.ToLower(email)
		atIdx := strings.LastIndex(emailLower, "@")
		if atIdx == -1 {
			continue
		}
		domainPart := emailLower[atIdx+1:]
		if domainPart == domainLower || strings.HasSuffix(domainPart, "."+domainLower) {
			result[domainPart]++
		}
	}
	return result, nil
}
