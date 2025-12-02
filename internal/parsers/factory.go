package parsers

import (
	"fmt"
	"log"
	"regexp"

	"github.com/crispyarty/novelparser/internal/parsers/novelbin"
)

type creator func() ParseHtml

func getDomain(url string) (string, error) {
	reg := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\n]+@)?(?:www\.)?([^:\/\n?]+)`)

	if matches := reg.FindStringSubmatch(url); len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("cant parse domain url: %v", url)

}

func ParserFactory(url string) creator {
	creators := map[string]creator{
		"novelbin.com": func() ParseHtml { return &novelbin.ParseHtmlNobelBin{} },
	}

	domain, err := getDomain(url)

	if err != nil {
		log.Fatalln(err)
	}

	return creators[domain]
}
