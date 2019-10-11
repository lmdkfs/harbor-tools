package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	global *Config
)

type Config struct {
	RunMode     string
	Store       string
	Swagger     string
	Log         Log
	LogGormHook LogGormHook
	HTTP        HTTP
	Gorm        Gorm
	Postgres    Postgres
	Harbor      Harbor
	Ldap        Ldap
}

func GetGlobalConfig() *Config {
	if global == nil {
		return &Config{}
	}
	return global
}

func NewGlobalConfig(ConfigPath string) (*Config, error) {
	vp := viper.New()
	vp.SetEnvPrefix("harbortools")                      // 环境变量前缀, 环境变量必须大写
	vp.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 为了兼容读取yaml文件
	vp.AutomaticEnv()

	// load config from yaml
	fmt.Println("configpath", ConfigPath)
	vp.AddConfigPath(ConfigPath)
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	if err := vp.ReadInConfig(); err != nil {
		return &Config{}, err
	}
	fmt.Println(">>>>>>>>",vp.ConfigFileUsed())

	// LoadServer
	//fmt.Println(vp.GetString("server.run_mode"))
	//fmt.Println(vp.GetInt("server.port"))
	global = new(Config)
	global.HTTP.Port = vp.GetInt("server.port")
	global.RunMode = vp.GetString("server.run_mode")
	global.Log.LogPath = vp.GetString("server.logpath")
	global.Log.LogName = vp.GetString("server.logname")

	// LoadHarbor
	global.Harbor.Addr = vp.GetString("harbor.addr")
	global.Harbor.User = vp.GetString("harbor.user")
	global.Harbor.Password = vp.GetString("harbor.password")

	// LoadLdap

	global.Ldap.Addr = vp.GetString("ldap.addr")
	global.Ldap.Dn = vp.GetString("ldap.dn")
	global.Ldap.BaseDn = vp.GetString("ldap.basedn")
	global.Ldap.Password = vp.GetString("ldap.password")
	global.Ldap.Uid = vp.GetString("ldap.uid")

	return global, nil

}

type Harbor struct {
	Addr     string
	User     string
	Password string
}

type Ldap struct {
	Addr     string
	Dn       string
	BaseDn   string
	Password string
	Uid      string
}

// Log 日志配置参数
type Log struct {
	LogPath       string
	LogName       string
	Level         int
	Format        string
	Output        string
	OutputFile    string
	EnableHook    bool
	Hook          string
	HookMaxThread int
	HookMaxBuffer int
}

// LogGormHook 日志gorm钩子配置
type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

// JWTAuth 用户认证
type JWTAuth struct {
	SigningMethod string
	SigningKey    string
	Expired       int
	Store         string
}

// HTTP http 配置参数
type HTTP struct {
	Host            string
	Port            int
	ShutdownTimeout int
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
}

// monitor 监控配置参数
type Monitor struct {
	Enable    bool
	Addr      string
	ConfigDir string
}

// Gorm 配置参数
type Gorm struct {
	Debug        bool
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdelConns int
	TablePrefix  string
}

// Postgres
type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// DSN 数据库连接串
func (pg Postgres) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		pg.Host, pg.Port, pg.User, pg.DBName, pg.Password)
}

// Sqlite3 sqlite3配置参数
type Sqlite3 struct {
	Path string
}

// DSN 数据库连接串
func (a Sqlite3) DSN() string {
	return a.Path
}
