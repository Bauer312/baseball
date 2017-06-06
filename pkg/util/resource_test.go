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

func TestDateResource(t *testing.T) {
	var dateTest = []struct {
		Date         time.Time
		RootURL      string
		RootFS       string
		ExpectedURL  *url.URL
		ExpectedPath string
	}{
		{
			time.Date(2017, time.April, 1, 5, 0, 0, 0, time.UTC),
			"http://www.test.com/components/game/mlb",
			"/example/baseball",
			helperForURL("http://www.test.com/components/game/mlb/year_2017/month_04/day_01/"),
			"/example/baseball/year_2017/month_04/day_01/index.html",
		},
	}

	for _, ex := range dateTest {
		var tstResource Resource
		err := tstResource.Roots(ex.RootURL, ex.RootFS)
		if err != nil {
			t.Error("Unable to set root elements")
		}
		tDefs, err := tstResource.Date(ex.Date)
		if err != nil {
			t.Error("Unable to get transfer definitions for date")
		}
		if len(tDefs) != 1 {
			t.Errorf("Expecting 1 transfer definition but instead got %d", len(tDefs))
		}
		if tDefs[0].Source.EscapedPath() != ex.ExpectedURL.EscapedPath() {
			t.Errorf("Expecting URL %s but got %s instead", ex.ExpectedURL, tDefs[0].Source)
		}
		if tDefs[0].Target != ex.ExpectedPath {
			t.Errorf("Expecting Path %s but got %s instead", ex.ExpectedPath, tDefs[0].Target)
		}
	}
}
