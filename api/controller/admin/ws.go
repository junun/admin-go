package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/util"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		CheckOrigin:      func(r *http.Request) bool { return true },
		HandshakeTimeout: time.Duration(time.Second * 5),
	}
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)


type Conn struct {
	Conn *websocket.Conn

	AfterReadFunc   func(messageType int, r io.Reader)
	BeforeCloseFunc func()

	once   sync.Once
	id     string
	stopCh chan struct{}
}

func wshandleError(ws *websocket.Conn, err error) bool {
	if err != nil {
		logrus.WithError(err).Error("handler ws ERROR:")
		dt := time.Now().Add(time.Second)
		if err := ws.WriteControl(websocket.CloseMessage, []byte(err.Error()), dt); err != nil {
			logrus.WithError(err).Error("websocket writes control message failed:")
		}
		return true
	}
	return false
}

func HandleWebsocket(c *gin.Context) {
	//升级get请求为webSocket协议
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	defer conn.Close()
	if err != nil {
		log.Println("cant upgrade connection:", err)
		return
	}

	for {
		msgType, msgData, err := conn.ReadMessage()
		if err != nil {
			log.Println("cant read message:", err)

			switch err.(type) {
			case *websocket.CloseError:
				return
			default:
				return
			}
		}

		// Skip binary messages
		if msgType != websocket.TextMessage {
			continue
		}

		if string(msgData) == "ping" {
			msgData = []byte("pong")
		}
		//写入ws数据
		err = conn.WriteMessage(msgType, msgData)
		if err != nil {
			break
		}
	}
}

// handle webSocket connection.
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func SshConsumer(c *gin.Context)  {
	middleware.WsTokenAuthMiddleware("host-console")

	var host models.Host
	models.DB.Find(&host, c.Param("id"))
	if host.ID <= 0 {
		util.JsonRespond(500, "Unknown Host！", "", c)
		return
	}

	clientConfig, _ := util.ReturnClientConfig(host.Username, "")
	hostIp := host.Addres + ":" + strconv.Itoa(host.Port)

	Scli, err := util.GetSshClient(hostIp, clientConfig)
	if err != nil {
		log.Fatalf("connect host err: %v", err)
	}
	defer Scli.Close()

	cols, _ := strconv.Atoi(c.DefaultQuery("cols", "512"))
	rows, _ := strconv.Atoi(c.DefaultQuery("rows", "512"))
	ssConn, err := util.NewSshConn(cols, rows, Scli)
	//if wshandleError(c, err) {
	//	//	return
	//	//}
	defer ssConn.Close()

	//升级get请求为webSocket协议
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("cant upgrade connection:", err)
		return
	}
	defer wsConn.Close()

	quitChan := make(chan bool, 3)

	var logBuff = new(bytes.Buffer)
	// most messages are ssh output, not webSocket input
	go ssConn.ReceiveWsMsg(wsConn, logBuff, quitChan)
	go ssConn.SendComboOutput(wsConn, quitChan)
	go ssConn.SessionWait(quitChan)

	<-quitChan
	//保存日志

	logrus.Info("websocket finished")
}

func SyncWsApp(c *gin.Context)  {
	middleware.WsTokenAuthMiddleware("config-app-async")

	//升级get请求为webSocket协议
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("cant upgrade connection:") )
		return
	}
	defer wsConn.Close()

	// 查询项目是否是否已经初始化过
	id 	:= c.Param("id")
	var hostapp []models.HostApp
	var app models.App
	var env models.ConfigEnv

	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", id).
		Find(&hostapp)

	if  len(hostapp) == 0 {
		util.WsWriteMessage("该项目没有绑定主机，请先绑定到对应主机上！", wsConn)
		return
	}

	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", id).
		Where("status = 0").
		Find(&hostapp)

	if  len(hostapp) == 0 {
		util.WsWriteMessage("该项目已经初始过，请勿重复初始化！", wsConn)
		return
	}

	models.DB.Model(&models.App{}).
		Where("id = ?", id).
		Find(&app)

	models.DB.Model(&models.ConfigEnv{}).
		Where("id = ?", app.EnvId).
		Find(&env)

	for _, v := range hostapp {
		if v.Status == 0 {
			var host models.Host
			models.DB.Model(&models.Host{}).
				Where("id = ?",  v.Hid).
				Find(&host)

			// 执行初始化操作
			util.WsWriteMessage("开始初始化项目到主机" + host.Name, wsConn)
			desDir 	:= util.ReturnSyncRunDir(env.Name, app.ID)
			desDir 	= desDir + "/" + app.Name

			// 建立 sftp 传文件
			util.WsWriteMessage("建立sftp文件传输通道！", wsConn)

			clientConfig, _ := util.ReturnClientConfig(host.Username, "")
			hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
			Scli, err 		:= util.GetSshClient(hostIp, clientConfig)
			if err 	!= nil {
				util.WsWriteMessage(err.Error(), wsConn)
				return
			}

			SftpClient, err := util.GetSftpClient(Scli)
			if err 	!= nil {
				util.WsWriteMessage(err.Error(), wsConn)
				return
			}

			util.WsWriteMessage("开始执行自定义初始化任务", wsConn)

			e 	:= syncBywsId(app , desDir, app.Name,  SftpClient, Scli, wsConn)
			if e!= nil {
				util.WsWriteMessage(err.Error(), wsConn)
				return
			}


			// 初始化完成 修改状态
			e 	= models.DB.Model(&models.HostApp{}).
				Where("id = ?", v.ID).
				Update("status", 1).
				Error
			if e!= nil {
				util.WsWriteMessage(err.Error(), wsConn)
				return
			}
		}
	}

	util.WsWriteMessage("项目初始化完毕!", wsConn)
}

func syncBywsId(app models.App, path string, name string, SftpClient *sftp.Client, Scli *ssh.Client, ws *websocket.Conn) error {
	switch app.Tid {
	// backend jar
	case 1:
		// copy 项目样本文件夹到 Sync RunTime Dir
		dir, _ 	:= os.Getwd()
		srcDir 	:= dir + "/files/AppInit/BackendJar/origin"
		e   	:= util.CopyDirectory(srcDir, path)
		if e != nil {
			util.WsWriteMessage(e.Error(), ws)
			return e
		}

		// 检查项目变量
		var value []models.AppSyncValue
		models.DB.Model(&models.AppSyncValue{}).
			Where("aid = ?", app.ID).
			Find(&value)

		if len(value) <= 0 {
			var e = fmt.Errorf("%s", "该项目没有设置环境变量，请先设置!")
			util.WsWriteMessage("该项目没有设置环境变量，请先设置!", ws)
			return e
		}

		// 生成本地变量文件 var
		f, e 	:= os.OpenFile(path + "/bin/var", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer f.Close()
		if e 	!= nil {
			log.Fatal(e)
			util.WsWriteMessage(e.Error(), ws)
			return e
		}


		for index, v := range value {
			str := v.Name + "='" + v.Value + "'\n"
			if index == 0 {
				e := ioutil.WriteFile(path + "/bin/var", []byte(str), 0644)
				if e != nil {
					util.WsWriteMessage(e.Error(), ws)
					return e
				}
			} else {
				_, e = f.Write([]byte(str))
				if e != nil {
					util.WsWriteMessage(e.Error(), ws)
					return e
				}
			}
		}

		// 上传文件
		util.WsWriteMessage("上传BackendJar项目通用系统依赖文件 jarFuncs", ws)
		srcFile := dir + "/files/AppInit/BackendJar/jarFuncs"
		e 		= util.PutFile(SftpClient, srcFile, "/etc/init.d/")
		if e != nil {
			util.WsWriteMessage(e.Error(), ws)
			return e
		}

		e 		= SftpClient.Chmod("/etc/init.d/jarFuncs", 0755)
		if e 	!= nil {
			util.WsWriteMessage(e.Error(), ws)
			return e
		}

		// 上传目录
		util.WsWriteMessage("上传项目基本结构", ws)
		e   	= util.PutDirectory(SftpClient, path, "/data/webapps/" + name)
		if e 	!= nil {
			util.WsWriteMessage(e.Error(), ws)
			return e
		}

		util.WsWriteMessage("改变目标主机项目目录权限", ws)
		pathDir := "/data/webapps/" + name
		cmdstr 	:= "chown -R tomcat:tomcat " + pathDir

		_, e 	= util.ExecuteCmd(cmdstr,  Scli)
		if e 	!= nil {
			util.WsWriteMessage(e.Error(), ws)
			return e
		}

		util.WsWriteMessage("脚本添加执行权限", ws)
		pathDir = "/data/webapps/" + name + "/bin/*"
		cmdstr 	= "chmod +x " + pathDir

		_, e 	= util.ExecuteCmd(cmdstr, Scli)
		if e 	!= nil {
			util.WsWriteMessage(e.Error(), ws)
			return e
		}

	case 2:
		fmt.Println(2)
		return nil
	}
	return nil
}





