package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const CountUrls = 10

const ParseUrl = "https://lux.fm/"

var mainUrl = getHostName(ParseUrl)

func main() {
	links := getLinksByUrl(ParseUrl)

	for i := 0; i <= CountUrls; i++ {
		parseLink := parseUrl(links[i])
		link := ""
		if parseLink.Host == "" {
			link = mainUrl + parseLink.Path
		} else {
			link = links[i]
		}
		printTitleUrl(link)
	}
}

func printTitleUrl(siteUrl string) {
	title := getTitleByUrl(siteUrl)
	fmt.Println("url: " + siteUrl + " | title: " + title)
}

func getLinksByUrl(siteUrl string) []string {

	links := make([]string, 0, CountUrls)

	resp, err := http.Get(siteUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	bodyString := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
	if err != nil {
		panic(err)
	}

	doc.Find("a").Each(func(index int, link *goquery.Selection) {
		href, ok := link.Attr("href")
		if ok && len(links) <= CountUrls {
			links = append(links, strings.Trim(href, " "))
		}
	})

	return links
}

func getTitleByUrl(siteUrl string) string {
	resp, err := http.Get(siteUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	bodyString := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
	if err != nil {
		panic(err)
	}
	href := ""
	doc.Find("title").Each(func(index int, link *goquery.Selection) {
		href = link.Text()
	})

	return href
}

func getHostName(siteUrl string) string {
	u := parseUrl(siteUrl)
	if u.Scheme == "" || u.Host == "" {
		err := errors.New("url not correct")
		if err != nil {
			return ""
		}
	}
	return u.Scheme + "://" + u.Host
}

func parseUrl(siteUrl string) *url.URL {
	u, err := url.Parse(siteUrl)
	if err != nil {
		panic(err)
	}
	return u
}
