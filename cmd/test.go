package main

import (
	"github.com/sirupsen/logrus"
	"grates/pkg/utils"
	"time"
)

func main() {
	date := utils.Date{Time: time.Now()}
	logrus.Info(date)
}
