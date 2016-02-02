package extract

import (
	urlparse "net/url"

	"github.com/rjsamson/rss"
)

func ExtractRss(url string, body []byte) (*PageInfo, error) {
	uri, err := urlparse.Parse(url)
	if err != nil {
		return nil, err
	}

	feed, err := rss.Parse(body)
	if err != nil {
		return nil, err
	}

	pageInfo := &PageInfo{
		Url:         uri.String(),
		Host:        uri.Host,
		Title:       feed.Title,
		Description: feed.Description,
	}

	for _, item := range feed.Items {
		itemUri, err := urlparse.Parse(item.Link)
		if err != nil {
			continue
		}

		link := &Link{}
		link.Url = itemUri.String()
		link.Anchor = item.Title
		link.Inner = (itemUri.Host == uri.Host)

		if item.Author != "" {
			link.Remarks = append(link.Remarks, "author:"+item.Author)
		}

		if item.PubDate != "" {
			link.Remarks = append(link.Remarks, "date:"+item.PubDate)
		}

		if item.Summary != "" {
			link.Remarks = append(link.Remarks, "summary:"+item.Summary)
		}

		pageInfo.Links = append(pageInfo.Links, link)
	}

	return pageInfo, nil
}
