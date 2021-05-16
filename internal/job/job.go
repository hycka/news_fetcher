package job

import (
	"github.com/hycka/news_fetcher/internal/data"
)

func RefreshNews() error {
	fr := data.NewFetcherRepo()
	err := fr.UpdateNews()
	if err != nil {
		return err
	}
	return err
}
