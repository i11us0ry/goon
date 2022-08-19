package gonmap

type request struct {
	//Probe <protocol> <probename> <probestring>
	Protocol string
	Name     string
	String   string
}

func newRequest() *request {
	return &request{
		Protocol: "",
		Name:     "",
		String:   "",
	}
}
