package gonmap

import (
	"crypto/x509"
	"github.com/PuerkitoBio/goquery"
	"goon3/lib/kscan/lib/urlparse"
	"io"
	"io/ioutil"
	"goon3/lib/gonmap/shttp"
	"goon3/lib/kscan/lib/httpfinger"
	"goon3/lib/kscan/lib/iconhash"
	"goon3/lib/kscan/lib/misc"
	"goon3/lib/kscan/lib/slog"
	"net/http"
)

type HttpFinger struct {
	URL              *urlparse.URL
	StatusCode       int
	Response         string
	ResponseDigest   string
	Title            string
	Header           string
	HeaderDigest     string
	HashFinger       string
	KeywordFinger    string
	PeerCertificates *x509.Certificate
}

func NewHttpFinger(url *urlparse.URL) *HttpFinger {
	return &HttpFinger{
		URL:              url,
		StatusCode:       0,
		Response:         "",
		ResponseDigest:   "",
		Title:            "",
		Header:           "",
		HashFinger:       "",
		KeywordFinger:    "",
		PeerCertificates: nil,
	}
}

func (h *HttpFinger) LoadHttpResponse(url *urlparse.URL, resp *http.Response) {
	h.Title = getTitle(shttp.GetBody(resp))
	h.StatusCode = resp.StatusCode
	h.Header = getHeader(resp.Header.Clone())
	h.HeaderDigest = getHeaderDigest(resp.Header.Clone())
	h.Response = getResponse(shttp.GetBody(resp))
	h.ResponseDigest = getResponseDigest(shttp.GetBody(resp))
	h.HashFinger = getFingerByHash(*url)
	h.KeywordFinger = getFingerByKeyword(h.Header, h.Title, h.Response)
	_ = resp.Body.Close()
}

func getTitle(resp io.Reader) string {
	query, err := goquery.NewDocumentFromReader(resp)
	if err != nil {
		slog.Debug(err.Error())
		return ""
	}
	result := query.Find("title").Text()
	result = misc.FixLine(result)
	//Body.Close()
	return result
}

func getHeader(header http.Header) string {
	return shttp.Header2String(header)
}

func getResponse(resp io.Reader) string {
	body, err := ioutil.ReadAll(resp)
	if err != nil {
		slog.Debug(err.Error())
		return ""
	}
	bodyStr := string(body)
	return bodyStr
}

func getResponseDigest(resp io.Reader) string {

	var result string

	query, err := goquery.NewDocumentFromReader(CopyIoReader(&resp))
	if err != nil {
		slog.Debug(err.Error())
		return ""
	}

	query.Find("script").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	query.Find("style").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	query.Find("textarea").Each(func(_ int, tag *goquery.Selection) {
		tag.Remove() // 把无用的 tag 去掉
	})
	query.Each(func(_ int, tag *goquery.Selection) {
		result = result + tag.Text()
	})

	result = misc.FixLine(result)

	result = misc.FilterPrintStr(result)

	result = misc.StrRandomCut(result, 20)

	if len(result) == 0 {
		b, _ := ioutil.ReadAll(CopyIoReader(&resp))
		result = string(b)
		result = misc.FixLine(result)
		result = misc.FilterPrintStr(result)
		result = misc.StrRandomCut(result, 20)
	}

	return result
}

func getHeaderDigest(header http.Header) string {
	if header.Get("SERVER") != "" {
		return "server:" + header.Get("SERVER")
	}
	return ""
}

func getFingerByKeyword(header string, title string, body string) string {
	return httpfinger.KeywordFinger.Match(header, title, body)
}

func getFingerByHash(url urlparse.URL) string {
	resp, err := shttp.GetFavicon(url)
	if err != nil {
		slog.Debug(url.UnParse() + err.Error())
		return ""
	}
	if resp.StatusCode != 200 {
		slog.Debug(url.UnParse() + "没有图标文件")
		return ""
	}
	hash, err := iconhash.Get(resp.Body)
	if err != nil {
		slog.Debug(url.UnParse() + err.Error())
		return ""
	}
	_ = resp.Body.Close()
	return httpfinger.FaviconHash.Match(hash)
}

func (h *HttpFinger) LoadCert(resp *http.Response) {
	h.PeerCertificates = resp.TLS.PeerCertificates[0]
}
