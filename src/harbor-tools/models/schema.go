package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

// ImageTag tags for image
type ImageTag struct {
	gorm.Model

	Tag string `gorm:"type:varchar(255);"`
	//CreatedAt time.Time

}

// DiffTag is diferent set betwen source harbor and target harbor
type DiffTag struct {
	gorm.Model
	Tag string `gorm:"type:varchar(255);"`
}

type JobStatus struct {
	gorm.Model
	JobID int `gorm:"AUTO_INCREMENT"`
	Start *time.Time
	End   *time.Time
}

type AccessLog struct {
	RepoName string `gorm:"column:repo_name"`
	RepoTag string `gorm:"column:repo_tag"`
	UserName string `gorm:"column:username"`
	Operation string `gorm:"column:operation"`
	OpTime time.Time `gorm:"column:op_time"`
}

func (AccessLog) TableName() string {
	return "access_log"
}

func (JobStatus) TableName() string {
	return "job_status"
}