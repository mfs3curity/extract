package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	var wg sync.WaitGroup

	for sc.Scan() {
		ur := sc.Text()
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			onColly(url, false)
		}(ur)
	}

	wg.Wait()
}

func extractHost(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

func onColly(u string, isProxy bool) {

	uuid := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + extractHost(u)
	fmt.Println("Folder Name:", uuid)

	err := os.Mkdir(uuid, os.ModePerm)
	if err != nil {
		log.Println("Error creating directory:", err)
		return
	}
	c := colly.NewCollector()

	if isProxy {
		proxyURL, err := url.Parse("http://127.0.0.1:8080")
		if err != nil {
			fmt.Println("Error parsing proxy URL:", err)
			return
		}
		transport := &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		c.WithTransport(transport)
	}

	c.OnHTML("script[src]", func(h *colly.HTMLElement) {
		l := h.Attr("src")

		f, _ := openFile(uuid + "/js.txt")
		defer f.Close()
		f.WriteString(l + "\n")
	})

	c.OnHTML("a[href]", func(h *colly.HTMLElement) {
		l := h.Attr("href")
		f, _ := openFile(uuid + "/href.txt")
		defer f.Close()
		f.WriteString(l + "\n")
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36")
	})

	c.OnResponse(func(r *colly.Response) {
		f, _ := openFile(uuid + "/headers.txt")
		defer f.Close()
		for k, v := range *r.Headers {
			f.WriteString(fmt.Sprintf("%s: %s\n", k, strings.Join(v, ", ")))
		}

		body := string(r.Body)
		f2, _ := openFile(uuid + "/body.txt")
		defer f2.Close()
		f2.WriteString(body + "\n")

		re := regexp.MustCompile(`<!--(.*?)-->`)

		comments := re.FindAllStringSubmatch(string(r.Body), -1)

		f3, _ := openFile(uuid + "/comments.txt")
		defer f3.Close()
		for _, comment := range comments {

			f3.WriteString(comment[1] + "\n")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error visiting:", r.Request.URL, "Error:", err)
	})

	c.Visit(u)
}

func openFile(filename string) (*os.File, error) {
	var file *os.File
	var err error

	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening file for appending: %v", err)
	}

	return file, nil
}
