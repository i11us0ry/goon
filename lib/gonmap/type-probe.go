package gonmap

import (
	"errors"
	"goon3/lib/gonmap/simplenet"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var PROBE_LOAD_REGEXP = regexp.MustCompile("^(UDP|TCP) ([a-zA-Z0-9-_./]+) (?:q\\|([^|]*)\\|)$")
var PROBE_INT_REGEXP = regexp.MustCompile(`^(\d+)$`)
var PROBE_STRING_REGEXP = regexp.MustCompile(`^([a-zA-Z0-9-_./]+)$`)

type probe struct {
	Rarity       int
	Ports        *port
	Sslports     *port
	Totalwaitms  time.Duration
	Tcpwrappedms time.Duration
	Request      *request
	MatchGroup   []*match
	Fallback     string

	Response        response
	SoftMatchFilter string
}

func newProbe() *probe {
	return &probe{
		Rarity:       1,
		Totalwaitms:  time.Duration(0),
		Tcpwrappedms: time.Duration(0),

		Ports:      newPort(),
		Sslports:   newPort(),
		Request:    newRequest(),
		MatchGroup: []*match{},
		Fallback:   "",

		Response:        newResponse(),
		SoftMatchFilter: "",
	}
}

func (p *probe) loads(sArr []string) {
	for _, s := range sArr {
		p.load(s)
	}
}

func (p *probe) scan(t target, ssl bool) (string, error) {
	if ssl {
		//fmt.Println("开始TLS探测")
		return simplenet.TLSSend(p.Request.Protocol, t.Uri, p.Request.String, p.Totalwaitms, 512)
	} else {
		//fmt.Println("开始TCP探测")
		return simplenet.Send(p.Request.Protocol, t.Uri, p.Request.String, p.Totalwaitms, 512)
	}
}

func (p *probe) match(s string) TcpFinger {
	var f = newFinger()
	if p.MatchGroup == nil {
		return f
	}
	for _, m := range p.MatchGroup {
		//实现软筛选
		if p.SoftMatchFilter != "" {
			if m.Service != p.SoftMatchFilter {
				continue
			}
		}
		//fmt.Println("开始匹配正则：", m.service, m.patternRegexp.String())
		if m.PatternRegexp.MatchString(s) {
			//fmt.Println("成功匹配指纹：", m.pattern, "所在probe为：", p.request.name)
			if m.Soft {
				//如果为软捕获，这设置筛选器
				f.Service = m.Service
				p.SoftMatchFilter = m.Service
				continue
			} else {
				//如果为硬捕获则直接获取指纹信息
				f = m.makeVersionInfo(s)
				f.Service = m.Service
				return f
			}
		}
	}
	//清空软匹配过滤器
	p.SoftMatchFilter = ""
	return f
}

func (p *probe) load(s string) {
	//分解命令
	i := strings.Index(s, " ")
	commandName := s[:i]
	commandArgs := s[i+1:]
	//逐行处理
	switch commandName {
	case "Probe":
		p.loadProbe(commandArgs)
	case "match":
		p.loadMatch(commandArgs, false)
	case "softmatch":
		p.loadMatch(commandArgs, true)
	case "ports":
		p.loadPorts(commandArgs, false)
	case "sslports":
		p.loadPorts(commandArgs, true)
	case "Totalwaitms":
		p.Totalwaitms = time.Duration(p.getInt(commandArgs)) * time.Millisecond
	case "Tcpwrappedms":
		p.Tcpwrappedms = time.Duration(p.getInt(commandArgs)) * time.Millisecond
	case "Rarity":
		p.Rarity = p.getInt(commandArgs)
	case "Fallback":
		p.Fallback = p.getString(commandArgs)
	}
}

func (p *probe) loadProbe(s string) {
	//Probe <protocol> <probename> <probestring>
	if !PROBE_LOAD_REGEXP.MatchString(s) {
		panic(errors.New("probe 语句格式不正确"))
	}
	args := PROBE_LOAD_REGEXP.FindStringSubmatch(s)
	if args[1] == "" || args[2] == "" {
		panic(errors.New("probe 参数格式不正确"))
	}
	p.Request.Protocol = args[1]
	p.Request.Name = args[1] + "_" + args[2]
	str := args[3]
	str = strings.ReplaceAll(str, `\0`, `\x00`)
	str = strings.ReplaceAll(str, `"`, `${double-quoted}`)
	str = `"` + str + `"`
	str, _ = strconv.Unquote(str)
	str = strings.ReplaceAll(str, `${double-quoted}`, `"`)
	p.Request.String = str
}

func (p *probe) loadMatch(s string, soft bool) {
	m := newMatch()
	//"match": misc.MakeRegexpCompile("^([a-zA-Z0-9-_./]+) m\\|([^|]+)\\|([is]{0,2}) (.*)$"),
	//match <Service> <pattern>|<patternopt> [<versioninfo>]
	//	"matchVersioninfoProductname": misc.MakeRegexpCompile("p/([^/]+)/"),
	//	"matchVersioninfoVersion":     misc.MakeRegexpCompile("v/([^/]+)/"),
	//	"matchVersioninfoInfo":        misc.MakeRegexpCompile("i/([^/]+)/"),
	//	"matchVersioninfoHostname":    misc.MakeRegexpCompile("h/([^/]+)/"),
	//	"matchVersioninfoOS":          misc.MakeRegexpCompile("o/([^/]+)/"),
	//	"matchVersioninfoDevice":      misc.MakeRegexpCompile("d/([^/]+)/"),
	if !m.load(s, soft) {
		panic(errors.New("match 语句参数不正确"))
	}
	p.MatchGroup = append(p.MatchGroup, m)
}

func (p *probe) loadPorts(s string, ssl bool) {
	if ssl {
		if !p.Sslports.LoadS(s) {
			panic(errors.New("sslports 语句参数不正确"))
		}
	} else {
		if !p.Ports.LoadS(s) {
			panic(errors.New("ports 语句参数不正确"))
		}
	}
}

func (p *probe) getInt(expr string) int {
	if !PROBE_INT_REGEXP.MatchString(expr) {
		panic(errors.New("Totalwaitms or Tcpwrappedms 语句参数不正确"))
	}
	i, _ := strconv.Atoi(PROBE_INT_REGEXP.FindStringSubmatch(expr)[1])
	return i
}

func (p *probe) getString(expr string) string {
	if !PROBE_STRING_REGEXP.MatchString(expr) {
		panic(errors.New("Fallback 语句参数不正确"))
	}
	return PROBE_STRING_REGEXP.FindStringSubmatch(expr)[1]
}

func (p *probe) Clean() {
	p.Ports = newPort()
	p.Sslports = newPort()

	p.Request = newRequest()
	p.MatchGroup = []*match{}
	p.Fallback = ""

	p.Response = newResponse()
	p.SoftMatchFilter = ""
}

//
//func (this *probe) Scan(target scan.target) bool {
//	response, err := this.send(target)
//	if err != nil {
//		slog.Debug(err.Error())
//		return false
//	}
//	this.response.string = response
//	return true
//}
//
//func (this *probe) Match() bool {
//	var regular *regexp.Regexp
//	var err error
//	var finger = newFinger()
//	for _, match := range this.MatchGroup {
//		if this.SoftMatchFilter != "" {
//			if match.Service != this.SoftMatchFilter {
//				continue
//			}
//		}
//		regular, err = regexp.Compile(match.pattern)
//		if err != nil {
//			//slog.Debug(fmt.Sprintf("%s:%s",err.Error(),match.pattern))
//			continue
//		}
//		if regular.MatchGrouptring(this.response.string) {
//			if match.soft {
//				//如果为软捕获，这设置筛选器
//				finger.Service = match.Service
//				this.SoftMatchFilter = match.Service
//			} else {
//				//如果为硬捕获则直接设置指纹信息
//				finger = this.makeFinger(regular.FindStringSubmatch(this.response.string), match.versioninfo)
//				finger.Service = match.Service
//				this.response.finger = finger
//				return true
//			}
//		}
//	}
//	if finger.Service != "" {
//		this.response.finger = finger
//		return true
//	} else {
//		return false
//	}
//}
//
//func (this *probe) makeFinger(strArr []string, versioninfo *finger) *finger {
//	versioninfo.info = this.fixFingerValue(versioninfo.info, strArr)
//	versioninfo.devicetype = this.fixFingerValue(versioninfo.devicetype, strArr)
//	versioninfo.hostname = this.fixFingerValue(versioninfo.hostname, strArr)
//	versioninfo.operatingsystem = this.fixFingerValue(versioninfo.operatingsystem, strArr)
//	versioninfo.productname = this.fixFingerValue(versioninfo.productname, strArr)
//	versioninfo.version = this.fixFingerValue(versioninfo.version, strArr)
//	return versioninfo
//}
//
//func (this *probe) fixFingerValue(value string, strArr []string) string {
//	return value
//}
//
//func (this *probe) send(target scan.target) (string, error) {
//	if this.sslports.Len() == 0 && this.ports.Len() == 0 {
//		return stcp.Send(this.request.protocol, target.netloc, this.request.string, this.Tcpwrappedms)
//	}
//	if this.sslports.IsExist(target.port) {
//		return stls.Send(this.request.protocol, target.netloc, this.request.string, this.Tcpwrappedms)
//	}
//	if this.ports.IsExist(target.port) {
//		return stcp.Send(this.request.protocol, target.netloc, this.request.string, this.Tcpwrappedms)
//	}
//	return "", errors.New("无匹配端口，故未进行扫描")
//	//return stcp.Send(this.request.protocol, target.netloc, this.request.string, this.Tcpwrappedms)
//}
