package hw10programoptimization

import (
	"bufio"
	"bytes"
	"encoding/json"
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
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	domainLower := strings.ToLower(domain)
	domainBytes := []byte(domainLower)
	dotDomainBytes := []byte("." + domainLower)
	emailKey := []byte(`"Email":"`)

	for scanner.Scan() {
		line := scanner.Bytes()
		idx := bytes.Index(line, emailKey)
		if idx == -1 {
			continue
		}

		start := idx + len(emailKey)
		end := bytes.IndexByte(line[start:], '"')
		if end == -1 {
			continue
		}

		email := line[start : start+end]
		atIdx := bytes.LastIndexByte(email, '@')
		if atIdx == -1 {
			continue
		}

		domainPart := email[atIdx+1:]
		domainPartLower := bytes.ToLower(domainPart)

		if bytes.Equal(domainPartLower, domainBytes) || bytes.HasSuffix(domainPartLower, dotDomainBytes) {
			result[string(domainPartLower)]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
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
