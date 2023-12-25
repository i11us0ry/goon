package check

import (
	"fmt"
	"goon3/module/brute"
	"goon3/module/scan"
	"goon3/public"
	"os"
	"strings"
)

/*
处理mode，并进行任务下发
*/
func CheckMode() {
	fmt.Println()
	public.Info.Println("checking mode……")
	fmt.Println()
	CheckModeStart()
}

func CheckModeStart() {
	/*
		先将ConfValue值传给扫描和爆破模块
	*/
	scan.Init(&public.ConfValue)
	brute.Init(&public.ConfValue)

	tempIps := public.ConfValue.Ip
	/*
		ip探活，如果mode为ip则直接返回结果
	*/
	if public.InputValue.PingPtr && len(public.ConfValue.Ip) != 0 {
		PingScan(public.ConfValue.Ip)
		if len(public.ConfValue.IpAlive) <= 0 {
			public.Warning.Println("ping cant find any ip 2!")
			os.Exit(0)
		} else {
			tempIps = public.ConfValue.IpAlive
		}
		if public.InputValue.ModePtr == "ip" {
			return
		}
	} else if !public.InputValue.NoPingPtr && len(public.ConfValue.Ip) != 0 {
		/*
		 如果输入Ping则调用ping探活，反之用icmp
		*/
		IcmpScan(public.ConfValue.Ip)
		if len(public.ConfValue.IpAlive) <= 0 {
			public.Warning.Println("icmp cant find any ip 2!")
			os.Exit(0)
		} else {
			tempIps = public.ConfValue.IpAlive
		}
		if public.InputValue.ModePtr == "ip" {
			return
		}
	}

	switch public.InputValue.ModePtr {
	case "port":
		PortScan(tempIps)
	case "title":
		TitleScan(public.ConfValue.Url)
	case "finger":
		FingerScan(public.ConfValue.Url)
	case "tomcat":
		TomcatScan(public.ConfValue.Url)
	case "dfuzz":
		DirScan(public.ConfValue.Url)
	case "fofa":
		Fofa(public.ConfValue.FofaWord)
	case "ftp":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Ftp.Info.Port)
		Ftp(bruteHost.Ftp.Host)
	case "mssql":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Mssql.Info.Port)
		Mssql(bruteHost.Mssql.Host)
	case "mysql":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Mysql.Info.Port)
		Mysql(bruteHost.Mysql.Host)
	case "postgres":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Postgres.Info.Port)
		Postgres(bruteHost.Postgres.Host)
	case "redis":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Redis.Info.Port)
		Redis(bruteHost.Redis.Host)
	case "ssh":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Ssh.Info.Port)
		Ssh(bruteHost.Ssh.Host)
	case "webscan":
		public.InputValue.WebPtr = true
		PortScan(tempIps)
		TitleScan(public.ConfValue.Url)
		FingerScan(public.ConfValue.Url)
	case "brute":
		bruteHost := PortInfo(tempIps, nil)
		Ftp(bruteHost.Ftp.Host)
		Mssql(bruteHost.Mssql.Host)
		Mysql(bruteHost.Mysql.Host)
		Postgres(bruteHost.Postgres.Host)
		Redis(bruteHost.Redis.Host)
		Ssh(bruteHost.Ssh.Host)
		SMB(bruteHost.Smb.Host)
		RDP(bruteHost.Rdp.Host)
	case "ms17010":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Smb.Info.Port)
		MS17010(bruteHost.Smb.Host)
	case "smb":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Smb.Info.Port)
		SMB(bruteHost.Smb.Host)
	case "netbios":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.NetBios.Info.Port)
		NetBios(bruteHost.NetBios.Host)
	case "rdp":
		bruteHost := PortInfo(tempIps, public.ConfValue.Brute.Rdp.Info.Port)
		RDP(bruteHost.Rdp.Host)
	case "all":
		public.InputValue.WebPtr = true
		PortScan(tempIps)
		TitleScan(public.ConfValue.Url)
		FingerScan(public.ConfValue.Url)
		bruteHost := PortInfo(tempIps, nil)
		Ftp(bruteHost.Ftp.Host)
		Mssql(bruteHost.Mssql.Host)
		Mysql(bruteHost.Mysql.Host)
		Postgres(bruteHost.Postgres.Host)
		Redis(bruteHost.Redis.Host)
		Ssh(bruteHost.Ssh.Host)
		SMB(bruteHost.Smb.Host)
		MS17010(bruteHost.Smb.Host)
		RDP(bruteHost.Rdp.Host)
		NetBios(bruteHost.NetBios.Host)
	default:
		public.OutMode()
	}
}

func PingScan(ips []string) {
	public.Info.Println("start ping scan……")
	public.Out2("\n------------------------------------ping------------------------------------\n")
	fmt.Println()
	scan.Ping(ips)
}

func IcmpScan(ips []string) {
	public.Info.Println("start icmp scan……")
	public.Out2("\n------------------------------------icmp------------------------------------\n")
	fmt.Println()
	scan.Icmp(ips)
}

func PortScan(ips []string) {
	if len(ips) > 0 {
		fmt.Println()
		public.Info.Println("start port scan……")
		public.Out2("\n------------------------------------port------------------------------------\n")
		fmt.Println()
		scan.PortScan(ips)
	}
}

func TitleScan(urls []string) {
	if len(urls) > 0 {
		fmt.Println()
		public.Info.Println("start title scan……")
		public.Out2("\n------------------------------------title------------------------------------\n")
		fmt.Println()
		scan.TitleScan(urls)
	}
}

func FingerScan(urls []string) {
	if len(urls) > 0 {
		fmt.Println()
		public.Info.Println("start finger scan……")
		public.Out2("\n------------------------------------finger------------------------------------\n")
		fmt.Println()
		scan.FingerScan(urls)
	}
}

func TomcatScan(urls []string) {
	if len(urls) > 0 {
		public.Info.Println("start tomcat brute……")
		public.Out2("\n------------------------------------tomcat------------------------------------\n")
		fmt.Println()
		scan.TomcatScan(urls)
	}
}

func DirScan(urls []string) {
	if len(urls) > 0 {
		public.Info.Println("start dir scan……")
		public.Out2("\n------------------------------------dir------------------------------------\n")
		fmt.Println()
		scan.DirScan(urls)
	}
}

func Fofa(words []string) {
	if len(words) > 0 {
		public.Info.Println("start fofa scan……")
		public.Out2("\n------------------------------------fofa------------------------------------\n")
		fmt.Println()
		scan.FofaScan(words)
	}
}

func Ftp(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		public.Info.Println("start ftp brute……")
		public.Out2("\n------------------------------------ftp------------------------------------\n")
		fmt.Println()
		brute.Ftp(hosts)
	}
}

func Mssql(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start mssql brute……")
		public.Out2("\n------------------------------------mssql------------------------------------\n")
		fmt.Println()
		brute.Mssql(hosts)
	}
}

func Mysql(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start mysql brute……")
		public.Out2("\n------------------------------------mysql------------------------------------\n")
		fmt.Println()
		brute.Mysql(hosts)
	}
}

func Postgres(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start postgres brute……")
		public.Out2("\n------------------------------------postgres------------------------------------\n")
		fmt.Println()
		brute.Postgres(hosts)
	}
}

func Redis(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start redis brute……")
		public.Out2("\n------------------------------------redis------------------------------------\n")
		fmt.Println()
		brute.Redis(hosts)
	}
}

func Ssh(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start ssh brute……")
		public.Out2("\n------------------------------------ssh------------------------------------\n")
		fmt.Println()
		brute.Ssh(hosts)
	}
}

func MS17010(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start MS17010探测……")
		public.Out2("\n------------------------------------ms17010------------------------------------\n")
		fmt.Println()
		brute.MS17010(hosts)
	}
}

func SMB(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start SMB探测……")
		public.Out2("\n------------------------------------smb------------------------------------\n")
		fmt.Println()
		brute.Smb(hosts)
	}
}

func NetBios(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start netbios scan……")
		public.Out2("\n------------------------------------netbios------------------------------------\n")
		fmt.Println()
		scan.NetBIOS(hosts)
	}
}

func RDP(hosts []string) {
	if len(hosts) > 0 {
		// 这里得到的原始ip或存活检测过后的ip，下一步判断是否做指纹识别
		fmt.Println()
		public.Info.Println("start rdp brute……")
		public.Out2("\n------------------------------------rdp------------------------------------\n")
		fmt.Println()
		brute.RDP(hosts)
	}
}

func PortInfo(ips, ports []string) *public.BruteHostType {
	if len(ips) > 0 {
		fmt.Println()
		public.Info.Println("start finger scan……")
		fmt.Println()
		hosts := []string{}
		// 如果单独输入了port则使用输入的port而不是用配置文件中的默认port
		if public.InputValue.PortPtr != "" {
			ports = public.ConfValue.Port
			for _, ip := range ips {
				if strings.Contains(ip, ":") {
					hosts = append(hosts, ip)
				} else {
					for _, port := range ports {
						hosts = append(hosts, fmt.Sprintf("%v:%v", ip, port))
					}
				}
			}
		} else if ports == nil {
			// brute模式下使用所有配置文件中的端口
			ports = getAllPort()
			for _, ip := range ips {
				if strings.Contains(ip, ":") {
					hosts = append(hosts, ip)
				} else {
					for _, port := range ports {
						hosts = append(hosts, fmt.Sprintf("%v:%v", ip, port))
					}
				}
			}
		} else {
			for _, ip := range ips {
				if strings.Contains(ip, ":") {
					hosts = append(hosts, ip)
				} else {
					for _, port := range ports {
						hosts = append(hosts, fmt.Sprintf("%v:%v", ip, port))
					}
				}
			}
		}
		return public.GetPortInfo(hosts, public.ConfValue.Thread)
	}
	return nil
}

// 在brute模式下，如果没有指定端口，则从配置文件获取所有默认端口
func getAllPort() []string {
	ports := []string{}
	for _, port := range public.ConfValue.Brute.Ftp.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.Mssql.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.Mysql.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.Postgres.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.Redis.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.Ssh.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.Smb.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.Rdp.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.Telnet.Info.Port {
		ports = append(ports, port)
	}
	for _, port := range public.ConfValue.Brute.NetBios.Info.Port {
		ports = append(ports, port)
	}
	return ports
}
