package gonmap

import (
	"fmt"
	"goon3/lib/kscan/lib/chinese"
	"goon3/lib/kscan/lib/misc"
	"strconv"
	"strings"
)

type AppBanner struct {
	//端口号
	Port int
	//IP地址
	IPAddr string
	//端口开放状态码
	StatusCode int
	//HTTP协议标题，其他协议正文摘要
	AppDigest string
	//返回包全文
	Response string
	//端口开放协议类型
	Protocol string
	//fingerprintMap
	fingerPrint map[string]string
}

func NewAppBanner() *AppBanner {
	banner := &AppBanner{}
	banner.fingerPrint = make(map[string]string)
	return banner
}

func (a *AppBanner) URL() string {
	return fmt.Sprintf("%s://%s:%d", a.Protocol, a.IPAddr, a.Port)
}

func (a *AppBanner) LoadHttpFinger(finger *HttpFinger) {
	a.IPAddr = finger.URL.Netloc
	a.Port = finger.URL.Port
	a.AppDigest = finger.Title
	a.StatusCode = finger.StatusCode
	a.Response = finger.Response
	a.Protocol = finger.URL.Scheme
	a.SetCertSubject(func() string {
		if finger.PeerCertificates != nil {
			strCert := finger.PeerCertificates.Subject.String()
			strCert = chinese.ToUTF8(strCert)
			return strCert
		} else {
			return ""
		}
	}())
	a.SetResponseDigest(finger.ResponseDigest)
	a.SetHashFinger(finger.HashFinger)
	a.SetKeywordFinger(finger.KeywordFinger)
	a.SetInfo(finger.HeaderDigest)
	if a.AppDigest == "" {
		switch a.StatusCode {
		case 100:
			a.AppDigest = "100 Continue"
		case 101:
			a.AppDigest = "101 Switching Protocols"
		case 201:
			a.AppDigest = "201 Created"
		case 202:
			a.AppDigest = "202 Accepted"
		case 203:
			a.AppDigest = "203 Non-Authoritative Information"
		case 204:
			a.AppDigest = "204 No Content"
		case 205:
			a.AppDigest = "205 Reset Content"
		case 206:
			a.AppDigest = "206 Partial Content"
		case 300:
			a.AppDigest = "300 Multiple Choices"
		case 301:
			a.AppDigest = "301 Moved Permanently"
		case 302:
			a.AppDigest = "302 Found"
		case 303:
			a.AppDigest = "303 See Other"
		case 304:
			a.AppDigest = "304 Not Modified"
		case 305:
			a.AppDigest = "305 Use Proxy"
		case 306:
			a.AppDigest = "306 Unused"
		case 307:
			a.AppDigest = "307 Temporary Redirect"
		case 400:
			a.AppDigest = "400 Bad Request"
		case 401:
			a.AppDigest = "401 Unauthorized"
		case 402:
			a.AppDigest = "402 Payment Required"
		case 403:
			a.AppDigest = "403 Forbidden"
		case 404:
			a.AppDigest = "404 Not Found"
		case 405:
			a.AppDigest = "405 Method Not Allowed"
		case 406:
			a.AppDigest = "406 Not Acceptable"
		case 407:
			a.AppDigest = "407 Proxy Authentication Required"
		case 408:
			a.AppDigest = "408 Request Time-out"
		case 409:
			a.AppDigest = "409 Conflict"
		case 410:
			a.AppDigest = "410 Gone"
		case 411:
			a.AppDigest = "411 Length Required"
		case 412:
			a.AppDigest = "412 Precondition Failed"
		case 413:
			a.AppDigest = "413 Request Entity Too Large"
		case 414:
			a.AppDigest = "414 Request-URI Too Large"
		case 415:
			a.AppDigest = "415 Unsupported Media Type"
		case 416:
			a.AppDigest = "416 Requested range not satisfiable"
		case 417:
			a.AppDigest = "417 Expectation Failed"
		case 500:
			a.AppDigest = "500 Internal Server Error"
		case 501:
			a.AppDigest = "501 Not Implemented"
		case 502:
			a.AppDigest = "502 Bad Gateway"
		case 503:
			a.AppDigest = "503 Service Unavailable"
		case 504:
			a.AppDigest = "504 Gateway Time-out"
		case 505:
			a.AppDigest = "505 HTTP Version not supported"
		default:
			a.AppDigest = "No Title"
		}
	}

	if finger.StatusCode == 0 {
		a.Protocol = "unknown"
	}

}

func (a *AppBanner) LoadTcpBanner(banner *TcpBanner) {
	if a.StatusCode == 0 {
		if banner.TcpFinger.Service == "http" || banner.TcpFinger.Service == "https" {
			a.StatusCode = 500
		} else {
			a.StatusCode = 200
		}
		a.Protocol = func() string {
			if banner.TcpFinger.Service == "" {
				banner.TcpFinger.Service = "unknow"
			}
			return banner.TcpFinger.Service
		}()
		a.Port = misc.Str2Int(strings.Split(banner.Target.Uri, ":")[1])
		a.IPAddr = strings.Split(banner.Target.Uri, ":")[0]

		a.Response = banner.Response.String
		a.AppDigest = func() string {
			appDigest := misc.FixLine(a.Response)
			appDigest = misc.FilterPrintStr(appDigest)
			appDigest = misc.MustLength(appDigest, 10)
			return appDigest
		}()
	}

	a.SetProductName(banner.TcpFinger.ProductName)
	a.SetInfo(banner.TcpFinger.Info)
	a.SetDeviceType(banner.TcpFinger.DeviceType)
	a.SetOperatingSystem(banner.TcpFinger.OperatingSystem)
	a.SetHostname(banner.TcpFinger.Hostname)
	a.SetVersion(banner.TcpFinger.Version)
}

func (a *AppBanner) Output() string {
	fingerPrint := misc.SprintStringMap(a.fingerPrint)

	a.AppDigest = chinese.ToUTF8(a.AppDigest)

	return fmt.Sprintf("%s\t%d\t%s\t%s", a.URL(), a.StatusCode, a.AppDigest, fingerPrint)
}

//返回包摘要
func (a *AppBanner) SetResponseDigest(s string) {
	a.fingerPrint["ResponseDigest"] = s
}

//Http IconHash指纹识别信息
func (a *AppBanner) SetHashFinger(s string) {
	a.fingerPrint["HashFinger"] = s
}

//Http 关键字指纹识别信息
func (a *AppBanner) SetKeywordFinger(s string) {
	a.fingerPrint["KeywordFinger"] = s
}

//Https 证书信息
func (a *AppBanner) SetCertSubject(s string) {
	a.fingerPrint["CertSubject"] = s
}

//端口开放 产品信息
func (a *AppBanner) SetProductName(s string) {
	a.fingerPrint["ProductName"] = s
}

//端口开放 产品版本信息
func (a *AppBanner) SetVersion(s string) {
	a.fingerPrint["Version"] = s
}

//主机名称
func (a *AppBanner) SetHostname(s string) {
	a.fingerPrint["Hostname"] = s
}

//操作系统名称
func (a *AppBanner) SetOperatingSystem(s string) {
	a.fingerPrint["OperatingSystem"] = s
}

//设备类型
func (a *AppBanner) SetDeviceType(s string) {
	a.fingerPrint["DeviceType"] = s
}

//端口其他信息
func (a *AppBanner) SetInfo(s string) {
	a.fingerPrint["Info"] = s
}

func (a *AppBanner) Map() map[string]string {
	bannerMap := make(map[string]string)
	bannerMap["Response"] = a.Response
	bannerMap["URL"] = a.URL()
	bannerMap["Port"] = strconv.Itoa(a.Port)
	bannerMap["IPAddr"] = a.IPAddr
	bannerMap["AppDigest"] = a.AppDigest
	bannerMap["Protocol"] = a.Protocol
	bannerMap["StatusCode"] = strconv.Itoa(a.StatusCode)
	for key, value := range a.fingerPrint {
		bannerMap[key] = value
	}
	return bannerMap
}
