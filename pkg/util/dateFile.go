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
	"os"
	"strings"

	"golang.org/x/net/html"
)

/*
ParseGameFile will open a locally-saved HTML file that is formatted like the
    MLB date index file.  It will return a slice containing all of the game
    IDs contained in the file.
*/
func ParseGameFile(path string) ([]string, error) {
	var gids []string
	fileReader, err := os.Open(path)
	if err != nil {
		return gids, err
	}
	defer fileReader.Close()
	htmlTokenizer := html.NewTokenizer(fileReader)
	for {
		tt := htmlTokenizer.Next()
		if tt == html.ErrorToken {
			break
		}

		if tt == html.StartTagToken {
			t := htmlTokenizer.Token()

			isAnchor := t.Data == "a"
			if isAnchor {
				for _, a := range t.Attr {
					if a.Key == "href" {
						if strings.HasPrefix(a.Val, "gid_") {
							gids = append(gids, a.Val)
						}
						break
					}
				}
			}
		}
	}

	return gids, nil
}
