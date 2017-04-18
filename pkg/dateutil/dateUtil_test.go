package dateutil

import (
	"net/url"
	"testing"
	"time"
)

func helperForURL(rawURL string) *url.URL {
	retURL, err := url.Parse(rawURL)
	if err != nil {
		return nil
	}
	return retURL
}

func TestDateToURL(t *testing.T) {
	var dateTest = []struct {
		Date        time.Time
		RootURL     string
		ExpectedURL *url.URL
	}{
		{time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC), "http://www.test.com/components/game/mlb", helperForURL("http://www.test.com/components/game/mlb/year_2017/month_04/day_01/")},
	}

	for _, ex := range dateTest {
		SetRoot(ex.RootURL, "")
		retURL, err := DateToURL(ex.Date)
		if err != nil {
			t.Errorf("Unable to convert date to URL")
		}
		retEscapedPath := retURL.EscapedPath()
		expectedEscapedPath := ex.ExpectedURL.EscapedPath()
		if retEscapedPath != expectedEscapedPath {
			t.Errorf("URLs do not match -> %s vs %s", retEscapedPath, expectedEscapedPath)
		}
	}
}
