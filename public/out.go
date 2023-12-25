package public

var Thread int

func Out(result chan string) {
	/*
		扫描进度
	*/
	i, y := 0, 0
	for host := range result {
		i++
		if host != "" {
			y++
			Success.Printf("%s", host)
			if !InputValue.NoOutputPtr {
				FileWrite(InputValue.OfilePtr, (host + "\n"))
			}
		}
		if i == cap(result) {
			close(result)
			break
		}
	}
}

func Out2(s string) {
	if !InputValue.NoOutputPtr {
		FileWrite(InputValue.OfilePtr, s)
	}
}
