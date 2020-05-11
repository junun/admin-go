package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg *ini.File

	RunMode string
	//
	//HTTPPort int
	//ReadTimeout time.Duration
	//WriteTimeout time.Duration

	//PageSize int
	//JwtSecret string
	AppSetting = &App{}
	ServerSetting = &Server{}
)

type App struct {
	JwtSecret string
	PageSize int
	RuntimeRootPath string

	ImagePrefixUrl string
	ImageSavePath string
	ImageMaxSize int
	ImageAllowExts []string

	LogSavePath string
	LogSaveName string
	LogFileExt string
	TimeFormat string

	GitLocalPath string
	IdRsaPath 	string
	SyncPath 	string
	DeployPath 	string
	GitSshKey 	string
}

type Server struct {
	RunMode string
	HttpPort int
	ReadTimeout time.Duration
	WriteTimeout time.Duration
}

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	RunMode = Cfg.Section("").Key("RUN_MODE").MustString("debug")

	if RunMode == "debug" {
		Cfg, err = ini.Load("conf/debug.ini")
		if err != nil {
			log.Fatalf("Fail to parse 'conf/debug.ini': %v", err)
		}
	} else {
		Cfg, err = ini.Load("conf/release.ini")
		if err != nil {
			log.Fatalf("Fail to parse 'conf/release.ini': %v", err)
		}
	}

	LoadServer()
	LoadApp()
}


func LoadServer() {
	//sec, err := Cfg.GetSection("server")
	//if err != nil {
	//	log.Fatalf("Fail to get section 'server': %v", err)
	//}

	//HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)
	//ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second
	//WriteTimeout =  time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second

	err := Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.ReadTimeout * time.Second
}

func LoadApp() {
	//sec, err := Cfg.GetSection("app")
	//if err != nil {
	//	log.Fatalf("Fail to get section 'app': %v", err)
	//}

	//JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	//PageSize = sec.Key("PAGE_SIZE").MustInt(10)

	err := Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024

}
