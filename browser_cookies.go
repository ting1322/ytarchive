package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all" // register cookie store finders!
)

func (di *DownloadInfo) ImportBrowserCookie() (*cookiejar.Jar, error) {
	//burl, err := url.Parse("https://youtube.com")
	//if err != nil {
	//	return nil, err
	//}
	cookies := kooky.ReadCookies(kooky.Valid)

	cookieMap := make(map[string][]*http.Cookie)
	for _, c := range cookies {
		domain := c.Domain

		if _, ok := cookieMap[domain]; !ok {
			cookieMap[domain] = make([]*http.Cookie, 0)
		}

		cookieMap[domain] = append(cookieMap[domain], &c.Cookie)
	}
	jar, err := cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
	if err != nil {
		return nil, err
	}
	if len(cookieMap) > 0 {
		for _, cookies := range cookieMap {
			url, err := url.Parse(fmt.Sprintf("https://%s", cookies[0].Domain))

			if err == nil {
				jar.SetCookies(url, cookies)
				if strings.HasSuffix(url.Host, "youtube.com") {
					di.CookiesURL = url
				}
			}
		}
	}
	return jar, nil
}