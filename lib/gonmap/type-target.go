package gonmap

type target struct {
	Port int
	Host string
	Uri  string
}

func newTarget() target {
	return target{0, "", ""}
}
