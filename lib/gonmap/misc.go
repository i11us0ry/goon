package gonmap

import (
	"bytes"
	"io"
	"io/ioutil"
)

func IsInIntArr(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func Xrange(args ...int) []int {
	var start, stop int
	var step = 1
	var r []int
	switch len(args) {
	case 1:
		stop = args[0]
		start = 0
	case 2:
		start, stop = args[0], args[1]
	case 3:
		start, step, step = args[0], args[1], args[2]
	default:
		return nil
	}
	if start > stop {
		return nil
	}
	if step < 0 {
		return nil
	}

	for i := start; i <= stop; i += step {
		r = append(r, i)
	}
	return r
}

func CopyIoReader(reader *io.Reader) io.Reader {

	bodyBuf, err := ioutil.ReadAll(*reader)
	if err != nil {
		return nil
	}
	*reader = bytes.NewReader(bodyBuf)
	return bytes.NewReader(bodyBuf)
}
