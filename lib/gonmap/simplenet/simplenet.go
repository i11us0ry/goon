package simplenet

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"strings"
	"time"
)

func Send(protocol string, netloc string, data string, duration time.Duration, size int) (string, error) {
	protocol = strings.ToLower(protocol)
	conn, err := net.DialTimeout(protocol, netloc, duration)
	if err != nil {
		//fmt.Println(conn)
		return "", errors.New(err.Error() + " STEP1:CONNECT")
	}
	buf := make([]byte, size)
	_, err = io.WriteString(conn, data)
	//_, err = io.WriteString(conn, data)
	if err != nil {
		_ = conn.Close()
		return "", errors.New(err.Error() + " STEP2:WRITE")
	}
	//设置读取超时Deadline
	_ = conn.SetReadDeadline(time.Now().Add(duration * 2))
	length, err := conn.Read(buf)
	if err != nil && err.Error() != "EOF" {
		_ = conn.Close()
		return "", errors.New(err.Error() + " STEP3:READ")
	}
	_ = conn.Close()
	if length == 0 {
		return "", errors.New("response is empty")
	}
	return string(buf[:length]), nil
}

func TLSSend(protocol string, netloc string, data string, duration time.Duration, size int) (string, error) {
	protocol = strings.ToLower(protocol)
	config := &tls.Config{
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS10,
	}
	dial := &net.Dialer{
		Timeout:  duration,
		Deadline: time.Now().Add(duration * 2),
	}
	conn, err := tls.DialWithDialer(dial, protocol, netloc, config)
	if err != nil {
		return "", errors.New(err.Error() + " STEP1:CONNECT")
	}
	_, err = io.WriteString(conn, data)
	if err != nil {
		_ = conn.Close()
		return "", errors.New(err.Error() + " STEP2:WRITE")
	}
	buf := make([]byte, size)
	_ = conn.SetReadDeadline(time.Now().Add(duration * 2))
	length, err := conn.Read(buf)
	if err != nil && err.Error() != " EOF" {
		_ = conn.Close()
		return "", errors.New(err.Error() + " STEP3:READ")
	}
	_ = conn.Close()
	if length == 0 {
		return "", errors.New("response is empty")
	}
	return string(buf[:length]), nil
}
