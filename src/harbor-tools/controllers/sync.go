package controllers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"harbor-tools/harbor-tools/common"
	DB "harbor-tools/harbor-tools/db"
	"harbor-tools/harbor-tools/models"
	log "harbor-tools/harbor-tools/utils/logger"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)
 //var  MyDB, harborDB, _ = DB.DB()
 func init() {
 	log.SetReportCaller(true)
 }
func AttemptCreatProjects(projectName string) {
	projectUrl := cfg.TargetHarbor.Addr + "/api/projects"

	//project := Projects{
	//	Name:     projectName,
	//	Metadata: map[string]interface{}{"public": true},
	//}
	project := ProjectInfo{
		ProjectName: projectName,
		Metadata:    map[string]interface{}{"public": "true"},
	}
	if jsonStr, err := json.Marshal(project); err != nil {
		log.Printf("json 编码失败:%s", err)
	} else {
		body := bytes.NewBuffer([]byte(jsonStr))
		log.Printf("Attemp to create project:%s, url:%s", projectName, projectUrl)
		//statusCode, err := common.HttpClient("POST", projectUrl, body, nil)
		switch statusCode, err := common.HttpClient("POST", projectUrl, body, nil); statusCode {
		case 201:
			log.Printf("status: %d, create project:%s success", statusCode, projectName)
		case 409:
			log.Printf("status: %d, project:%s has already exist", statusCode, projectName)
		default:
			log.Printf("status: %d, create project:%s error:%s", statusCode, projectName, err)

		}

	}

}

func RetagAndPush(inCh <-chan Tags, wg *sync.WaitGroup, goroutineNums int) {
	defer wg.Done()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Printf("New client error:%s", err)

	} else {
		authConfig := types.AuthConfig{
			Username: "admin",
			Password: "Harbor12345",
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			log.Printf("encode err:%s", err)
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		imgPushFunc := func(tag string) error {
			ret, err := cli.ImagePush(context.Background(), tag, types.ImagePushOptions{RegistryAuth: authStr})
			if err != nil {
				return err
			}
			defer ret.Close()

			io.Copy(os.Stdout, ret)
			return err

		}
		for tag := range inCh {

			srcTag := strings.Replace(cfg.Harbor.Addr, "https://", "", 1) + tag.Name
			dstTag := strings.Replace(cfg.TargetHarbor.Addr, "https://", "", 1) + tag.Name
			log.Printf("srcTag:%s, dstTag:%s", srcTag, dstTag)
			projectName := strings.Split(tag.Name, "/")[1]
			log.Printf("goroutine Num: %d, Start push: %s", goroutineNums, dstTag)
			AttemptCreatProjects(projectName)
			if err := cli.ImageTag(context.Background(), srcTag, dstTag); err != nil {
				log.Printf("ReTag image from %s to %s fail, error:%s", srcTag, dstTag, err)
				continue
			}

			if err := imgPushFunc(dstTag); err != nil {
				log.Printf("push image: %s fail: %s", dstTag, err)
				continue
			}
			log.Printf("pushed image %s, success", tag.Name)


		}
	}
}

func Merge(cs ...<-chan Tags) <-chan Tags {
	out := make(chan Tags, 30)
	var wg sync.WaitGroup

	collect := func(in <-chan Tags) {
		defer wg.Done()
		for n := range in {
			out <- n
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go collect(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func Pull(inCh <-chan Tags, goroutineNums int) (<-chan Tags) {
	outCh := make(chan Tags, 50)
	go func() {
		defer close(outCh)

		cli, err := client.NewEnvClient()
		if err != nil {
			log.Printf("New client error:%s", err)

		} else {
			authConfig := types.AuthConfig{
				Username: "admin",
				Password: "Harbor12345",
			}
			encodedJSON, err := json.Marshal(authConfig)
			if err != nil {
				log.Printf("encode json error:%s", err)
			}
			authStr := base64.URLEncoding.EncodeToString(encodedJSON)
			imgPullFunc := func(tag string) error {
				ret, err := cli.ImagePull(context.Background(), tag, types.ImagePullOptions{RegistryAuth: authStr})
				if err != nil {
					return err
				}
				defer ret.Close()

				io.Copy(os.Stdout, ret)
				return err

			}
			for tag := range inCh {
				tagStr := strings.Replace(cfg.Harbor.Addr, "https://", "", 1) + tag.Name
				log.Printf("goroutine num: %d, Start pull image: %s", goroutineNums, tagStr)
				if err := imgPullFunc(tagStr); err != nil {
					log.Printf("image: %s pull error: %s", tag, err)
					continue
				}
				outCh <- tag
				//log.Println(tagStr)
			}
		}

	}()

	return outCh

}
func ProductTagsFromHarborDB() <-chan Tags {
	MyDB, harborDB, err := DB.DB()
	if err != nil {
		log.Panic("获取db失败:", err)
	}
	var lastRecord models.JobStatus
	MyDB.Last(&models.JobStatus{}).Select("*").Order("job_id desc").Limit(1).Scan(&lastRecord)

	if lastRecord.Start == nil {
		lasteTime, _ := time.Parse("2006-01-02 15:04:05", "2019-11-21 23:50:51")
		lastRecord.Start = &lasteTime
		lastRecord.End = &lasteTime
	}
	fmt.Println(lastRecord.Start.Format("2006-01-02 15:04:05"), lastRecord.End.Format("2006-01-02 15:04:05"), lastRecord.JobID)
	rows, err := harborDB.Model(&models.AccessLog{}).Select("*").Where("operation = ? AND op_time >= ?", "push", fmt.Sprint(lastRecord.Start.Format("2006-01-02 15:04:05"))).Rows()
	if err != nil {
		log.Panic("查询失败:", err)
	}
	//startJobTime := time.Now()
	tagsCh := make(chan Tags, 50)
	go func() {
		for rows.Next() {
			var tag Tags
			var accessLog models.AccessLog
			err := harborDB.ScanRows(rows, &accessLog)
			if err != nil {
				log.Println("scanRow 失败:", err)
				continue
			}

			tag.Name = "/" + accessLog.RepoName + ":" + accessLog.RepoTag
			tagsCh <- tag
			log.Println(tag.Name)
			//fmt.Println(tag.Name)


		}
	}()

	return tagsCh

}

func ProductTagsFromFile(wg *sync.WaitGroup) <-chan Tags {


	tagsCh := make(chan Tags, 50)

	go func() {
		defer wg.Done()
		defer close(tagsCh)

		_, err := os.Stat(cfg.File)
		if err != nil {
			if os.IsNotExist(err) {
				log.Panic("File does not exist", err)
			}
		}
		f, err := os.Open(cfg.File)
		if err != nil {
			log.Panic("Error:", err)
		}
		defer f.Close()
		br := bufio.NewReader(f)
		for {
			a, _, c := br.ReadLine()
			if c == io.EOF{
				break
			}

			var tag Tags
			//fmt.Println(strings.Split(string(a), ":")[1])
			tag.Name =  string(a)
			tagsCh <- tag
			log.Printf("image tag: %s", tag.Name)

		}

	}()


	return tagsCh
}

func ProductTagsFromToolsDB(wg *sync.WaitGroup) <-chan Tags {
	db, _, err := DB.DB()
	if err != nil {
		log.Panic("获取db失败:", err)
	}
	rows, err := db.Model(&models.DiffTag{}).Select("*").Rows()
	if err != nil {
		log.Panic("查询失败:", err)
	}
	tagsCh := make(chan Tags, 50)
	go func() {
		defer rows.Close()
		defer close(tagsCh)
		defer wg.Done()

		for rows.Next() {
			var tag Tags
			var diffTag models.DiffTag

			err := db.ScanRows(rows, &diffTag)
			if err != nil {
				log.Println("scanRow 失败:", err)
			}
			t, err := url.Parse(cfg.Harbor.Addr + "/" + diffTag.Tag)
			if err != nil {
				log.Printf("url:%s 解析错误,error:%s", t, err)
				continue
			}

			tag.Name = t.RequestURI()
			tagsCh <- tag
			log.Println(tag.Name)
		}

	}()

	return tagsCh

}

func StartSync() {
	var wg sync.WaitGroup
	var productWG sync.WaitGroup
	log.SetReportCaller(true)
	var tagsCh <-chan Tags
	productWG.Add(1)
	MyDB, _, _ := DB.DB()
	startJobTime := time.Now()
	if cfg.File != "" {
		fmt.Println("&&&&&&7")
		tagsCh = ProductTagsFromFile(&productWG)
	} else {
		if !cfg.Sync.Manual {
			tagsCh = ProductTagsFromToolsDB(&productWG)
		}	else {
			tagsCh = ProductTagsFromHarborDB()
		}
	}



	tagsChs := []<-chan Tags{}
	for i := 0; i < cfg.Goroutines.PullWorkers; i++ {
		pullCh := Pull(tagsCh, i)
		tagsChs = append(tagsChs, pullCh)
	}
	mergeCh := Merge(tagsChs...)
	for i := 0; i < cfg.Goroutines.PushWorkers; i++ {
		go RetagAndPush(mergeCh, &wg, i)
	}
	wg.Add(cfg.Goroutines.PushWorkers)
	productWG.Wait()
	wg.Wait()
	endJodbTIme := time.Now()

	jobStatus := models.JobStatus{Start: &startJobTime, End: &endJodbTIme}
	MyDB.NewRecord(jobStatus)
	MyDB.Create(&jobStatus)

}

func DBdemo() 	{
	MyDB,_, err := DB.DB()
    if err != nil {
		log.Panic("获取db失败:", err)
	}
    var lastRecord models.JobStatus
    MyDB.Last(&models.JobStatus{}).Select("*").Order("job_id desc").Limit(1).Scan(&lastRecord)
	if lastRecord.Start == nil {
		lasteTime, _ := time.Parse("2006-01-02 15:04:05", "2019-11-21 23:54:51")
		lastRecord.Start = &lasteTime
	}
    fmt.Printf("%+v, %+v\n", lastRecord.Start, lastRecord.End)
}
