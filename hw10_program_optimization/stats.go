package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	reader := bufio.NewReader(r)
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return result, err
		}
		email := emailRegex.FindString(line)
		if strings.Contains(email, "."+domain) {
			result[strings.ToLower(strings.Split(email, "@")[1])]++
		}
		if err == io.EOF {
			break
		}
	}
	return result, nil
}
