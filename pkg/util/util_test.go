/*
	Copyright 2017 Brian Bauer

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package util

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
		retURL, err := DateToURLNoSideEffects(ex.Date, ex.RootURL)
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

func TestURLToFSPath(t *testing.T) {
	var pathTest = []struct {
		URL          *url.URL
		RootFS       string
		ExpectedPath string
	}{
		{helperForURL("http://www.test.com/components/game/mlb/year_2017/month_04/day_01/"), "/root", "/root/year_2017/month_04/day_01/index.html"},
		{helperForURL("http://www.test.com/components/game/mlb/year_2017/month_04/day_01/index.html"), "/root", "/root/year_2017/month_04/day_01/index.html"},
		{helperForURL("http://www.test.com/components/game/mlb/year_2017/month_04/day_01/test.xml"), "/root", "/root/year_2017/month_04/day_01/test.xml"},
	}

	for _, ex := range pathTest {
		retPath, err := URLToFSPathNoSideEffects(ex.URL, ex.RootFS)
		if err != nil {
			t.Errorf("Unable to convert URL to path")
		}
		if retPath != ex.ExpectedPath {
			t.Errorf("Paths do not match -> %s vs %s", retPath, ex.ExpectedPath)
		}
	}
}
