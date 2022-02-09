package main

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	urlToParse := os.Args[1]
	mainUrl := getHostName(urlToParse)
	links := getLinksByUrl(urlToParse)

	countUrls := 10
	if len(os.Args) > 2 {
		cs, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
		countUrls = cs
	}

	countLinks := 1
	for i := 1; countLinks <= countUrls; i++ {
		parseLink := parseUrl(links[i])
		link := ""
		if parseLink.Host == "" {
			link = mainUrl + parseLink.Path
		} else {
			link = links[i]
		}
		countLinks++
		go printTitleUrl(link)
	}

	time.Sleep(time.Second)
	fmt.Println("done")
}

func printTitleUrl(siteUrl string) {
	title := getTitleByUrl(siteUrl)
	fmt.Println("url: " + siteUrl + " | title: " + title)
}

func getLinksByUrl(siteUrl string) []string {

	links := make([]string, 0)

	resp, err := http.Get(siteUrl)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	bodyString := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
	if err != nil {
		panic(err)
	}

	doc.Find("a").Each(func(index int, link *goquery.Selection) {
		href, ok := link.Attr("href")
		if ok {
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	bodyString := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(bodyString))
	if err != nil {
		panic(err)
	}
	title := ""
	doc.Find("title").Each(func(index int, attr *goquery.Selection) {
		title = attr.Text()
	})

	return title
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
