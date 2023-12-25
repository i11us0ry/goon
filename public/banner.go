package public

import (
	"fmt"
)

func Banner() {
	var goon = "" +
		"---------------------------------------------------------------------------------------------\n" +
		"				  ____  ____  ____  ____ \n" +
		"				 / __ `/ __ \\/ __ \\/ __ \\\n" +
		"				/ /_/ / /_/ / /_/ / / / /\n" +
		"				\\__, /\\____/\\____/_/ /_/\n" +
		"				/___/\n\n" +
		`
					goon v3
					by:i11us0ry
---------------------------------------------------------------------------------------------                                                    
	`
	fmt.Println(goon)
}

func OutMode() {
	var outMode = `
可选mode如下:

all:		默认选项, 输入支持ip、domain（自动转换为url）、url, 流程包含ip-port(web)-title-finger-ftp-ms17010-mssql-mysql-postgres-redis-ssh-smb-rdp-telnet-netbios
webscan:	包含ip-port(web)-title-finger
brute:		包含ip-ftp-ms17010-mssql-mysql-postgres-redis-ssh-smb-rdp-telnet
ip:		ip探活, 默认使用icmp，执行-ping使用ping，支持/8-/31之间任意CIDR，/8-/15之间自动生成所有c段，先探测每个c段的.1;/16-/23之间自动生成所有c段，先探测每个c段的.1和.254，/24先探测.1和.24，/25-/31探测所有ip
port:		端口扫描,执行-web直接探测http/https
fofa:		fofa资产获取,执行-web输出host时添加http(fields为多个时host放在最后一位)
title:		title扫描，输入支持 domain、url
finger:		web指纹探测，输入支持 domain、url
dfuzz:		路径fuzz,适用于对批量url进行单个dir探测，支持post发包，支持正则匹配，可探测简单poc
tomcat:		tomcat爆破，输入支持 url
ftp:		ftp爆破, 输入支持ip、domain(127.0.0.1:21 21端口被指定为爆破端口),其他ms17010,mssql,mysql,postgres,redis,ssh,smb,rdp,telnet同理
netbios:	netbios探测

其他：
ip：127.0.0.1
domain：127.0.0.1:8080
url：http://127.0.0.1:8080

详细参考：-h
`
	fmt.Println(outMode)
}
