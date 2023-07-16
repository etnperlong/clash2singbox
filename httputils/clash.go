package httputils

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/etnperlong/clash2singbox/model/clash"
	"gopkg.in/yaml.v3"
)

func filter(isinclude bool, reg string, sl []clash.Proxies) []clash.Proxies {
	r := regexp.MustCompile(reg)
	return getForList(sl, func(v clash.Proxies) (clash.Proxies, bool) {
		has := r.MatchString(v.Name)
		if has && isinclude {
			return v, true
		}
		if !isinclude && !has {
			return v, true
		}
		return clash.Proxies{}, false
	})
}

func getForList[K, V any](l []K, check func(K) (V, bool)) []V {
	sl := make([]V, 0, len(l))
	for _, v := range l {
		s, ok := check(v)
		if !ok {
			continue
		}
		sl = append(sl, s)
	}
	return sl
}

func GetClash(cxt context.Context, hc *http.Client, url string, include, exclude string) (clash.Clash, error) {
	urls := strings.Split(url, "|")

	c := clash.Clash{}

	for _, v := range urls {
		b, err := HttpGet(context.TODO(), hc, v)
		if err != nil {
			return c, fmt.Errorf("GetClash: %w", err)
		}
		lc := clash.Clash{}
		err = yaml.Unmarshal(b, &lc)
		if err != nil {
			return c, fmt.Errorf("GetClash: %w", err)
		}
		if include != "" {
			lc.Proxies = filter(true, include, lc.Proxies)
		}
		if exclude != "" {
			lc.Proxies = filter(false, exclude, lc.Proxies)
		}
		c.Proxies = append(c.Proxies, lc.Proxies...)
	}
	return c, nil
}
