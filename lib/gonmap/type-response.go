package gonmap

type response struct {
	String string
}

func newResponse() response {
	return response{
		String: "",
	}
}

func (r response) Length() int {
	return len(r.String)
}
