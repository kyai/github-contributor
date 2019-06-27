package creator

import (
	"encoding/json"
	"fmt"
)

type Contributor struct {
	Author struct {
		Login  string `json:"login"`
		Avatar string `json:"avatar"`
	} `json:"author"`
	Total int `json:"total"`
}

func getContributorsByRepo(repo string) (contributors []*Contributor, err error) {
	url := fmt.Sprintf("https://github.com/%s/graphs/contributors-data", repo)

	res, err := httpGetWithHeader(url, map[string]string{"accept": "application/json"})
	if err != nil {
		return
	}

	err = json.Unmarshal(res, &contributors)
	return
}
