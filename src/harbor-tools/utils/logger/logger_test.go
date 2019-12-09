package logger

import "testing"

func TestLog(t *testing.T) {
	log.SetReportCaller(true)

	log.Println(">>>>")
}

