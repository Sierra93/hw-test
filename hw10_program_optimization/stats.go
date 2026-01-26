package hw10programoptimization

import (
	"bufio"
	"bytes"
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
