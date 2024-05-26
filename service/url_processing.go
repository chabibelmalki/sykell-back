package service

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sykell-back/utils"

	"golang.org/x/net/html"
)

type URLInfo struct {
	URL                    string      `json:"-"`
	HTMLVersion            string      `json:"htmlVersion"`
	Title                  string      `json:"title"`
	HeadingTagsCount       map[int]int `json:"headingTagsCount"`
	InternalLinksCount     int         `json:"internalLinksCount"`
	ExternalLinksCount     int         `json:"externalLinksCount"`
	InaccessibleLinksCount int         `json:"inaccessibleLinksCount"`
	IsLoginPage            bool        `json:"isLoginPage"`

	//technical fields
	rootNode *html.Node `json:"-"`
}

func UrlProcess(url string) (info *URLInfo, err error) {

	// Fetch the URL content
	responseBody, err := utils.FetchURLContent(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return
	}

	// Parse the HTML content
	rootNode, err := html.Parse(strings.NewReader(responseBody))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Init url info
	info = &URLInfo{
		URL:              url,
		HeadingTagsCount: make(map[int]int),
		rootNode:         rootNode,
	}

	//process
	info.extractTitle()
	info.extractHTMLVersion()
	info.countHeadingTags()
	info.checkLoginForm()
	info.countLinks()

	return info, nil
}

// traverseNodes traverse the HTML nodes and applies the provided function
func traverseNodes(n *html.Node, fn func(*html.Node) bool) {
	if fn(n) {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNodes(c, fn)
	}
}

// extractTitle extracts the title from the HTML document
func (urlinfo *URLInfo) extractTitle() {
	traverseNodes(urlinfo.rootNode, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				urlinfo.Title = n.FirstChild.Data
				return true
			}
		}
		return false
	})
}

// extractHTMLVersion extracts the HTML version from the doctype declaration
func (urlinfo *URLInfo) extractHTMLVersion() {
	traverseNodes(urlinfo.rootNode, func(n *html.Node) bool {
		if n.Type == html.DoctypeNode {
			doctype := strings.ToLower(n.Data)
			if strings.Contains(doctype, "html") {
				if strings.Contains(doctype, "html 4.01") {
					urlinfo.HTMLVersion = "HTML 4.01"
				} else if strings.Contains(doctype, "xhtml 1.0") {
					urlinfo.HTMLVersion = "XHTML 1.0"
				} else if strings.Contains(doctype, "xhtml 1.1") {
					urlinfo.HTMLVersion = "XHTML 1.1"
				} else {
					urlinfo.HTMLVersion = "HTML5"
				}
				return true
			}
		}
		return false
	})
}

// countHeadingTags
func (urlinfo *URLInfo) countHeadingTags() {
	headingRegex := regexp.MustCompile(`^h([1-6])$`)

	traverseNodes(urlinfo.rootNode, func(n *html.Node) bool {
		if n.Type == html.ElementNode && headingRegex.MatchString(n.Data) {
			level, _ := strconv.Atoi(string(n.Data[1]))
			urlinfo.HeadingTagsCount[level]++
		}
		return false
	})
}

// countLinks counts internal and external links
func (urlinfo *URLInfo) countLinks() {
	parsedBaseURL, _ := url.Parse(urlinfo.URL)

	traverseNodes(urlinfo.rootNode, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {

					parsedLinkURL, _ := url.Parse(attr.Val)

					resolvedURL := parsedBaseURL.ResolveReference(parsedLinkURL)

					if resolvedURL.Hostname() != parsedBaseURL.Hostname() {
						urlinfo.ExternalLinksCount++

						//Check inaccessible link
						statut, err := utils.CheckUrl(attr.Val)
						if err != nil || (statut >= 400 && statut < 600) {
							urlinfo.InaccessibleLinksCount++
						}

					} else {
						urlinfo.InternalLinksCount++
					}
				}
			}
		}
		return false
	})
}

// checkLoginForm checks if a page contains a login form
func (urlinfo *URLInfo) checkLoginForm() {
	hasPasswordField := false
	hasSubmitButton := false
	formDetected := false

	traverseNodes(urlinfo.rootNode, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "form" {
			formDetected = true
		}
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, attr := range n.Attr {
				if attr.Key == "type" && attr.Val == "password" {
					hasPasswordField = true
				}
				if attr.Key == "type" && (attr.Val == "submit" || attr.Val == "button") {
					hasSubmitButton = true
				}
			}
		}
		return false
	})

	if formDetected && hasPasswordField && hasSubmitButton {
		urlinfo.IsLoginPage = true
	}
}
