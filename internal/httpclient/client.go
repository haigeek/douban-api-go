package httpclient

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	originHeader  = "https://movie.douban.com"
	refererHeader = "https://movie.douban.com/"
	uaHeader      = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.131 Safari/537.36"
)

var (
	reTok        = regexp.MustCompile(`id="tok"[^>]*value="([^"]+)"`)
	reCha        = regexp.MustCompile(`id="cha"[^>]*value="([^"]+)"`)
	reRed        = regexp.MustCompile(`id="red"[^>]*value="([^"]+)"`)
	reDifficulty = regexp.MustCompile(`difficulty\s*=\s*(\d+)`)
)

type Client struct {
	httpClient *http.Client
}

func New(cookieValue string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(cookieValue) != "" {
		baseURL, _ := url.Parse("https://douban.com/")
		cookies := make([]*http.Cookie, 0)
		for _, part := range strings.Split(cookieValue, ";") {
			item := strings.TrimSpace(part)
			if item == "" {
				continue
			}
			kv := strings.SplitN(item, "=", 2)
			if len(kv) != 2 {
				continue
			}
			cookies = append(cookies, &http.Cookie{
				Name:   strings.TrimSpace(kv[0]),
				Value:  strings.TrimSpace(kv[1]),
				Domain: "douban.com",
				Path:   "/",
			})
		}
		jar.SetCookies(baseURL, cookies)
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &Client{
		httpClient: &http.Client{
			Timeout:   30 * time.Second,
			Jar:       jar,
			Transport: transport,
		},
	}, nil
}

func (c *Client) Get(ctx context.Context, rawURL string, query map[string]string, strictStatus bool) (*http.Response, error) {
	return c.get(ctx, rawURL, query, strictStatus, true)
}

func (c *Client) get(ctx context.Context, rawURL string, query map[string]string, strictStatus bool, allowChallenge bool) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}

	if len(query) > 0 {
		q := req.URL.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Origin", originHeader)
	req.Header.Set("Referer", refererHeader)
	req.Header.Set("User-Agent", uaHeader)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if allowChallenge {
		ok, err := c.solveChallengeIfNeeded(ctx, resp)
		if err != nil {
			return nil, err
		}
		if ok {
			return c.get(ctx, rawURL, query, strictStatus, false)
		}
	}

	if strictStatus && resp.StatusCode >= http.StatusBadRequest {
		defer resp.Body.Close()
		return nil, errors.New(fmt.Sprintf("upstream status %d", resp.StatusCode))
	}

	return resp, nil
}

func (c *Client) solveChallengeIfNeeded(ctx context.Context, resp *http.Response) (bool, error) {
	if resp == nil || resp.Request == nil {
		return false, nil
	}
	if !strings.Contains(resp.Request.URL.Host, "sec.douban.com") {
		return false, nil
	}

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return false, err
	}
	html := string(body)

	tok := capture(reTok, html)
	cha := capture(reCha, html)
	red := capture(reRed, html)
	if tok == "" || cha == "" || red == "" {
		return false, errors.New("douban sec challenge parse failed")
	}

	difficulty := 4
	if v := capture(reDifficulty, html); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			difficulty = n
		}
	}
	sol := solvePow(cha, difficulty)

	form := url.Values{}
	form.Set("tok", tok)
	form.Set("cha", cha)
	form.Set("sol", strconv.Itoa(sol))
	form.Set("red", red)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://sec.douban.com/c", strings.NewReader(form.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", uaHeader)
	req.Header.Set("Referer", resp.Request.URL.String())
	req.Header.Set("Origin", "https://sec.douban.com")

	chkResp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer chkResp.Body.Close()
	_, _ = io.Copy(io.Discard, chkResp.Body)
	return true, nil
}

func solvePow(data string, difficulty int) int {
	target := strings.Repeat("0", difficulty)
	for nonce := 1; ; nonce++ {
		sum := sha512.Sum512([]byte(data + strconv.Itoa(nonce)))
		if strings.HasPrefix(hex.EncodeToString(sum[:]), target) {
			return nonce
		}
	}
}

func capture(re *regexp.Regexp, text string) string {
	m := re.FindStringSubmatch(text)
	if len(m) < 2 {
		return ""
	}
	return strings.TrimSpace(m[1])
}
