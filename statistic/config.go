package statistic

type Config struct {
	OutPath string `mapstructure:"outpath" json:"outpath" yaml:"outpath"`
	Mysql   Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Gap     int    `mapstructure:"gap" json:"gap" yaml:"gap"`
}

type Mysql struct {
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	DBname   string `mapstructure:"dbname" json:"dbname" yaml:"dbname"`
	Charset  string `mapstructure:"charset" json:"charset" yaml:"charset"`
	// Config   string
}
