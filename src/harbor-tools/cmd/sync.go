package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"harbor-tools/harbor-tools/controllers"
	"harbor-tools/harbor-tools/db"
	"harbor-tools/harbor-tools/models"
	"log"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "from source harbor sync tags to target harbor ",
	Run: func(cmd *cobra.Command, args []string) {
		MyDB, HarborDB, err :=db.DB()
		defer HarborDB.Close()
		defer MyDB.Close()
		if  err != nil {
			log.Panicf("db init fail error:", err)
		}
		fmt.Println("Start sync images")
		MyDB.AutoMigrate(&models.JobStatus{})
		controllers.StartSync()
		//controllers.ProductTagsFromFile()
		log.Println("sync End")
		//controllers.DBdemo()
		//controllers.ProductTagsFromHarborDB()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().String("pull-worker", "5","pull workers")
	syncCmd.Flags().String("push-worker", "5","push workers")
	syncCmd.Flags().Bool("manual", true, "manual execute sync")

	syncCmd.Flags().String("file","", "from file")
	viper.BindPFlag("goroutines.pull-worker", syncCmd.Flags().Lookup("pull-worker"))
	viper.BindPFlag("goroutines.pull-worker", syncCmd.Flags().Lookup("push-worker"))
	viper.BindPFlag("sync.manual", syncCmd.Flags().Lookup("manual"))
	viper.BindPFlag("file", syncCmd.Flags().Lookup("file"))
}
