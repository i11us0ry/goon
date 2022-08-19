package telnet

import (
	"net"
	"time"
)

type TelnetClient struct {
	Host  			 string
	IsAuthentication bool
	UserName         string
	Password         string
	Time 			 int
}

func (t TelnetClient)Conn() bool{
	return telnet_Creat(t.Host,t.UserName,t.Password,t.Time)
}

func telnet_Creat(host , user ,pass string,time int)  bool {
	telnetClientObj := new(TelnetClient)
	telnetClientObj.Host = host
	telnetClientObj.IsAuthentication = true
	telnetClientObj.UserName = user
	telnetClientObj.Password = pass
	return telnetClientObj.Telnet(time)
}
func (this *TelnetClient) Telnet(timeout int) bool {
	conn, err := net.DialTimeout("tcp", this.Host, time.Duration(timeout)*time.Second)
	if nil != err {
		return false
	}
	defer conn.Close()
	if false == this.telnetProtocolHandshake(conn) {
		return false
	}
	return true
}

func (this *TelnetClient) telnetByWindows(conn net.Conn,buf [8192]byte) bool {
	// 发起telnet连接请求
	buf[0] = 0xff
	buf[1] = 0xfc
	buf[2] = 0x25
	n, err := conn.Write(buf[0:3])
	if nil != err {
		return false
	}
	//fmt.Println("0. ",this.Host, this.UserName, this.Password, n,buf[0:10])
	n, err = conn.Read(buf[0:])
	if nil != err {
		//fmt.Println("1. ",err)
		return false
	}
	//fmt.Println("1. ",this.Host, this.UserName, this.Password, n,buf[0:10])
	// 第二次请求
	buf[0] = 0xff
	buf[1] = 0xfe
	buf[2] = 0x01
	buf[3] = 0xff
	buf[4] = 0xfe
	buf[5] = 0x03
	buf[6] = 0xff
	buf[7] = 0xfc
	buf[8] = 0x27
	buf[9] = 0xff
	buf[10] = 0xfc
	buf[11] = 0x1f
	buf[12] = 0xff
	buf[13] = 0xfc
	buf[14] = 0x00
	buf[15] = 0xff
	buf[16] = 0xfe
	buf[17] = 0x00

	n, err = conn.Write(buf[0:18])
	if nil != err {
		return false
	}
	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	//fmt.Println("2. ",this.Host, this.UserName, this.Password, n,buf[0:10])

	// 输入用户名
	n, err = conn.Write([]byte(this.UserName + "\r\n"))
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * 1000)
	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	//fmt.Println("3. ",this.Host, this.UserName, this.Password, n,buf[0:10])

	// 输入密码
	n, err = conn.Write([]byte(this.Password+ "\r\n"))
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * 1000)
	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	//fmt.Println("4. ",this.Host, this.UserName, this.Password, n,buf[0:10])

	// 登录成功n为3，buf为255 253 24,失败为n为39，buf为13 10 190...
	if n!=3 || (buf[0]!=255 || buf[1]!=253 || buf[2]!=24){
		return false
	}

	//// 执行获取登录成功后页面命令
	//buf[0] = 0xff
	//buf[1] = 0xfc
	//buf[2] = 0x18
	//n, err = conn.Write(buf[0:3])
	//if nil != err {
	//
	//	return false
	//}
	//// 这里读取的是登录成功后的页面
	//n, err = conn.Read(buf[0:])
	//if nil != err {
	//
	//	return false
	//}
	//// 第456位为*==，对应2a 3d 3d
	//if buf[4]!=0x2a && buf[5]!=0x3d && buf[6]!=0x3d{
	//	return false
	//}
	this.closeConn(conn)
	return true
}

func (this *TelnetClient) telnetByLinux(conn net.Conn,buf [8192]byte) bool {
	// 第一次发送数据
	buf[0] = 0xff
	buf[1] = 0xfc
	buf[2] = 0x18
	buf[3] = 0xff
	buf[4] = 0xfc
	buf[5] = 0x20
	buf[6] = 0xff
	buf[7] = 0xfc
	buf[8] = 0x23
	buf[9] = 0xff
	buf[10] = 0xfc
	buf[11] = 0x27
	n, err := conn.Write(buf[0:12])
	if nil != err {
		
		return false
	}
	n, err = conn.Read(buf[0:])
	if nil != err {
		
		return false
	}
	// 第二次发送数据
	buf[0] = 0xff
	buf[1] = 0xfe
	buf[2] = 0x03
	buf[3] = 0xff
	buf[4] = 0xfc
	buf[5] = 0x01
	buf[6] = 0xff
	buf[7] = 0xfc
	buf[8] = 0x1f
	buf[9] = 0xff
	buf[10] = 0xfe
	buf[11] = 0x05
	buf[12] = 0xff
	buf[13] = 0xfc
	buf[14] = 0x21
	n, err = conn.Write(buf[0:15])
	if nil != err {
		
		return false
	}
	n, err = conn.Read(buf[0:])
	if nil != err {
		
		return false
	}

	// 第三次发送数据，获取系统信息
	buf[0] = 0xff
	buf[1] = 0xfe
	buf[2] = 0x03
	n, err = conn.Write(buf[0:3])
	if nil != err {
		
		return false
	}
	n, err = conn.Read(buf[0:])
	if nil != err {
		
		return false
	}

	// 第四次发送数据，获取登录点
	buf[0] = 0xff
	buf[1] = 0xfe
	buf[2] = 0x01
	n, err = conn.Write(buf[0:3])
	if nil != err {
		
		return false
	}
	n, err = conn.Read(buf[0:])
	if nil != err {
		
		return false
	}

	// 输入用户名
	n, err = conn.Write([]byte(this.UserName + "\n"))
	if nil != err {
		
		return false
	}
	time.Sleep(time.Millisecond * 1000)

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}

	// 输入密码
	n, err = conn.Write([]byte(this.Password+ "\n"))
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * 1000)
	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}

	//fmt.Println(this.Host,this.UserName,this.Password,n,buf[0], buf[1], buf[2], buf[0:10])
	// 成功标志n=,buf=[13 10 87 101 108 99...],失败标志n=2,buf=[13 10 13 10 80 97...]
	if n==2 || (buf[0]!=13 || buf[1]!=10 || buf[2]!=87){
		//fmt.Println("失败，返回")
		return false
	}

	//// 执行exit
	buf[0] = 0x65
	buf[1] = 0x78
	buf[2] = 0x69
	buf[3] = 0x74
	buf[4] = 0x0a
	conn.Write(buf[0:5])
	this.closeConn(conn)
	return true
}

func (this *TelnetClient) closeConn(conn net.Conn) {
	var buf [8192]byte
	// 执行exit
	//fmt.Println("发送exit.")
	buf[0] = 0x65
	buf[1] = 0x78
	buf[2] = 0x69
	buf[3] = 0x74
	buf[4] = 0x0a
	//n, err = conn.Write(buf[0:5])
	//if nil != err {
	//
	//	return false
	//}
	conn.Write(buf[0:5])
}

func (this *TelnetClient) telnetProtocolHandshake(conn net.Conn) bool {
	var buf [8192]byte
	_, err := conn.Read(buf[0:])
	if nil != err {
		
		return false
	}
	// 判断系统类型,0x18为linux系统，0x25为win系统
	//fmt.Println(this.Host, this.UserName, this.Password, buf[2])
	if buf[2] == 37 || buf[2] == 178 {
		return this.telnetByWindows(conn,buf)
	} else if buf[2] == 24 {
		return this.telnetByLinux(conn,buf)
	}
	
	return false
}