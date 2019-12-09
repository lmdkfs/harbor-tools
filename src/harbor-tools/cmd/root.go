package cmd

import (
	"fmt"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path"
	"time"

	"log"
	"os"
	"strings"

	"harbor-tools/harbor-tools/utils/logger"
	"harbor-tools/harbor-tools/config"

	"github.com/mitchellh/go-homedir"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "harbor-tools",
	Short: "harbor 辅助工具",
	Long:  "harbor 辅助工具",
}

// Execute ...
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cfg.yaml)")

	rootCmd.PersistentFlags().String("harbor-addr", "", "harbor-addr")
	rootCmd.PersistentFlags().String("harbor-user", "", "harbor user")
	rootCmd.PersistentFlags().String("harbor-pwd", "", "harbor password")
	rootCmd.PersistentFlags().String("logname", "harbor-tools.log", "harbor-tools logname")
	rootCmd.PersistentFlags().String("logpath", "", "harbor-tools logname")

	rootCmd.PersistentFlags().String("target-addr", "", "target harbor-addr")
	rootCmd.PersistentFlags().String("target-user", "", "target harbor user")
	rootCmd.PersistentFlags().String("target-pwd", "", "target harbor password")

	rootCmd.PersistentFlags().String("tag-worker", "5", "get tags goroutines number")
	rootCmd.PersistentFlags().String("db-worker", "5", "inserter into db goroutines number")

	viper.BindPFlag("harbor.addr", rootCmd.PersistentFlags().Lookup("harbor-addr"))
	viper.BindPFlag("harbor.user", rootCmd.PersistentFlags().Lookup("harbor-user"))
	viper.BindPFlag("harbor.password", rootCmd.PersistentFlags().Lookup("harbor-pwd"))

	viper.BindPFlag("target.addr", rootCmd.PersistentFlags().Lookup("target-addr"))
	viper.BindPFlag("target.user", rootCmd.PersistentFlags().Lookup("target-user"))
	viper.BindPFlag("target.password", rootCmd.PersistentFlags().Lookup("target-pwd"))

	viper.BindPFlag("log.logname", rootCmd.PersistentFlags().Lookup("logname"))
	viper.BindPFlag("log.logpath", rootCmd.PersistentFlags().Lookup("logpath"))

	viper.BindPFlag("goroutines.tagworker", rootCmd.PersistentFlags().Lookup("tag-worker"))
	viper.BindPFlag("goroutines.dbworker", rootCmd.PersistentFlags().Lookup("db-worker"))
	viper.SetDefault("log.logpath", "/var/log/harbor-tools/")

}

func initConfig() {
	var cfg = config.NewConfig()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Println("Get currentDir  Fail", err)
	}
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		viper.AddConfigPath(currentDir)
		viper.AddConfigPath(home)
		viper.SetConfigName(".cfg")
	}
	viper.SetEnvPrefix("harbortools")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}

	cfg.HTTP.Port = viper.GetInt("server.port")
	cfg.HTTP.Host = viper.GetString("server.host")
	cfg.HTTP.ReadTimeout = viper.GetDuration("server.readtimeout")
	cfg.HTTP.WriteTimeout = viper.GetDuration("server.writetimeout")

	cfg.Harbor.Addr = viper.GetString("harbor.addr")
	cfg.Harbor.User = viper.GetString("harbor.user")
	cfg.Harbor.Password = viper.GetString("harbor.pwd")

	cfg.TargetHarbor.Addr = viper.GetString("target.addr")
	cfg.TargetHarbor.User = viper.GetString("target.user")
	cfg.TargetHarbor.Password = viper.GetString("target.pwd")

	cfg.Postgres.DBName = viper.GetString("db.name")
	cfg.Postgres.Host = viper.GetString("db.host")
	cfg.Postgres.User = viper.GetString("db.user")
	cfg.Postgres.Password = viper.GetString("db.password")
	cfg.Postgres.Port = viper.GetInt("db.port")

	cfg.HarborPostgres.DBName = viper.GetString("harbordb.name")
	cfg.HarborPostgres.Host = viper.GetString("harbordb.host")
	cfg.HarborPostgres.User = viper.GetString("harbordb.user")
	cfg.HarborPostgres.Password = viper.GetString("harbordb.password")
	cfg.HarborPostgres.Port = viper.GetInt("harbordb.port")

	cfg.Log.LogName = viper.GetString("log.logname")
	cfg.Log.LogPath = viper.GetString("log.logpath")

	cfg.Goroutines.DBWorkers = viper.GetInt("goroutines.dbworker")
	cfg.Goroutines.TagWorkers = viper.GetInt("goroutines.tagworker")

	cfg.Goroutines.PullWorkers = viper.GetInt("goroutines.pull-worker")
	cfg.Goroutines.PushWorkers = viper.GetInt("goroutines.pull-worker")

	cfg.Sync.Manual = viper.GetBool("sync.manual")
	cfg.File = viper.GetString("file")

	fmt.Println("logpath", cfg.Log.LogPath)

	//utils.NewLogger()
	exist, err := PathExists(cfg.Log.LogPath)
	logger.Println("日志路径:%s", cfg.Log.LogPath)



	if err != nil {
		logger.Errorf("Get dir error![%v]\n", err)
	}
	if !exist {
		if err := os.MkdirAll(cfg.Log.LogPath, 0755); err != nil {
			logger.Errorf("MkdirsFailed![%v]\n", err)
		} else {
			logger.Info("Mkdirs Success!\n")
		}

	}

	logFileName := path.Join(cfg.Log.LogPath, cfg.Log.LogName)
	logWriter, err := rotatelogs.New(
		logFileName+".%Y-%m-%d-%H-%M.log",
		rotatelogs.WithLinkName(logFileName),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		logger.Errorf("config local file system logger err. %v", err)
	}


	writeMap := lfshook.WriterMap{
		logrus.DebugLevel: logWriter,
		logrus.InfoLevel:  logWriter,
		logrus.WarnLevel:  logWriter,
		logrus.ErrorLevel: logWriter,
		logrus.FatalLevel: logWriter,
		logrus.PanicLevel: logWriter,
	}
	lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{})
	logger.AddHook(lfHook)
	//logger.SetReportCaller(true)




}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}