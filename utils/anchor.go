package utils

import "strings"

// GetText extracts text content between anchor tags: <a href=...>TEXT</a>
func GetText(tag string) string {
	start := strings.Index(tag, ">") + 1
	finish := strings.Index(tag, "</a")
	if start > 0 && finish >= start {
		return tag[start:finish]
	}
	return ""
}

// GetHref extracts the href value from an anchor tag
func GetHref(tag string) string {
	start := strings.Index(tag, "href") + 5
	finish := strings.Index(tag, ">")
	if start >= 5 && finish >= start {
		href := tag[start:finish]
		href = strings.ReplaceAll(href, "'", "")
		href = strings.ReplaceAll(href, "\"", "")
		return href
	}
	return ""
}
