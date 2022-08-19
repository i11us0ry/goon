package scan

import (
	"goon3/public"
	"net"
	"runtime"
	"time"
)

func Icmp (ips []string){
	runtime.GOMAXPROCS(runtime.NumCPU())
	input := make(chan string, len(ips))
	result := make(chan string, len(ips))
	defer close(input)
	for _, ip := range(ips){
		input <- ip
	}
	thread := 10
	if len(ips) < Par.Thread {
		thread = len(ips)
	} else {
		thread = Par.Thread
	}
	for i := 0; i< thread; i++{
		go icmpWork(input,result)
	}
	public.Out(result,Par.Ofile)
}

func icmpWork(input,result chan string){
	for {
		task,ok := <-input
		if !ok{
			return
		}
		alive := IcmpScan(task, Par.Timeout)
		if alive{
			result <- "[icmp] " + task
			Par.IpAlive = append(Par.IpAlive,task)
		} else {
			result <- ""
		}
	}
}

func IcmpScan(host string,t int) bool{
	var size int
	var seq int16 = 1
	const ECHO_REQUEST_HEAD_LEN = 8

	size = 32
	starttime := time.Now()
	conn, err := net.DialTimeout("ip4:icmp", host, time.Duration(t) * time.Second)
	if err != nil {
		return false
	}
	defer conn.Close()
	id0, id1 := genidentifier(host)

	var msg []byte = make([]byte, size+ECHO_REQUEST_HEAD_LEN)
	msg[0] = 8                        // echo
	msg[1] = 0                        // code 0
	msg[2] = 0                        // checksum
	msg[3] = 0                        // checksum
	msg[4], msg[5] = id0, id1         //identifier[0] identifier[1]
	msg[6], msg[7] = gensequence(seq) //sequence[0], sequence[1]

	length := size + ECHO_REQUEST_HEAD_LEN

	check := checkSum(msg[0:length])
	msg[2] = byte(check >> 8)
	msg[3] = byte(check & 255)

	conn.SetDeadline(starttime.Add(time.Duration(t) * time.Second))
	_, err = conn.Write(msg[0:length])

	const ECHO_REPLY_HEAD_LEN = 20

	var receive []byte = make([]byte, ECHO_REPLY_HEAD_LEN+length)
	n, err := conn.Read(receive)
	_ = n
	var endduration int = int(int64(time.Since(starttime)) / (1000 * 1000))

	if err != nil || receive[ECHO_REPLY_HEAD_LEN+4] != msg[4] || receive[ECHO_REPLY_HEAD_LEN+5] != msg[5] || receive[ECHO_REPLY_HEAD_LEN+6] != msg[6] || receive[ECHO_REPLY_HEAD_LEN+7] != msg[7] || endduration >= int(1000) || receive[ECHO_REPLY_HEAD_LEN] == 11 {
		return false
	}
	return true
}

/* 检验和 */
func checkSum(msg []byte) uint16 {
	sum := 0
	length := len(msg)
	for i := 0; i < length-1; i += 2 {
		sum += int(msg[i])*256 + int(msg[i+1])
	}
	if length%2 == 1 {
		sum += int(msg[length-1]) * 256 // notice here, why *256?
	}

	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	var answer uint16 = uint16(^sum)
	return answer
}


func gensequence(v int16) (byte, byte) {
	ret1 := byte(v >> 8)
	ret2 := byte(v & 255)
	return ret1, ret2
}

func genidentifier(host string) (byte, byte) {
	return host[0], host[1]
}