package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func main() {
	i := 0
	defer fmt.Println("a:", i)
	//闭包调用，将外部i传到闭包中进行计算，不会改变i的值，如上边的例3
	defer func(i int) {
		fmt.Println("b:", i)
	}(i)
	//闭包调用，捕获同作用域下的i进行计算
	defer func() {
		fmt.Println("c:", i)
	}()
	i++
	//res := ""
	//page, _ := com.StrTo(res).Int()
	//fmt.Println(page)
	//
	//dir, _ 		:= os.Getwd()
	//fmt.Println(dir)
	//
	//
	////layout 	:= "2006-01-02T15:04:05.000Z"
	//str 	:= "2014-11-12T11:45:26.371Z"
	//t, err 	:= time.Parse(time.RFC3339, str)
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//fmt.Println(t)

	//fmt.Println(1<<0)


	//##### go run -race main.go
	//a := 1
	//go func(){
	//	a = 2
	//}()
	//a = 3
	//fmt.Println("a is ", a)
	//
	//time.Sleep(2 * time.Second)


	//##### 变量和指针
	//var a = 3
	//double(&a)
	//fmt.Println(a) // 6

	// cron
	//maps 	:= make(map[string]interface{})
	//c := cron.New()
	//c.AddFunc("*/1 * * * *", func() { fmt.Println("Every 10s ") })
	//c.AddFunc("*/1 * * * *", func() { fmt.Println("Every 11s ") })
	//c.AddFunc("*/1 * * * *", func() { fmt.Println("Every 12s ") })
	//c.AddFunc("*/2 * * * *", func() { fmt.Println("Every 2min ") })
	//c.Start()
	//
	//maps["c"] = c
	//
	//fmt.Println(c.Entries())
	//
	//for k,v := range c.Entries() {
	//	fmt.Println(k)
	//	fmt.Println(v.Schedule)
	//	fmt.Println(v.Job)
	//}
	//
	//c1 := cron.New()
	//id,_	:= c1.AddFunc("*/1 * * * *", func() { fmt.Println("Every 1min ") })
	//fmt.Println(id)
	//
	//fmt.Println(c1)
	//c1.Remove(id)
	//maps["c1"] = c1
	//
	//c1.Stop()
	//
	//fmt.Println(c)
	//fmt.Println(c1.Entries())
	////c.Stop()
	//
	////models.Rdb.Set("1111", util.JSONMarshalToString(maps), 0)
	////select {
	////}
	//
	//defer c.Stop()
	//defer c1.Stop()

}

func RemoveJobByName(name *cron.Cron, id cron.EntryID)  {
	name.Remove(id)
}

func double(x *int)  {
	fmt.Println(x)
	*x++
	fmt.Println(&x)
	x = nil
	fmt.Println(x)
}
