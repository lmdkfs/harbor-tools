package config

import (
	"fmt"
	"sync"
	"time"
)

var (
	global *Config
	once   sync.Once
)

type Config struct {
	RunMode        string
	Store          string
	Swagger        string
	Log            Log
	LogGormHook    LogGormHook
	HTTP           HTTP
	Gorm           Gorm
	Postgres       Postgres
	HarborPostgres Postgres
	Harbor         Harbor
	TargetHarbor   TargetHarbor
	Ldap           Ldap
	Goroutines     Goroutines
	Sync           Sync
	File           string
}

func NewConfig() *Config {
	once.Do(func() {
		global = &Config{}
	})

	return global
}

type Sync struct {
	Manual bool
}
type Goroutines struct {
	TagWorkers  int
	DBWorkers   int
	PullWorkers int
	PushWorkers int
}
type Harbor struct {
	Addr     string
	User     string
	Password string
}

type TargetHarbor struct {
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
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
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
