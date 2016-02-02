package extract

import (
	"bytes"
	"log"
	urlparse "net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var (
	logger    = log.New(os.Stderr, "[extract] ", log.LstdFlags|log.Lshortfile)
	titleTags = "h1,h2,h3,h4,h5,h6"
)

type Link struct {
	Url     string
	Anchor  string
	Inner   bool
	Remarks []string
}

type PageInfo struct {
	Url         string
	Host        string
	Title       string
	Description string
	Content     string
	Remarks     []string
	Links       []*Link
}

type ExtractFunc func(string, []byte) (*PageInfo, error)

func Extract(url string, body []byte) (*PageInfo, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	doc.Url, err = urlparse.Parse(url)
	if err != nil {
		return nil, err
	}
	return ExtractFromDocument(doc)
}

func ExtractFromDocument(doc *goquery.Document) (*PageInfo, error) {
	pageInfo := &PageInfo{
		Url:   doc.Url.String(),
		Host:  doc.Url.Host,
		Links: getLinks(doc),
	}

	contentNode := getContentNode(doc)
	if contentNode == nil {
		logger.Printf("get content node fail")
		return pageInfo, nil
	}

	pageInfo.Content = getNodeText(contentNode, -1)

	titleNode := getTitleNode(doc)
	if titleNode == nil {
		logger.Printf("get title node fail")
		return pageInfo, nil
	}

	pageInfo.Title = getNodeText(titleNode, 1)

	if titleNode.Data != "title" {
		remarkNodes := getMiddleNodes(titleNode, contentNode)
		for _, remarkNode := range remarkNodes {
			if remarkNode.Type == html.TextNode {
				continue
			}
			remark := getNodeText(remarkNode, 1)
			if remark == "" {
				continue
			}
			pageInfo.Remarks = append(pageInfo.Remarks, remark)
		}
	}
	return pageInfo, nil
}

func getLinks(doc *goquery.Document) []*Link {
	links := []*Link{}
	urlMap := make(map[string]int)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok {
			return
		}

		href = strings.TrimSpace(href)
		uri, err := urlparse.Parse(href)
		if err != nil {
			logger.Printf("urlparse error: %s", err)
			return
		}

		if uri.IsAbs() && uri.Scheme != "http" && uri.Scheme != "https" {
			return
		}

		if !uri.IsAbs() {
			uri.Host = doc.Url.Host
			uri.Scheme = doc.Url.Scheme
		}

		url := uri.String()
		if n, ok := urlMap[url]; ok {
			urlMap[url] = n + 1
			return
		}

		link := &Link{
			Url:    url,
			Anchor: strings.TrimSpace(s.Text()),
			Inner:  uri.Host == doc.Url.Host,
		}
		links = append(links, link)
		urlMap[url] = 1
	})

	return links
}

func getContentNode(doc *goquery.Document) *html.Node {
	maybe := make(map[*html.Node]int)

	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		score := len(text) - len(s.Find("a").Text())
		if score <= 0 {
			return
		}

		parentSel := s.Parent()
		if len(parentSel.Nodes) == 0 {
			return
		}

		parent := parentSel.Get(0)
		grandParent := parent.Parent
		if grandParent == nil {
			return
		}

		if _, ok := maybe[parent]; !ok {
			maybe[parent] = 0
		}
		if _, ok := maybe[grandParent]; !ok {
			maybe[grandParent] = 0
		}

		maybe[parent] += score
		maybe[grandParent] += score / 2
	})

	bestNode := &html.Node{}
	bestScore := 0
	for node, score := range maybe {
		logger.Printf("node<%s>(%d): %s", node.Data, score, getNodeText(node, -1))
		if score > bestScore {
			bestNode = node
			bestScore = score
		}
	}

	return bestNode
}

func getTitleNode(doc *goquery.Document) *html.Node {
	titleNode := doc.Find("title").Get(0)
	if titleNode == nil {
		return nil
	}

	defaultTitle := getNodeText(titleNode, 1)
	defaultTitle = strings.Replace(defaultTitle, " ", "", -1)

	bestNode := titleNode
	bestScore := 0

	doc.Find(titleTags).Each(func(i int, s *goquery.Selection) {
		text := strings.Replace(s.Text(), " ", "", -1)
		textLen := len(text)
		if textLen < 5 || textLen > len(defaultTitle)*2 {
			return
		}

		score, _ := LCS([]rune(text), []rune(defaultTitle))
		if score > bestScore {
			bestNode = s.Get(0)
			bestScore = score
		}
	})

	return bestNode
}

func getNodeText(node *html.Node, depth int) string {
	text := ""
	if node.Type == html.TextNode {
		text += strings.TrimSpace(node.Data)
	}

	if depth != 0 {
		child := node.FirstChild
		for child != nil {
			text += getNodeText(child, depth-1)
			child = child.NextSibling
		}
	}
	return text
}

func getMiddleNodes(startNode *html.Node, endNode *html.Node) []*html.Node {
	middleNodes := make([]*html.Node, 0, 100)
	node := startNode
	for node != nil {
		if node == endNode {
			break
		}

		middleNodes = append(middleNodes, node)

		if node.FirstChild != nil {
			node = node.FirstChild
		} else if node.NextSibling != nil {
			node = node.NextSibling
		} else if node.Parent != nil {
			node = node.Parent
			for node != nil {
				if node.NextSibling != nil {
					node = node.NextSibling
					break
				}
				node = node.Parent
			}
		} else {
			break
		}
	}
	return middleNodes[1:]
}

func LCS(a, b []rune) (int, []rune) {
	lengths := make([][]int, len(a)+1)
	for i := 0; i <= len(a); i++ {
		lengths[i] = make([]int, len(b)+1)
	}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b); j++ {
			if a[i] == b[j] {
				lengths[i+1][j+1] = lengths[i][j] + 1
			} else if lengths[i+1][j] > lengths[i][j+1] {
				lengths[i+1][j+1] = lengths[i+1][j]
			} else {
				lengths[i+1][j+1] = lengths[i][j+1]
			}
		}
	}

	s := make([]rune, 0, lengths[len(a)][len(b)])
	for x, y := len(a), len(b); x != 0 && y != 0; {
		if lengths[x][y] == lengths[x-1][y] {
			x--
		} else if lengths[x][y] == lengths[x][y-1] {
			y--
		} else {
			s = append(s, a[x-1])
			x--
			y--
		}
	}

	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return len(s), s
}
