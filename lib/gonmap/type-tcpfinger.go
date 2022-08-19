package gonmap

type TcpFinger struct {
	Service         string
	ProductName     string
	Version         string
	Info            string
	Hostname        string
	OperatingSystem string
	DeviceType      string
	//  p/vendorproductname/
	//	v/version/
	//	i/info/
	//	h/hostname/
	//	o/operatingsystem/
	//	d/devicetype/
}

func newFinger() TcpFinger {
	return TcpFinger{
		Service:         "",
		ProductName:     "",
		Version:         "",
		Info:            "",
		Hostname:        "",
		OperatingSystem: "",
		DeviceType:      "",
	}
}
