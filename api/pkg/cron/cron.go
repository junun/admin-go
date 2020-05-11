package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func AddCronJob()  {
	c 	:= cron.New()
	fmt.Println(c)
}
