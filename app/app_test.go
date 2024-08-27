package app

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	t.Parallel()

	type Request struct {
		path string // must contain a leading slash
		resp []byte
	}

	type testCase struct {
		name string

		requests []Request

		expectedLinks []string
	}

	testCases := []testCase{
		{
			name: "onePage",
			requests: []Request{
				{
					path: "/",
					resp: []byte(`
				<html>
					<body>
						<a href="/link1">Link 1</a>
						<a href="/link2">Link 2</a>
					</body>
				</html>
`),
				},
			},
			expectedLinks: []string{"/link1", "/link2"},
		},
		{
			name: "twoPages",
			requests: []Request{
				{
					path: "/",
					resp: []byte(`
				<html>
					<body>
						<a href="/basic1">Link 1</a>
						<a href="/basic2">Link 2</a>
					</body>
				</html>
`),
				},
				{
					path: "/basic1",
					resp: []byte(`
				<html>
					<body>
						<a href="/basic3">Link 1</a>
						<a href="/basic4">Link 2</a>
					</body>
				</html>
`),
				},
			},
			expectedLinks: []string{"/basic1", "/basic2", "/basic3", "/basic4"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			allResponses := map[string][]byte{}
			for _, req := range tc.requests {
				allResponses[req.path] = req.resp
			}

			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if resp, ok := allResponses[r.URL.Path]; ok {
					w.Write(resp)
				}
			}))
			defer testServer.Close()

			domain := strings.TrimPrefix(testServer.URL, "http://")

			Run(domain)

			resultFileName := fmt.Sprintf("%s_result.csv", strings.ReplaceAll(domain, "/", "_"))
			resultFile, err := os.Open(resultFileName)
			if err != nil {
				t.Fatal(err)
			}
			defer resultFile.Close()

			csvReader := csv.NewReader(resultFile)
			results, err := csvReader.ReadAll()
			if err != nil {
				t.Fatal(err)
			}

			require.NoError(t, os.Remove(resultFileName))

			var resultSlice []string
			for _, result := range results {
				resultSlice = append(resultSlice, strings.TrimPrefix(result[0], testServer.URL))
			}

			require.ElementsMatch(t, tc.expectedLinks, resultSlice)
		})
	}
}
