package controllers

import (
	"fmt"
	"harbor-tools/harbor-tools/common"
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	t.Log("Start")
	var projects []Projects
	_, err := common.HttpClient("GET", "https://harbor-5.finupgroup.com/api/projects", nil, &projects)
	if err != nil {
		log.Panic("Get projects err:", err)
	}

	repoCh := GetRepos(projects...)
	//fmt.Println("count:", count)
	fmt.Println(len(repoCh))
	repoChs := []<-chan string{}
	for i := 0; i < 3; i++ {
		c := GetTags(repoCh)
		repoChs = append(repoChs, c)
	}
	count := 0
	for n := range MergeCh(repoChs...) {
		fmt.Println(n)
		count++

	}
	t.Log("count:", count)
	t.Log("End")
}
