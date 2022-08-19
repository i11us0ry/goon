package public

func InitInput() *InputType{
	return &InputType{
		"all",
		"",
		"",
		"",
		0,
		0,
		"",
		false,
		"",
		"get",
		"",
		"",
		200,
		"",
		"",
		"",
		0,
		"",
		"",
		false,
		false,
		"",
		"",
		"",
		"",
	}
}

func InitPar(conf *Conf) *AutoParType{
	return &AutoParType{
		conf.Thread,
		conf.Timeout,
		true,
		conf.PortScan.Ports,
		[]string{},
		[]string{},
		[]string{},
		[]string{},
		[]string{},
		conf.PortScan.StatusCode,
		false,
		false,
		GetCurrentDir()+"/result.txt",
		DirScanConf{
			conf.DirScan.Dir,
			conf.DirScan.Mode,
			conf.DirScan.Code,
			conf.DirScan.Header,
			conf.DirScan.Body,
			conf.DirScan.RHeader,
			conf.DirScan.RBody,
		},
		FofaScanConf{
			conf.FofaScan.Email,
			conf.FofaScan.Key,
			conf.FofaScan.Num,
			conf.FofaScan.Fields,
		},
		[]string{},
		BruteHostType{
			BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Ssh.User,
					conf.Brute.Ssh.Pass,
					conf.Brute.Ssh.Port,
				},
			},BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Redis.User,
					conf.Brute.Redis.Pass,
					conf.Brute.Redis.Port,
				},
			},BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Ftp.User,
					conf.Brute.Ftp.Pass,
					conf.Brute.Ftp.Port,
				},
			},BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Mysql.User,
					conf.Brute.Mysql.Pass,
					conf.Brute.Mysql.Port,
				},
			},BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Mssql.User,
					conf.Brute.Mssql.Pass,
					conf.Brute.Mssql.Port,
				},
			},BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Postgres.User,
					conf.Brute.Postgres.Pass,
					conf.Brute.Postgres.Port,
				},
			},BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Smb.User,
					conf.Brute.Smb.Pass,
					conf.Brute.Smb.Port,
				},
			},
			BruteInfoType{
				[]string{},
				BruteParConf{
					[]string{},
					[]string{},
					conf.Brute.NetBios.Port,
				},
			},BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Rdp.User,
					conf.Brute.Rdp.Pass,
					conf.Brute.Rdp.Port,
				},
			}, BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Telnet.User,
					conf.Brute.Telnet.Pass,
					conf.Brute.Telnet.Port,
				},
			},BruteInfoType{
				[]string{},
				BruteParConf{
					conf.Brute.Tomcat.User,
					conf.Brute.Tomcat.Pass,
					conf.Brute.Tomcat.Port,
				},
			},
		},
	}
}

var InputValue  *InputType
var ConfValue   *AutoParType

func Init(conf *Conf){
	InputValue = InitInput()
	ConfValue = InitPar(conf)
}
