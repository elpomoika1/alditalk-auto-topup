package browser

import (
	"context"
	"fmt"
	"net/http"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
)

type Cookie struct {
	Name  string
	Value string
	Path  string
}

func GetCookies(domainPrefix string) []Cookie {
	cookiesSeq := kooky.TraverseCookies(
		context.TODO(),
		kooky.Valid,
		kooky.DomainHasPrefix(domainPrefix)).OnlyCookies()

	cookies := make([]Cookie, 0)

	for cookie := range cookiesSeq {
		cookies = append(cookies, Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
			Path:  cookie.Path,
		})
	}

	return cookies
}

func ApplyCookiesToReq(req *http.Request, domainPrefix string) {
	cookies := GetCookies(domainPrefix)
	fmt.Print(cookies)

	for _, cookie := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}
}
