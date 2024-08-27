package crawler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindAllLinks(t *testing.T) {
	t.Parallel()

	monzoHTML, err := os.ReadFile("./testdata/monzo.html")
	if err != nil {
		t.Fatal(err)
	}

	type testCase struct {
		name string

		path string // unique
		resp []byte

		expectedLinks []string
	}

	testCases := []testCase{
		{
			name: "basic",
			path: "/basic",
			resp: []byte(`
				<html>
					<body>
						<a href="https://example.com/link1">Link 1</a>
						<a href="https://example.com/link2">Link 2</a>
					</body>
				</html>
			`),
			expectedLinks: []string{"https://example.com/link1", "https://example.com/link2"},
		},
		{
			name: "basicWithImg",
			path: "/basic/img",
			resp: []byte(`
				<html>
					<body>
						<a href="https://example.com/link1">Link 1</a>
						<a href="https://example.com/link2">Link 2</a>
						<img src="https://example.com/image1" />
					</body>
				</html>
`),
			expectedLinks: []string{"https://example.com/link1", "https://example.com/link2"},
		},
		{
			name: "nested",
			path: "/nested",
			resp: []byte(`
				<html>
					<body>
						<a href="https://example.com/link1">Link 1</a>		
						<div>	
							<a href="https://example.com/link2">Link 2</a>
							<ul>
								<li>
									<a href="https://example.com/link3">Link 3</a>
								</li>
							</ul>
						</div>	
					</body>	
				</html>
`),
			expectedLinks: []string{"https://example.com/link1", "https://example.com/link2", "https://example.com/link3"},
		},
		{
			name: "monzo",
			path: "/monzo",
			resp: monzoHTML,
			expectedLinks: []string{
				"/features/travel#what-is-the-european-economic-area",
				"https://monzo.com/fraud/",
				"/savings/instant-access",
				"/help/monzo-premium",
				"https://monzo.com/i/helping-everyone-belong-at-monzo/",
				"https://monzo.com/tone-of-voice/",
				"https://monzo.com/legal/privacy-notice/",
				"https://twitter.com/monzo",
				"/investments",
				"/help",
				"/service-quality-results#personal-great-britain",
				"https://monzo.com/help/",
				"/refer-a-friend",
				"https://www.youtube.com/monzobank",
				"/security",
				"/loans",
				"/switch",
				"/service-quality-results#personal-northern-ireland",
				"/accessibility",
				"https://monzo.com/modern-slavery-statements/",
				"https://monzo.com/legal/fscs-information/",
				"https://monzo.com/legal/browser-support-policy/",
				"/current-account/under-16s",
				"https://monzo.com/information-about-current-account-services/",
				"https://we83.adj.st/?adj_t=ydi27sn&adj_engagement_type=fallback_click&adj_redirect=https%3A%2F%2Fmonzo.com%2Fdownload",
				"https://app.adjust.com/9mq4ox7?engagement_type=fallback_click&fallback=https%3A%2F%2Fmonzo.com%2Fdownload&redirect_macos=https%3A%2F%2Fmonzo.com%2Fdownload",
				"#",
				"https://web.monzo.com/",
				"https://monzo.com/business-banking/",
				"https://www.facebook.com/monzobank",
				"/pots",
				"/overdrafts",
				"https://monzo.com/investor-information/",
				"/features/travel",
				"https://www.psr.org.uk/app-fraud-data",
				"https://monzo.com/faq/",
				"/money-worries",
				"https://monzo.com/legal/terms-and-conditions/",
				"https://monzo.com/legal/cookie-notice/",
				"https://monzo.com/legal/mobile-operating-system-support-policy/",
				"https://www.instagram.com/monzo",
				"/help/monzo-plus",
				"/features/see-your-mortgage",
				"/fraud",
				"/business-banking",
				"#mainContent",
				"/current-account/joint-account",
				"/features/get-paid-early",
				"https://we83.adj.st/home?adj_t=ydi27sn_9mq4ox7&adj_redirect=https%3A%2F%2Fmonzo.com%2Fsign-up%2F&adj_engagement_type=fallback_click",
				"/",
				"https://monzo.com/service-information/",
				"https://monzo.com/blog/",
				"https://monzo.com/press/",
				"https://www.linkedin.com/company/monzo-bank",
				"/features/travel#how-do-i-know-if-monzo-is-my-main-bank",
				"/savingwithmonzo",
				"/pensions",
				"/ecb-rates",
				"https://uk.trustpilot.com/review/www.monzo.com",
				"https://monzo.com/about/",
				"/current-account/plans",
				"/current-account/16-17",
				"/features/cashback",
				"/features/savings",
				"/flex",
				"https://monzo.com/us/",
				"/our-social-programme",
				"/current-account",
				"https://monzo.com/supporting-all-our-customers/",
				"/careers",
				"/shared-tabs-more",
			},
		},
	}

	allResponses := map[string][]byte{}
	for _, tc := range testCases {
		allResponses[tc.path] = tc.resp
	}

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if resp, ok := allResponses[r.URL.Path]; ok {
			w.Write(resp)
		}
	}))

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := FindAllLinks(testServer.URL + tc.path)

			resultSlice := make([]string, 0, len(result))
			for k := range result {
				resultSlice = append(resultSlice, k)
			}

			require.ElementsMatch(t, tc.expectedLinks, resultSlice)
		})
	}
}
