package curlwhois

import (
	"api/pkg/logging"
	"encoding/json"
	"fmt"
	"github.com/andelf/go-curl"
	"math/rand"
	"os"
	"regexp"
	"strings"
)

var (
	USER_AGENTS = []string{"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; AcooBrowser; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
		"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0; Acoo Browser; SLCC1; .NET CLR 2.0.50727; Media Center PC 5.0; .NET CLR 3.0.04506)",
		"Mozilla/4.0 (compatible; MSIE 7.0; AOL 9.5; AOLBuild 4337.35; Windows NT 5.1; .NET CLR 1.1.4322; .NET CLR 2.0.50727)",
		"Mozilla/5.0 (Windows; U; MSIE 9.0; Windows NT 9.0; en-US)",
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 2.0.50727; Media Center PC 6.0)",
		"Mozilla/5.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0; WOW64; Trident/4.0; SLCC2; .NET CLR 2.0.50727; .NET CLR 3.5.30729; .NET CLR 3.0.30729; .NET CLR 1.0.3705; .NET CLR 1.1.4322)",
		"Mozilla/4.0 (compatible; MSIE 7.0b; Windows NT 5.2; .NET CLR 1.1.4322; .NET CLR 2.0.50727; InfoPath.2; .NET CLR 3.0.04506.30)",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN) AppleWebKit/523.15 (KHTML, like Gecko, Safari/419.3) Arora/0.3 (Change: 287 c9dfb30)",
		"Mozilla/5.0 (X11; U; Linux; en-US) AppleWebKit/527+ (KHTML, like Gecko, Safari/419.3) Arora/0.6",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.8.1.2pre) Gecko/20070215 K-Ninja/2.1.1",
		"Mozilla/5.0 (Windows; U; Windows NT 5.1; zh-CN; rv:1.9) Gecko/20080705 Firefox/3.0 Kapiko/3.0",
		"Mozilla/5.0 (X11; Linux i686; U;) Gecko/20070322 Kazehakase/0.4.5",
		"Mozilla/5.0 (X11; U; Linux i686; en-US; rv:1.9.0.8) Gecko Fedora/1.9.0.8-1.fc10 Kazehakase/0.5.6",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_3) AppleWebKit/535.20 (KHTML, like Gecko) Chrome/19.0.1036.7 Safari/535.20",
		"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; fr) Presto/2.9.168 Version/11.52"}
)

//const bodyfilename   = "/Users/angus/Documents/md_document/myproject/github/go_spug/api/pkg/test/body.out"

type ResultInfo struct {
	CreationDate			string
	ExpirationDate  		string
}

// check domain info by curl
// api url  https://lookup.icann.org/api/whois?q=xx.xx
func Whois(domian string) (error, ResultInfo) {
	c := curl.EasyInit()
	defer c.Cleanup()

	c.Setopt(curl.OPT_URL, "https://lookup.icann.org/api/whois?q=" + domian)
	c.Setopt(curl.OPT_SSL_VERIFYPEER, false)
	c.Setopt(curl.OPT_HTTPHEADER, SetHttpHeader())

	//c.Setopt(curl.OPT_VERBOSE, true)

	var content []byte
	readContent := func (buf []byte, userdata interface{}) bool {
		//println("DEBUG: size=>", len(buf))
		//println("DEBUG: content=>", string(buf))
		content = buf
		return true
	}

	c.Setopt(curl.OPT_WRITEFUNCTION, readContent)


	// write file
	//c.Setopt(curl.OPT_WRITEFUNCTION, writeDataToFile)
	//fp, _ := os.OpenFile(bodyfilename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	//defer fp.Close()
	//c.Setopt(curl.OPT_WRITEDATA, fp)

	if e := c.Perform(); e != nil {
		logging.Error("ERROR: ", e.Error())
		return e, ResultInfo{}
	}

	return nil , parseContent(content)
}

func writeDataToFile(ptr []byte, userdata interface{}) bool {
	fp := userdata.(*os.File)
	if _, err := fp.Write(ptr); err == nil {
		return true
	}
	return false
}


func SetHttpHeader() []string{
	header := []string{"Connection: keep-alive", "Cache-Control: max-age=0", "Upgrade-Insecure-Requests: 1", "Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"}

	userAgent := "User-Agent:"+USER_AGENTS[rand.Intn(len(USER_AGENTS))]
	header = append(header, userAgent)

	return header
}

func parseContent(byt []byte) ResultInfo {
	var dat map[string]interface{}
	var result ResultInfo

	if err := json.Unmarshal(byt, &dat); err != nil {
		logging.Error(err)
	}
	res := dat["records"]

	// interface to string
	str := fmt.Sprintf("%v", res)

	// 取过期时间
	rExpirationDate := regexp.MustCompile(`Registrar Registration Expiration Date:.*`)
	expirationDate := rExpirationDate.FindStringSubmatch(str)

	// 取创建时间
	rCreationDate := regexp.MustCompile(`Creation Date:.*`)
	creationDate := rCreationDate.FindStringSubmatch(str)

	result.ExpirationDate = strings.Split(expirationDate[0],":")[1]
	result.CreationDate = strings.Split(creationDate[0],":")[1]


	return result
}