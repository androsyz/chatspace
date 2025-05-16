package config

type Config struct {
	Server   Server   `yaml:"server" json:"server"`
	Database Database `yaml:"database" json:"database"`
	Settings Settings `yaml:"settings" json:"settings"`
}

type Server struct {
	Address  string `yaml:"address" json:"address"`
	Debug    bool   `yaml:"debug" json:"debug"`
	Name     string `yaml:"name" json:"name"`
	NameAddr string `yaml:"nameaddr" json:"nameaddr"`
}

type Settings struct {
	JWTSecret            string `yaml:"jwtsecret" json:"jwtsecret"`
	TokenDuration        int    `yaml:"tokenduration" json:"tokenduration"`
	RefreshTokenDuration int    `yaml:"refreshtokenduration" json:"refreshtokenduration"`
}

type Database struct {
	Address         string `yaml:"address" json:"address"`
	Port            string `yaml:"port" json:"port"`
	Username        string `yaml:"username" json:"username"`
	Password        string `yaml:"password" json:"password"`
	DBName          string `yaml:"dbname" json:"dbname"`
	MaxOpenConn     int    `yaml:"maxopenconn" json:"maxopenconn"`
	MaxIdleConn     int    `yaml:"maxidleconn" json:"maxidleconn"`
	ConnMaxIdleTime int    `yaml:"connmaxidletime" json:"connmaxidletime"`
	ConnMaxLifeTime int    `yaml:"connmaxlifetime" json:"connmaxlifetime"`
}
