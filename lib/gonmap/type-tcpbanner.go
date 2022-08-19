package gonmap

type TcpBanner struct {
	Target    target
	Response  response
	TcpFinger TcpFinger
	Status    string
	ErrorMsg  error
}

func NewTcpBanner(target target) TcpBanner {
	return TcpBanner{
		Target:    target,
		Response:  newResponse(),
		TcpFinger: newFinger(),
		Status:    "UNKNOWN",
		ErrorMsg:  nil,
	}
}

func (p *TcpBanner) Load(np *TcpBanner) {
	if p.Status == "CLOSED" || p.Status == "MATCHED" {
		return
	}
	if p.Status == "UNKNOWN" {
		*p = *np
	}
	if p.Status == "OPEN" && np.Status != "CLOSED" && np.Status != "UNKNOWN" {
		*p = *np
	}
	//fmt.Println("加载完成后端口状态为：",p.status)
}

func (p *TcpBanner) Length() int {
	return p.Response.Length()
}

func (p *TcpBanner) CLOSED() *TcpBanner {
	p.Status = "CLOSED"
	return p
}

func (p *TcpBanner) OPEN() *TcpBanner {
	p.Status = "OPEN"
	return p
}

func (p *TcpBanner) MATCHED() *TcpBanner {
	p.Status = "MATCHED"
	return p
}
