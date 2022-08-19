package check

import (
	"fmt"
	"goon3/public"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func GetConfigDir() string{
	dir := public.GetCurrentDir()
	return fmt.Sprintf("%v/conf.yml",dir)
}

/* 检查配置文件 */
func ConfigCheck(){
	config_Dir := GetConfigDir()
	_, exist := os.Stat(config_Dir)
	// 文件不存在
	if os.IsNotExist(exist) {
		fmt.Println()
		public.Warning.Println("conf.yml is created, please run it again！")
		configWrite(config_Dir)
		os.Exit(0)
	}
}

func configWrite(fileName string){
	/* portscan */
	conf := &public.Conf{}
	conf.Thread = 500
	conf.Timeout = 3
	conf.Follow = true

	// 端口扫描
	conf.PortScan.Ports = public.WebPort
	conf.PortScan.StatusCode = []int{200,302,401}

	/* dirscan */
	conf.DirScan.Dir = ""
	conf.DirScan.Mode = "get"
	conf.DirScan.Code = 200
	conf.DirScan.Body = ""
	conf.DirScan.Header = []string{"Accept-Language:zh-CN,zh;q=0.9","Accept:*/*"}
	conf.DirScan.RBody = ""
	conf.DirScan.RHeader = ""

	/* fofascan */
	conf.FofaScan.Email = ""
	conf.FofaScan.Key = ""
	conf.FofaScan.Num = 100
	conf.FofaScan.Fields = "host"

	/* brute */
	conf.Brute.Thread = 500
	conf.Brute.Timeout = 30
	conf.Brute.Redis.Pass = public.RedisDic["pass"]
	conf.Brute.Redis.Port = public.RedisDic["port"]

	conf.Brute.Ssh.User = public.SshDic["user"]
	conf.Brute.Ssh.Pass = public.SshDic["pass"]
	conf.Brute.Ssh.Port = public.SshDic["port"]

	conf.Brute.Ftp.User = public.FtpDic["user"]
	conf.Brute.Ftp.Pass = public.FtpDic["pass"]
	conf.Brute.Ftp.Port = public.FtpDic["port"]

	conf.Brute.Mysql.User = public.MysqlDic["user"]
	conf.Brute.Mysql.Pass = public.MysqlDic["pass"]
	conf.Brute.Mysql.Port = public.MysqlDic["port"]

	conf.Brute.Mssql.User = public.MssqlDic["user"]
	conf.Brute.Mssql.Pass = public.MssqlDic["pass"]
	conf.Brute.Mssql.Port = public.MssqlDic["port"]

	conf.Brute.Postgres.User = public.PostgreDic["user"]
	conf.Brute.Postgres.Pass = public.PostgreDic["pass"]
	conf.Brute.Postgres.Port = public.PostgreDic["port"]

	conf.Brute.Smb.User = public.SmbDic["user"]
	conf.Brute.Smb.Pass = public.SmbDic["pass"]
	conf.Brute.Smb.Port = public.SmbDic["port"]

	conf.Brute.Rdp.User = public.RdpDic["user"]
	conf.Brute.Rdp.Pass = public.RdpDic["pass"]
	conf.Brute.Rdp.Port = public.RdpDic["port"]

	conf.Brute.NetBios.Port = public.NetBiosPort

	conf.Brute.Telnet.User = public.TelnetDic["user"]
	conf.Brute.Telnet.Pass = public.TelnetDic["pass"]
	conf.Brute.Telnet.Port = public.TelnetDic["port"]

	conf.Brute.Tomcat.User = public.TomcatDic["user"]
	conf.Brute.Tomcat.Pass = public.TomcatDic["pass"]

	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	defer file.Close()
	enc := yaml.NewEncoder(file)
	err := enc.Encode(conf)
	if err != nil {
		public.Error.Printf("%v" ,err)
		os.Exit(1)
	}
}

func ConfigRead(c *public.Conf) *public.Conf {
	yamlFile, err := ioutil.ReadFile(GetConfigDir())
	if err != nil {
		fmt.Println()
		public.Error.Printf("conf.yml read err!")
		//os.Exit(1)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println()
		public.Error.Printf("conf.yml read err!")
		os.Exit(1)
	}
	return c
}