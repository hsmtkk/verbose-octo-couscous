package html

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func SelectDataSrc(html string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML; %w", err)
	}
	srcs := []string{}
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		dataSrc, ok := s.Attr("data-src")
		if ok {
			srcs = append(srcs, dataSrc)
		}
	})
	return srcs, nil
}
