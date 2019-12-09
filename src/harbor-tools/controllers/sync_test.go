package controllers

import (
	"sync"
	"testing"
	"github.com/robfig/cron/v3"
)

func TestSync(t *testing.T) {
	t.Log("start TestSync")

	t.Log("end TestSync")
}

func TestAddBeforRunning(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	cron := cron.New()
}