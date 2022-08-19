package gonmap

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var NMAP *Nmap

//r["PROBE"] 总探针数、r["MATCH"] 总指纹数 、r["USED_PROBE"] 已使用探针数、r["USED_MATCH"] 已使用指纹数
func Init(filter int, timeout time.Duration) map[string]int {
	//初始化NMAP探针库
	InitNMAP()
	//fmt.Println("初始化了")
	NMAP = &Nmap{
		Exclude:     newPort(),
		ProbeGroup:  make(map[string]*probe),
		ProbeSort:   []string{},
		PortMap:     make(map[int][]string),
		AllPortMap:  []string{},
		ProbeFilter: 0,
		Target:      newTarget(),
		Response:    newResponse(),
		Finger:      nil,
		Filter:      5,
	}
	NMAP.Filter = filter
	for i := 6370 ; i <= 6379; i++ {
		NMAP.PortMap[i] = []string{}
	}
	NMAP.loads(NMAP_SERVICE_PROBES)
	NMAP.AddAllProbe("TCP_GetRequest")
	NMAP.setTimeout(timeout)
	NMAP.AddMatch("TCP_GetRequest", `http m|^HTTP/1\.[01] \d\d\d (?:[^\r\n]+\r\n)*?Server: ([^\r\n]+)| p/$1/`)
	NMAP.AddMatch("TCP_GetRequest", `http m|^HTTP/1\.[01] \d\d\d|`)
	return NMAP.Status()
}

func InitNMAP() {
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, "${backquote}", "`")
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `\1`, `$1`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?=\\)`, `(?:\\)`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?=[\w._-]{5,15}\r?\n$)`, `(?:[\w._-]{5,15}\r?\n$)`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?:[^\r\n]*r\n(?!\r\n))*?`, `(?:[^\r\n]+\r\n)*?`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?:[^\r\n]*\r\n(?!\r\n))*?`, `(?:[^\r\n]+\r\n)*?`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?:[^\r\n]+\r\n(?!\r\n))*?`, `(?:[^\r\n]+\r\n)*?`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?!2526)`, ``)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?!400)`, ``)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?!\0\0)`, ``)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?!/head>)`, ``)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?!HTTP|RTSP|SIP)`, ``)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?!.*[sS][sS][hH]).*`, `.*`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?!\xff)`, `.`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?!x)`, `[^x]`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?<=.)`, `(?:.)`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `(?<=\?)`, `(?:\?)`)
	NMAP_SERVICE_PROBES = strings.ReplaceAll(NMAP_SERVICE_PROBES, `match rtmp`, `# match rtmp`)
}

func New() *Nmap {
	n := &Nmap{}
	*n = *NMAP
	return n
}

type Nmap struct {
	Exclude     *port
	ProbeGroup  map[string]*probe
	ProbeSort   []string
	ProbeFilter int
	PortMap     map[int][]string
	AllPortMap  []string

	Target target
	Filter int

	Response response
	Finger   *TcpFinger
}

func (n *Nmap) Scan(ip string, port int) TcpBanner {
	n.Target.Host = ip
	n.Target.Port = port
	n.Target.Uri = fmt.Sprintf("%s:%d", ip, port)

	//fmt.Println(n.portMap[port])
	//拼接端口探测队列，全端口探测放在最后
	b := NewTcpBanner(n.Target)
	//开始特定端口探测
	for _, requestName := range n.PortMap[port] {
		tls := n.ProbeGroup[requestName].Sslports.Exist(n.Target.Port)
		//fmt.Println(tls)
		//fmt.Println("开始探测：", requestName, "权重为", tls,n.probeGroup[requestName].rarity)
		b.Load(n.getTcpBanner(n.ProbeGroup[requestName], tls))
		if b.Status == "CLOSED" || b.Status == "MATCHED" {
			break
		}

		if n.Target.Port == 53 {
			if DnsScan(n.Target.Uri) {
				b.TcpFinger.Service = "dns"
				b.Response.String = "dns"
				b.MATCHED()
			} else {
				b.CLOSED()
			}
			break
		}

	}
	//fmt.Println(b.status)
	if b.Status != "MATCHED" && b.Status != "CLOSED" {
		//开始全端口探测
		for _, requestName := range n.AllPortMap {
			//fmt.Println("开始全端口探测：", requestName, "权重为", n.probeGroup[requestName].rarity)
			b.Load(n.getTcpBanner(n.ProbeGroup[requestName], false))
			if b.Status == "CLOSED" || b.Status == "MATCHED" {
				break
			}
			b.Load(n.getTcpBanner(n.ProbeGroup[requestName], true))
			if b.Status == "CLOSED" || b.Status == "MATCHED" {
				break
			}
		}
	}
	//进行最后输出修饰
	if b.TcpFinger.Service == "ssl/http" {
		b.TcpFinger.Service = "https"
	}
	if b.TcpFinger.Service == "ssl/https" {
		b.TcpFinger.Service = "https"
	}
	if b.TcpFinger.Service == "ms-wbt-server" {
		b.TcpFinger.Service = "rdp"
	}
	if b.TcpFinger.Service == "ssl" && n.Target.Port == 3389 {
		b.TcpFinger.Service = "rdp"
	}
	return b
}

func (n *Nmap) getTcpBanner(p *probe, tls bool) *TcpBanner {
	b := NewTcpBanner(n.Target)
	//fmt.Println("开始发送数据:",p.request.name,"超时时间为：",p.totalwaitms,p.tcpwrappedms)
	data, err := p.scan(n.Target, tls)
	//fmt.Println(data,err)
	if err != nil {
		b.ErrorMsg = err
		if strings.Contains(err.Error(), "STEP1") {
			return b.CLOSED()
		}
		//if p.request.protocol == "UDP" {
		//	return b.CLOSED()
		//}
		return b.OPEN()
	} else {
		b.Response.String = data
		//若存在返回包，则开始捕获指纹
		//fmt.Printf("成功捕获到返回包，返回包为：%x\n", data)
		//fmt.Printf("成功捕获到返回包，返回包长度为：%x\n", len(data))
		b.TcpFinger = n.getFinger(data, p.Request.Name)
		if b.TcpFinger.Service == "" {
			return b.OPEN()
		} else {
			if tls {
				if b.TcpFinger.Service == "http" {
					b.TcpFinger.Service = "https"
				}
			}
			return b.MATCHED()
		}
		//如果成功匹配指纹，则直接返回指纹
	}
}

func (n *Nmap) AddMatch(probeName string, expr string) {
	n.ProbeGroup[probeName].loadMatch(expr, false)
}

func (n *Nmap) AddAllProbe(probeName string) {
	n.AllPortMap = append(n.AllPortMap, probeName)
}

func (n *Nmap) filter(i int) {
	n.Filter = i
}

func (n *Nmap) Status() map[string]int {
	r := make(map[string]int)
	r["PROBE"] = len(NMAP.ProbeSort)
	r["MATCH"] = 0
	for _, p := range NMAP.ProbeGroup {
		r["MATCH"] += len(p.MatchGroup)
	}
	//fmt.Printf("成功加载探针：【%d】个,指纹【%d】条\n", PROBE_COUNT,MATCH_COUNT)
	r["USED_PROBE"] = len(NMAP.PortMap[0])
	r["USED_MATCH"] = 0
	for _, p := range NMAP.PortMap[0] {
		r["USED_MATCH"] += len(NMAP.ProbeGroup[p].MatchGroup)
	}
	//fmt.Printf("本次扫描将使用探针:[%d]个,指纹[%d]条\n", USED_PROBE_COUNT,USED_MATCH_COUNT)
	return r
}

func (n *Nmap) setTimeout(timeout time.Duration) {
	if timeout == 0 {
		return
	}
	for _, p := range n.ProbeGroup {
		p.Totalwaitms = timeout
		p.Tcpwrappedms = timeout
	}
}

func (n *Nmap) isCommand(line string) bool {
	//删除注释行和空行
	if len(line) < 2 {
		return false
	}
	if line[:1] == "#" {
		return false
	}
	//删除异常命令
	commandName := line[:strings.Index(line, " ")]
	commandArr := []string{
		"Exclude", "Probe", "match", "softmatch", "ports", "sslports", "totalwaitms", "tcpwrappedms", "rarity", "fallback",
	}
	for _, item := range commandArr {
		if item == commandName {
			return true
		}
	}
	return false
}

func (n *Nmap) getFinger(data string, requestName string) TcpFinger {
	data = n.convResponse(data)
	//fmt.Println(data)
	f := n.ProbeGroup[requestName].match(data)
	if f.Service == "" {
		if n.ProbeGroup[requestName].Fallback != "" {
			return n.ProbeGroup["TCP_"+n.ProbeGroup[requestName].Fallback].match(data)
		}
	}
	return f
}

func (n *Nmap) convResponse(s1 string) string {
	//	为了适配go语言的沙雕正则，只能讲二进制强行转换成UTF-8
	b1 := []byte(s1)
	var r1 []rune
	for _, i := range b1 {
		r1 = append(r1, rune(i))
	}
	s2 := string(r1)
	return s2
}

func (n *Nmap) loads(s string) {
	lines := strings.Split(s, "\n")
	var probeArr []string
	p := newProbe()
	for _, line := range lines {
		if !n.isCommand(line) {
			continue
		}
		commandName := line[:strings.Index(line, " ")]
		if commandName == "Exclude" {
			n.loadExclude(line)
			continue
		}
		if commandName == "Probe" {
			if len(probeArr) != 0 {
				p.loads(probeArr)
				n.pushProbe(p)
			}
			probeArr = []string{}
			p.Clean()
		}
		probeArr = append(probeArr, line)
	}
	p.loads(probeArr)
	n.pushProbe(p)
}

func (n *Nmap) loadExclude(expr string) {
	var exclude = newPort()
	expr = expr[strings.Index(expr, " ")+1:]
	for _, s := range strings.Split(expr, ",") {
		if !exclude.Load(s) {
			panic(errors.New("exclude 语句格式错误"))
		}
	}
	n.Exclude = exclude
}

func (n *Nmap) pushProbe(p *probe) {
	PROBE := newProbe()
	*PROBE = *p
	//if p.ports.length == 0 && p.sslports.length == 0 {
	//	fmt.Println(p.request.name)
	//}
	n.ProbeSort = append(n.ProbeSort, p.Request.Name)
	n.ProbeGroup[p.Request.Name] = PROBE

	//建立端口扫描对应表，将根据端口号决定使用何种请求包
	//如果端口列表为空，则为全端口
	if p.Rarity > n.Filter {
		return
	}
	//0记录所有使用的探针
	n.PortMap[0] = append(n.PortMap[0], p.Request.Name)

	if p.Ports.Length+p.Sslports.Length == 0 {
		p.Ports.Fill()
		p.Sslports.Fill()
		n.AllPortMap = append(n.AllPortMap, p.Request.Name)
		return
	}
	//分别压入sslports,ports
	for _, i := range p.Ports.Value {
		n.PortMap[i] = append(n.PortMap[i], p.Request.Name)
	}
	for _, i := range p.Sslports.Value {
		n.PortMap[i] = append(n.PortMap[i], p.Request.Name)
	}

}
