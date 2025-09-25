package pkg

import (
	"net/url"
	"strconv"
)

func UpdatePageQuery(rawURL string, page int) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return rawURL
	}
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return u.String()
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 1
	}
	return i
}
