package controllers

import (
	"fmt"
	DB "harbor-tools/harbor-tools/db"
	"harbor-tools/harbor-tools/models"
	"time"

	//"log"

	log "harbor-tools/harbor-tools/utils/logger"
	"strconv"
	"strings"
	"sync"

	"harbor-tools/harbor-tools/common"
	"harbor-tools/harbor-tools/config"
)

var cfg = config.NewConfig()

//var log = logger.NewLogger()
func DealTags(wg *sync.WaitGroup, goroutineNum int, inCh <-chan string) {
	defer wg.Done()
	db, _, _ := DB.DB()

	TargetUrl := cfg.TargetHarbor.Addr + "/api/repositories/"
	for t := range inCh {
		log.Println(t)
		tag := models.ImageTag{Tag: t}
		db.NewRecord(tag)
		db.Create(&tag)
		log.Printf("%s insert to db", t)
		s := strings.Split(t, ":")
		repoName := s[0]
		tagName := s[1]
		tagUrl := TargetUrl + repoName + "/tags/" + tagName
		log.Printf("Goroutine num: %d", goroutineNum)
		statusCode, err := common.HttpClient("GET", tagUrl, nil, nil)
		log.Printf("targetUrl: %s, status: %d ", tagUrl, statusCode)
		if err != nil {
			//log.Printf("target harbor:error %s", err)
			if statusCode == 404 {
				log.Printf("%s get err:", err)
				diffTag := models.DiffTag{Tag: t}
				db.NewRecord(diffTag)
				db.Create(&diffTag)
			}

		}
		fmt.Println(tagUrl)

	}

}
func FetchTags() {
	var projects []Projects
	var wg sync.WaitGroup
	//db, _ := DB.DB()
	harborUrl := cfg.Harbor.Addr + "/api/projects"
	//TargetUrl := cfg.TargetHarbor.Addr + "/api/repositories/"
	_, err := common.HttpClient("GET", harborUrl, nil, &projects)

	if err != nil {
		log.Panic("Get projects err:", err)
	}

	repoCh := GetRepos(projects...)
	repoChs := []<-chan string{}
	for i := 0; i < cfg.Goroutines.TagWorkers; i++ {
		c := GetTags(repoCh)
		repoChs = append(repoChs, c)
	}
	//count := 0

	mergeCh := MergeCh(repoChs...)
	wg.Add(cfg.Goroutines.DBWorkers)
	for i := 0; i < cfg.Goroutines.DBWorkers; i++ {
		go DealTags(&wg, i, mergeCh)
	}
	wg.Wait()

	log.Println("Fetch tags End")
}

func GetRepos(proj ...Projects) <-chan string {

	out := make(chan string, 100)
	var count int
	go func() {
		defer close(out)
		for _, p := range proj {
			repos := []Repositories{}
			harburUrl := cfg.Harbor.Addr + "/api/repositories?project_id=" + strconv.Itoa(p.ProjectID)
			log.Println("projectUrl:", harburUrl)
			_, err := common.HttpClient("GET", harburUrl, nil, &repos)
			if err != nil {
				log.Println("productRepos, err:", err)
			} else {
				count = len(repos)
				log.Println("count", count)
				for _, repo := range repos {
					log.Println("repo", repo.Name, repo.ProjectID)
					out <- repo.Name
				}
			}
		}

	}()
	return out
}

func GetTags(repoCh <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for repo := range repoCh {
			tags := []Tags{}
			harborUrl := cfg.Harbor.Addr + "/api/repositories/" + repo + "/tags"
			log.Println(harborUrl)
			log.Println("start getTags,", time.Now())
			startTime := time.Now()
			_, err := common.HttpClient("GET", harborUrl, nil, &tags)
			if err != nil {
				log.Printf("GetTags,err:%s, tagsUrl:%s", err, harborUrl)
			} else {
				for _, tag := range tags {
					out <- repo + ":" + tag.Name
				}
			}
			endTime := time.Now()
			log.Println("SpendTime:", endTime.Sub(startTime))

		}
	}()
	return out
}

func MergeCh(cs ...<-chan string) <-chan string {
	out := make(chan string, 350)
	var wg sync.WaitGroup

	collect := func(in <-chan string) {
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
