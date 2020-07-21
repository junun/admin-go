package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"runtime"
	"strings"
)

func foo() {
	fmt.Println(runtime.GOOS)
}

func main() {
	str := `{"title":"setup1","data|wc -l":"echo \"need reload\""}|{"title":"setup2","data":"ping -c 1 127.0.0.1\necho \"reloading\"\nw"}`

	fmt.Println(strings.Split(str, "|"))
	//models.DingtalkSentChannel(4, "测试企业微信应用", []string{"15818699723"}, false)

	//fmt.Println(models.ReturnWebChatAccessToken("ww418a1a1d690697ab", "zlQVK5rJM0LE3u95Zr4rhZjdDEGRVRJGNrzHDs-9H_4"))

	//key, e := rsa.GenerateKey(rand.Reader, 2048)
	//if e != nil {
	//	fmt.Println("Private key cannot be created.", e.Error())
	//}
	//
	//publickey := &key.PublicKey
	//
	//fmt.Println(util.DumpPrivateKeyBuffer(key))
	//fmt.Println(util.DumpPublicKeyBuffer(publickey))

	//str := "1_2_123"
	//fmt.Println(strings.Join(strings.Split(str, "_")[0:2],"_"))

	//fmt.Println(util.HumanNowTime())
	//str := `{"version":"1.0.0","name":"\"sweet hony\"","port":"9090"}`
	//data := make(map[string]string)
	//json.Unmarshal([]byte(str), &data)
	//fmt.Println(data)
	//fmt.Println(util.ReturnGitTagByCommand(2, "ssh://git@git.zhien88.com:2222/zeyw/nginx"))

	//var about models.About
	//
	//about.Golangversion = runtime.Version()
	//about.SystemInfo = runtime.GOOS
	//fmt.Println(about)
	//fmt.Println(gin.Version)
	//admin.CheckDomainAndCret()

	//_,  res := curlwhois.Whois("qq.com")
	//fmt.Println(res)

	//easy := curl.EasyInit()
	//defer easy.Cleanup()
	//
	//easy.Setopt(curl.OPT_URL, "https://lookup.icann.org/api/whois?q=wikipedia.org")
	//easy.Setopt(curl.OPT_SSL_VERIFYPEER, false)

	//easy.Setopt(curl.OPT_HTTPHEADER,[]string{"Connection: keep-alive", "Cache-Control: max-age=0", "Upgrade-Insecure-Requests: 1", "User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.106 Safari/537.36", "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"})

	//easy.Setopt(curl.OPT_VERBOSE, true)
	//
	//res := easy.Perform()
	//fmt.Println(res)

	//domain := "baidu.com"
	//result, e := whois.Whois(domain)
	//fmt.Println(e)
	//
	//fmt.Println(result)
	//
	//res, _ := whoisparser.Parse(result)
	//fmt.Println(res.Administrative)
	//fmt.Println(res.Domain)
	//fmt.Println(res.Domain.ExpirationDate)
	//fmt.Println(res.Domain.Status)

	//cert, e := util.ParseRemoteCertificate("images.baidu.com:443",10)
	//if e!=nil {
	//	fmt.Println(e)
	//}
	//fmt.Println(cert.Jsonify())
	//
	//str1 := "4.75387ms"
	//
	//fmt.Println(strconv.Atoi(str1))

//	startTime 	:= time.Now()
//	fmt.Println(admin.CheckIfLocalIp("127.0.0.1"))
//
//	cmd := `
//echo "hello world"
//who
//ifconfig
//`
//	util.ExecRuntimeCmd(cmd)
//
//	fmt.Println(time.Since(startTime).String())
	//str := "1"
	//a, _ := strconv.Atoi(str)
	//fmt.Println(a)
	//
	//var err = fmt.Errorf("%s", "the error test for fmt.Errorf")
	//fmt.Println(err)
	//
	//
	//i := 0
	//defer fmt.Println("a:", i)
	////闭包调用，将外部i传到闭包中进行计算，不会改变i的值，如上边的例3
	//defer func(i int) {
	//	fmt.Println("b:", i)
	//}(i)
	////闭包调用，捕获同作用域下的i进行计算
	//defer func() {
	//	fmt.Println("c:", i)
	//}()
	//i++
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
