package public

// 输入指令
type InputType struct {
	ModePtr     string
	IpsPtr      string
	IfilePtr    string
	OfilePtr    string
	ThreadPtr   int
	TimePtr     int
	PortPtr     string
	WebPtr      bool
	DirPtr      string
	DModePtr    string
	HeaderPtr   string
	BodyPtr     string
	RCodePtr    int
	RHeaderPtr  string
	RBodyPtr    string
	KeyPtr      string
	NumPtr      int
	FieldsPtr   string
	UrlPtr      string
	HelpPtr     bool
	PingPtr     bool
	NoPingPtr   bool
	NoOutputPtr bool
	UserFilePtr string
	PassFilePtr string
	UserPtr     string
	PassPtr     string
}

// 存放识别到的扫描、爆破、保存的数据
type AutoParType struct {
	Thread   int
	Timeout  int
	Follow   bool
	Port     []string
	Ip       []string
	IpAlive  []string
	Url      []string
	User     []string
	Pass     []string
	WebCode  []int
	DirInfo  DirScanConf
	FofaInfo FofaScanConf
	FofaWord []string
	Brute    BruteHostType
}

/* 端口指纹识别后存放的资产，爆破时从这里提取 */
type BruteHostType struct {
	Ssh      BruteInfoType
	Redis    BruteInfoType
	Ftp      BruteInfoType
	Mysql    BruteInfoType
	Mssql    BruteInfoType
	Postgres BruteInfoType
	Smb      BruteInfoType
	NetBios  BruteInfoType
	Rdp      BruteInfoType
	Telnet   BruteInfoType
	Tomcat   BruteInfoType
}

type BruteInfoType struct {
	Host []string
	Info BruteParConf
}

// 指纹扫描配置
type RuleDataType struct {
	Name string
	Type string
	Url  string
	Rule string
}

// -- 配置文件 --
type Conf struct {
	Thread   int          `yaml:"thread"`
	Timeout  int          `yaml:"timeout"`
	Follow   bool         `yaml:"follow_redirects"`
	PortScan PortScanConf `yaml:"portScan"`
	DirScan  DirScanConf  `yaml:"dirScan"`
	FofaScan FofaScanConf `yaml:"fofaScan"`
	Brute    BruteConf    `yaml:"brute"`
}

// 端口扫描配置
type PortScanConf struct {
	Ports      []string `yaml:"ports,flow"`
	StatusCode []int    `yaml:"code,flow"`
}

// dir扫描配置
type DirScanConf struct {
	Dir     string   `yaml:"dir"`
	Mode    string   `yaml:"mode"`
	Code    int      `yaml:"code"`
	Header  []string `yaml:"header"`
	Body    string   `yaml:"body"`
	RHeader string   `yaml:"rheader"`
	RBody   string   `yaml:"rbody"`
}

// fofa配置
type FofaScanConf struct {
	Email  string `yaml:"email"`
	Key    string `yaml:"key"`
	Num    int    `yaml:"num"`
	Fields string `yaml:"fields"`
}

// 爆破模块配置
type BruteConf struct {
	Thread   int          `yaml:"thread"`
	Timeout  int          `yaml:"timeout"`
	Redis    BruteParConf `yaml:"redis"`
	Ssh      BruteParConf `yaml:"ssh"`
	Ftp      BruteParConf `yaml:"ftp"`
	Mysql    BruteParConf `yaml:"mysql"`
	Mssql    BruteParConf `yaml:"mssql"`
	Postgres BruteParConf `yaml:"postgres"`
	Smb      BruteParConf `yaml:"smb"`
	NetBios  BruteParConf `yaml:"netbios"`
	Rdp      BruteParConf `yaml:"rdp"`
	Telnet   BruteParConf `yaml:"telnet"`
	Tomcat   BruteParConf `yaml:"tomcat"`
}

// 爆破参数配置
type BruteParConf struct {
	User []string `yaml:"user,flow"`
	Pass []string `yaml:"pass,flow"`
	Port []string `yaml:"port,flow"`
}
