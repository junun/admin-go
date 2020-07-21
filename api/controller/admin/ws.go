package admin

import (
	"api/middleware"
	"api/models"
	"api/pkg/help"
	"api/pkg/logging"
	"api/pkg/util"
	"bytes"
	"encoding/json"
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
	"strings"
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

	GlobalVersion = ""
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

type CustomDeployParse struct {
	Title 	string 	`json:"title"`
	Data 	string	`json:"data"`
}

type Conn struct {
	Conn *websocket.Conn

	AfterReadFunc   func(messageType int, r io.Reader)
	BeforeCloseFunc func()

	once   sync.Once
	id     string
	stopCh chan struct{}
}

// handle webSocket connection.
// first,we establish a ssh connection to ssh server when a webSocket comes;
// then we deliver ssh data via ssh connection between browser and ssh server.
// That is, read webSocket data from browser (e.g. 'ls' command) and send data to ssh server via ssh connection;
// the other hand, read returned ssh data from ssh server and write back to browser via webSocket API.
func SshConsumer(c *gin.Context)  {
	middleware.WsTokenAuthMiddleware("host-console", c)

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
		logging.Error("connect host err: %v", err)
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

func SyncApp(c *gin.Context)  {
	middleware.WsTokenAuthMiddleware("config-app-sync", c)

	//升级get请求为webSocket协议
	wsConn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		wsConn.WriteMessage(websocket.TextMessage, []byte("cant upgrade websocket connection:") )
		return
	}
	defer wsConn.Close()


	id 	:= c.Param("id")
	var hostapp []models.HostApp
	var app models.App
	var env models.ConfigEnv

	models.DB.Model(&models.HostApp{}).
		Where("aid = ?", id).
		Where("status = 0").
		Find(&hostapp)

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
			util.WsWriteMessage("######\r\n开始初始化项目到主机 " + host.Name, wsConn)
			desDir 	:= util.ReturnSyncRunDir(app.ID)
			desDir 	= desDir + "/" + app.Name

			// 建立 sftp 传文件
			util.WsWriteMessage("准备同步初始化文件到目标主机！", wsConn)

			clientConfig, _ := util.ReturnClientConfig(host.Username, "")
			hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
			Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
			if e 	!= nil {
				util.WsWriteMessage(e.Error(), wsConn)
				return
			}

			SftpClient, e := util.GetSftpClient(Scli)
			if e 	!= nil {
				util.WsWriteMessage(e.Error(), wsConn)
				return
			}

			e 	= SyncAppHost(app, desDir,  SftpClient, Scli, wsConn)
			if e!= nil {
				//util.WsWriteMessage(e.Error(), wsConn)
				return
			}

			// 初始化完成 修改状态
			e 	= models.DB.Model(&models.HostApp{}).
				Where("id = ?", v.ID).
				Update("status", 1).Error
			if e!= nil {
				util.WsWriteMessage(err.Error(), wsConn)
				return
			}
		}
	}

	util.WsWriteMessage("项目初始化完毕!", wsConn)
}

func SyncAppHost(app models.App, path string, SftpClient *sftp.Client, Scli *ssh.Client, ws *websocket.Conn) error {
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
		util.WsWriteMessage("开始上传BackendJar项目通用系统依赖文件 jarFuncs", ws)
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
		e   	= util.PutDirectory(SftpClient, path, "/data/webapps/" + app.Name)
		if e 	!= nil {
			util.WsWriteMessage(e.Error(), ws)
			return e
		}

		pathDir := "/data/webapps/" + app.Name
		cmdstr 	:= "chown -R tomcat:tomcat " + pathDir
		util.WsWriteMessage("改变目标主机项目目录权限，执行： "  + cmdstr, ws)
		out, e 	:= util.ExecuteCmdRemote(cmdstr,  Scli)
		if e 	!= nil {
			util.WsWriteMessage(e.Error() + "\r\n请检查是否存在tomcat用户", ws)
			return e
		}
		util.WsWriteMessage(out, ws)


		pathDir = "/data/webapps/" + app.Name + "/bin/*"
		cmdstr 	= "chmod +x " + pathDir
		util.WsWriteMessage("脚本添加执行权限执行：" + cmdstr, ws)

		out, e 	= util.ExecuteCmdRemote(cmdstr, Scli)
		if e 	!= nil {
			util.WsWriteMessage(e.Error(), ws)
			return e
		}
		util.WsWriteMessage(out, ws)

	case 2:
		fmt.Println(2)
		return nil
	}
	return nil
}

func wsWriteMsg(ws *websocket.Conn, step int, key, message string)  {
	var msg  help.Msg
	msg.Step = strconv.Itoa(step)
	msg.Key  = key
	message = "\r\n" + message
	msg.Data = message
	e := ws.WriteMessage(websocket.TextMessage, []byte(util.JSONMarshalToString(msg)))
	if e != nil {
		logging.Error(e)
	}
}

func sendResNotify(id, isOk, dtype int, name, strerr string)  {
	if isOk == 1 {
		msg := "应用 ：" + name  + "发布成功，请验证！"
		if dtype == 4 {
			msg = "应用 ：" + name  + "回滚成功，请验证！"
		}
		e := models.DingtalkSentChannel(id, msg, models.DingList, false)
		if e != nil {
			models.MakeNotify(1, 2, "发送通知异常" + e.Error(), msg, "" )
		}
	} else {
		msg := "应用 ：" + name  + "发布失败，---" + strerr
		if dtype == 4 {
			msg = "应用 ：" + name  + "回滚失败，---" + strerr
		}
		e := models.DingtalkSentChannel(id, msg, models.DingList, false)
		if e != nil {
			models.MakeNotify(1, 2, "发送通知异常" + e.Error(), msg,"" )
		}
	}
}
// 发布逻辑
func AppDeployRedo(c *gin.Context)  {
	middleware.WsTokenAuthMiddleware("deploy-app-redo", c)

	id 			:= c.Param("id")
	userid,_ 	:= c.Get("Uid")
	uid 		:= userid.(int)

	//升级get请求为webSocket协议
	wsConn, e := upgrader.Upgrade(c.Writer, c.Request, nil)
	if e != nil {
		msg := "Cannot upgrade connection \r\n"
		wsWriteMsg(wsConn, 1, "local", msg)
		return
	}
	defer wsConn.Close()

	var helper help.Helper
	helper.Rdbkey = models.DeployInfoKey + id
	helper.Wsconn = wsConn

	msg  := "完成\r\n" + util.HumanNowTime() + " 发布准备...  "
	helper.WsWriteMsg("local", 1, msg, "")
	helper.SendSetup("local", 1, msg)

	// 执行发布逻辑
	var deploy models.AppDeploy
	models.DB.Model(&models.AppDeploy{}).
		Where("id = ?", id).
		Where("status = ? or status = ?", ReviewSuccess, UndoNeedDeploy).
		Find(&deploy)

	// 检查发布模板信息
	var det models.DeployExtend
	models.DB.Model(&models.DeployExtend{}).
		Where("dtid = ?", deploy.Tid).
		Find(&det)

	// 检查项目信息
	var app models.App
	models.DB.Model(&models.App{}).
		Where("id = ?", det.Aid).
		Where("active = 1").
		Find(&app)

	GlobalVersion = id + "_" + strconv.Itoa(deploy.Tid) + "_" + util.HumanNowDate()

	if app.DeployType != 0 {
		// 检查发布类型
		switch app.Tid {
		// backend jar
		case 1:
			e := BackendJarDeploy(id, deploy, app, det, helper)
			models.Rdb.Expire(helper.Rdbkey,   7 * 24 * time.Hour)
			if e != nil {
				status := DeployFail
				if deploy.Status == 4 {
					status = UndoFail
				}

				if det.NotifyId > 0 {
					sendResNotify(det.NotifyId, 0, deploy.Status, app.Name, e.Error())
				}

				// 修改数据库信息 stauts
				e 	= models.DB.Model(&models.AppDeploy{}).
					Where("id = ?", id).
					Update("status", status).
					Update("deploy", uid).Error

				panic(e)
			}

			status := DeploySuccess
			if deploy.Status == 4 {
				status = UndoSuccess
			}

			// 修改数据库信息 stauts 5
			e 	= models.DB.Model(&models.AppDeploy{}).
				Where("id = ?", id).
				Update("status", status).
				Update("version", GlobalVersion).
				Update("deploy", uid).Error

			if det.NotifyId > 0 {
				if det.NotifyId > 0 {
					sendResNotify(det.NotifyId, 1, deploy.Status, app.Name, "")
				}
			}

			//BackendJarDeploy(id, app, det, deploy, uid, wsConn)
		case 2:
			fmt.Println(2)
		}
	} else {
		// 通用方式
		if det.Extend == 1 {
			e := CommonDeploy(id, deploy, det, helper)
			models.Rdb.Expire(helper.Rdbkey,   7 * 24 * time.Hour)
			if e != nil {
				status := DeployFail
				if deploy.Status == 4 {
					status = UndoFail
				}

				if det.NotifyId > 0 {
					sendResNotify(det.NotifyId, 0, deploy.Status, app.Name, e.Error())
				}

				// 修改数据库信息 stauts -3
				e 	= models.DB.Model(&models.AppDeploy{}).
					Where("id = ?", id).
					Update("status", status).
					Update("deploy", uid).Error

				panic(e)
			}

			status := DeploySuccess
			if deploy.Status == 4 {
				status = UndoSuccess
			}
			// 修改数据库信息 stauts 5
			e 	= models.DB.Model(&models.AppDeploy{}).
				Where("id = ?", id).
				Update("status", status).
				Update("version", GlobalVersion).
				Update("deploy", uid).Error

			if det.NotifyId > 0 {
				if det.NotifyId > 0 {
					sendResNotify(det.NotifyId, 1, deploy.Status, app.Name, "")
				}
			}
		}

		if det.Extend == 2 {
			e := CustomDeploy(id, deploy, det, helper)
			models.Rdb.Expire(helper.Rdbkey,   7 * 24 * time.Hour)
			if e != nil {
				status := DeployFail
				if deploy.Status == 4 {
					status = UndoFail
				}

				if det.NotifyId > 0 {
					sendResNotify(det.NotifyId, 0, deploy.Status, app.Name, e.Error())
				}

				// 修改数据库信息 stauts -3
				e 	= models.DB.Model(&models.AppDeploy{}).
					Where("id = ?", id).
					Update("status", status).
					Update("deploy", uid).Error

				panic(e)
			}

			status := DeploySuccess
			if deploy.Status == 4 {
				status = UndoSuccess
			}
			// 修改数据库信息 stauts 5
			e 	= models.DB.Model(&models.AppDeploy{}).
				Where("id = ?", id).
				Update("status", status).
				Update("version", GlobalVersion).
				Update("deploy", uid).Error

			if det.NotifyId > 0 {
				if det.NotifyId > 0 {
					sendResNotify(det.NotifyId, 1, deploy.Status, app.Name, "")
				}
			}
		}
	}
}

func BackendJarDeploy(id string, deploy models.AppDeploy, app models.App, det models.DeployExtend, helper help.Helper) error {
	// 用户自定义环境变量
	env := make(map[string]string)

	// 回滚
	if deploy.Status == 4 {
		msg	:= "完成\r\n" + util.HumanNowTime() + " 回滚发布...        跳过"
		helper.WsWriteMsg("local", 6, msg, "")
		helper.SendSetup("local", 6, msg)

		var olddeploy models.AppDeploy
		models.DB.Model(&models.AppDeploy{}).
			Where("id < ?", id).
			Where("tid = ?", deploy.Tid).
			Where("status = ? or status = ? ", DeploySuccess, UndoSuccess).
			Find(&olddeploy)
		GlobalVersion = olddeploy.Version
		// 主机回滚
		hosts := strings.Split(det.HostIds,",")

		for _, v := range hosts {
			e := BackendJarUndoHost(v, helper, det, app)
			if e != nil {
				return  e
			}
		}
		return  nil
	}

	if det.CustomEnvs != "" {
		if e := json.Unmarshal([]byte(det.CustomEnvs), &env); e != nil {
			msg := "载入用户变量出错:" + e.Error()
			helper.WsWriteMsg("local", 2, msg, "")
			return e
		}
	}

	treeHash := ""
	repoDir := util.ReturnGitLocalPath(det.Aid, det.RepoUrl)
	appRoot := util.ReturnAppGitRoot(det.Aid)
	if deploy.GitType == "branch" {
		treeHash = strings.Split(deploy.Commit," ")[0]
	} else {
		treeHash = strings.Split(deploy.TagBranch," ")[0]
	}

	msg  := "cd " + appRoot + " && rm -rf " + id + "_*"
	if e := helper.Local(msg, 2, "local", env); e != nil {
		return e
	}

	if det.PreCode != "" {
		msg = util.HumanNowTime() + " 检出前任务...  "
		helper.SendSetup("local", 2, msg)
		helper.WsWriteMsg("local", 2, msg, "" )
		if e :=  helper.Local("cd /tmp && " + det.PreCode, 2, "local", env); e != nil {
			return e
		}
	}

	msg = "完成\r\n" + util.HumanNowTime() + " 执行检出...   "
	helper.SendSetup("local", 3, msg)
	helper.WsWriteMsg("local", 3, msg, "")
	cmd := "cd "  + repoDir + " && git archive --prefix=" + GlobalVersion + "/ " + treeHash+ " | (cd .. && tar xf -)"

	if e := helper.Local(cmd, 3, "local", env); e != nil {
		return e
	}

	msg = "完成\r\n"
	helper.SendSetup("local", 3, msg)
	helper.WsWriteMsg( "local", 3, msg, "")

	if det.PostCode != "" {
		msg = util.HumanNowTime() + " 检出后任务...   "
		helper.SendSetup("local", 4, msg)
		helper.WsWriteMsg("local", 4, msg, "")
		cmd = "cd " + appRoot + "/"+ GlobalVersion + " && " + det.PostCode
		if e :=  helper.Local(cmd, 4, "local", env); e != nil {
			return e
		}
	}

	msg = util.HumanNowTime() + " 执行打包...   "
	helper.SendSetup("local", 5, msg)
	helper.WsWriteMsg("local",5,  msg, "")

	rule := make(map[string]string)
	if e := json.Unmarshal([]byte(det.FilterRule), &rule); e != nil {
		msg := "载入打包规则出错:" + e.Error()
		helper.WsWriteMsg("local", 5, msg, "error")
		helper.SendError("local", msg)
		return e
	}

	exclude := ""
	contain	:= GlobalVersion
	if len(rule) > 1 {
		if rule["type"] == "exclude" {
			for _, v := range strings.Split(rule["data"],"\n") {
				exclude = exclude + " --exclude=" + v
			}
		} else {
			contain = ""
			for _, v := range strings.Split(rule["data"],"\n") {
				contain = contain + " " + GlobalVersion + "/" + v
			}
		}
	}

	localFile :=  GlobalVersion + "/target/" + app.Name + "-" +
		det.Tag + ".jar"

	if !util.Exists(appRoot + "/" + localFile) {
		localFile 	= GlobalVersion  + "/" + app.Name +
			"/target/" + app.Name + "-" + det.Tag + ".jar"
	}

	cmd = "cd " + appRoot + " && cp " + localFile  + " " + GlobalVersion + " && tar zcf " + GlobalVersion + ".tar.gz "  + GlobalVersion + "/" + app.Name + "-" + det.Tag + ".jar"
	if e := helper.Local(cmd, 5, "local", env); e != nil {
		return e
	}

	helper.SendSetup("local", 6, "完成")
	helper.WsWriteMsg("local",6,  "完成", "")

	// 主机发布
	hosts := strings.Split(det.HostIds,",")
	for _, v := range hosts {
		e := BackendJarDeployHost(v, helper, det, app)
		if e != nil {
			return  e
		}
	}

	return  nil
}

func BackendJarUndoHost(hid string, helper help.Helper, det models.DeployExtend, app models.App) error {
	msg := util.HumanNowTime() + " 回滚数据准备...   "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid,1,  msg, "")

	msg = " 完成 \r\n"
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid, 1, msg, "")

	var host models.Host
	models.DB.Model(&models.Host{}).
		Where("id = ?", hid).
		Where("status = 1").
		Find(&host)

	if host.ID == 0 {
		helper.WsWriteMsg(hid, 1, "No such host", "error")
		helper.SendError(hid, "No such host")
		return fmt.Errorf("No such host")
	}

	clientConfig, _ := util.ReturnClientConfig(host.Username, "")
	hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
	Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
	if e != nil {
		helper.WsWriteMsg(hid, 1, e.Error() + "\r\n", "error")
		helper.SendError(hid, e.Error() + "\r\n")
		return e
	}

	// pre host
	if det.PreDeploy != "" {
		msg = util.HumanNowTime() + " 发布前任务...   "
		helper.SendSetup(hid,2, msg)
		helper.WsWriteMsg(hid, 2, msg, "")
		cmd := "cd " + det.DstRepo + " && " + det.PreDeploy
		if e := helper.Remote(hid,2, Scli, cmd); e != nil {
			return e
		}
	}

	// deploy
	msg = util.HumanNowTime() + " 执行发布...   "
	helper.SendSetup(hid, 3, msg)
	helper.WsWriteMsg(hid, 3, msg, "")

	appDir 		:= " /data/webapps/" + app.Name + "/lib"
	cmd := "cp -f " + det.DstRepo + "/" + GlobalVersion + "/" + app.Name + "-" + det.Tag + ".jar "  + appDir

	if e := helper.Remote(hid,3, Scli, cmd); e != nil {
		return e
	}
	msg = util.HumanNowTime() + " 完成  "
	helper.SendSetup(hid, 3, msg)
	helper.WsWriteMsg(hid, 3, msg, "")

	// post host
	if det.PostDeploy != "" {
		msg = util.HumanNowTime() + " 发布后任务...   "
		helper.SendSetup(hid,4, msg)
		helper.WsWriteMsg(hid, 4, msg, "")
		cmd = "cd " + det.DstRepo + " && " + det.PostDeploy
		if e := helper.Remote(hid,4, Scli, cmd); e != nil {
			return e
		}
	}

	msg = "\r\n" + util.HumanNowTime() + " *** 发布成功 ***"
	helper.SendSetup(hid, 5, msg)
	helper.WsWriteMsg(hid,5, msg, "")

	return nil
}

func BackendJarDeployHost(hid string, helper help.Helper, det models.DeployExtend, app models.App) error {
	msg := util.HumanNowTime() + " 数据准备...   "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid,1,  msg, "")

	var host models.Host
	models.DB.Model(&models.Host{}).
		Where("id = ?", hid).
		Where("status = 1").
		Find(&host)

	if host.ID == 0 {
		helper.WsWriteMsg(hid, 1, "No such host \r\n", "error")
		helper.SendError(hid, "No such host \r\n")
		return fmt.Errorf("No such host")
	}

	clientConfig, _ := util.ReturnClientConfig(host.Username, "")
	hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
	Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
	if e != nil {
		helper.WsWriteMsg(hid, 1, e.Error() + "\r\n", "error")
		helper.SendError(hid, e.Error() + "\r\n")
		return e
	}

	cmd := "mkdir -p " + det.DstRepo + "  && mkdir -p " + det.DstDir
	_,e = util.ExecuteCmdRemote(cmd, Scli)

	// 保留备份
	index	:= strings.Join(strings.Split(GlobalVersion, "_")[0:2],"_") + "_*"
	cmd 	= "cd " + det.DstRepo + " && rm -rf "+ GlobalVersion + " && ls -1tr -d " + index + " 2> /dev/null | head -n -" + strconv.Itoa(det.Versions) + " | xargs -d '\\n' rm -rf --"
	if e := helper.Remote(hid, 1, Scli, cmd); e != nil {
		return e
	}

	// transfer files
	sourceDir 		:= util.ReturnAppGitRoot(det.Aid) + "/"
	sourceTarFile 	:= GlobalVersion + ".tar.gz"
	SftpClient, e 	:= util.GetSftpClient(Scli)
	e 	= util.PutFile(SftpClient, sourceDir+sourceTarFile, det.DstRepo)
	if e != nil {
		helper.WsWriteMsg(hid, 1, e.Error() + "\r\n", "error")
		helper.SendError(hid, e.Error() + "\r\n")
		return e
	}

	cmd	= "cd " + det.DstRepo + " && tar xf " + sourceTarFile + " && rm -f " + index + "*.tar.gz"
	if e := helper.Remote(hid, 1, Scli, cmd); e != nil {
		return e
	}
	msg = util.HumanNowTime() + " 完成 "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid, 1, msg, "")

	// pre host
	if det.PreDeploy != "" {
		msg = util.HumanNowTime() + " 发布前任务...   "
		helper.SendSetup(hid,2, msg)
		helper.WsWriteMsg(hid, 2, msg, "")
		cmd = "cd " + det.DstRepo + " && " + det.PreDeploy
		if e := helper.Remote(hid,2, Scli, cmd); e != nil {
			return e
		}
	}

	// deploy
	msg = util.HumanNowTime() + " 执行发布...   "
	helper.SendSetup(hid, 3, msg)
	helper.WsWriteMsg(hid, 3, msg, "")
	appDir 		:= " /data/webapps/" + app.Name + "/lib"
	appTagName 	:= app.Name +"-" + det.Tag + ".jar "
	cmd = "cp -f  " + det.DstRepo + "/" + GlobalVersion + "/" + appTagName + appDir  + " && ln -sfn " + appDir + "/" + appTagName + appDir + "/" + app.Name + ".jar"

	if e := helper.Remote(hid,3, Scli, cmd); e != nil {
		return e
	}
	msg = util.HumanNowTime() + " 完成  "
	helper.SendSetup(hid, 3, msg)
	helper.WsWriteMsg(hid, 3, msg, "")

	// post host
	if det.PostDeploy != "" {
		msg = util.HumanNowTime() + " 发布后任务...  "
		helper.SendSetup(hid,4, msg)
		helper.WsWriteMsg(hid, 4, msg, "")
		cmd = "cd " + det.DstRepo + " && " + det.PostDeploy
		if e := helper.Remote(hid,4, Scli, cmd); e != nil {
			fmt.Println(e.Error())
			return e
		}
	}

	msg = "\r\n" + util.HumanNowTime() + " *** 发布成功 ***"
	helper.SendSetup(hid, 5, msg)
	helper.WsWriteMsg(hid,5, msg, "")

	return nil
}

func CommonDeploy(id string, deploy models.AppDeploy, det models.DeployExtend, helper help.Helper) error {
	// 用户自定义环境变量
	env := make(map[string]string)
	// 回滚
	if deploy.Status == 4 {
		msg	:= "完成\r\n" + util.HumanNowTime() + " 回滚发布...        跳过"
		helper.WsWriteMsg("local", 6, msg, "")
		helper.SendSetup("local", 6, msg)

		var olddeploy models.AppDeploy
		models.DB.Model(&models.AppDeploy{}).
			Where("id < ?", id).
			Where("tid = ?", deploy.Tid).
			Where("status = ? or status = ? ", DeploySuccess, UndoSuccess).
			Find(&olddeploy)
		GlobalVersion = olddeploy.Version

		// 主机回滚
		hosts := strings.Split(det.HostIds,",")
		for _, v := range hosts {
			e := CommonUndoHost(v, helper, det)
			if e != nil {
				return  e
			}
		}

		return  nil
	}

	if det.CustomEnvs != "" {
		if e := json.Unmarshal([]byte(det.CustomEnvs), &env); e != nil {
			msg := "载入用户变量出错:" + e.Error() + "\r\n"
			helper.WsWriteMsg("local", 2, msg, "")
			return e
		}
	}

	treeHash := ""
	repoDir := util.ReturnGitLocalPath(det.Aid, det.RepoUrl)
	appRoot := util.ReturnAppGitRoot(det.Aid)
	if deploy.GitType == "branch" {
		treeHash = strings.Split(deploy.Commit," ")[0]
	} else {
		treeHash = strings.Split(deploy.TagBranch," ")[0]
	}

	msg  := "cd " + appRoot + " && rm -rf " + id + "_*"
	if e := helper.Local(msg, 2, "local", env); e != nil {
		return e
	}

	if det.PreCode != "" {
		msg = util.HumanNowTime() + " 检出前任务...   "
		helper.SendSetup("local", 2, msg)
		helper.WsWriteMsg("local", 2, msg, "" )
		if e :=  helper.Local("cd /tmp && " + det.PreCode, 2, "local", env); e != nil {
			return e
		}
	}

	msg = "完成\r\n" + util.HumanNowTime() + " 执行检出...   "
	helper.SendSetup("local", 3, msg)
	helper.WsWriteMsg("local", 3, msg, "")
	cmd := "cd "  + repoDir + " && git archive --prefix=" + GlobalVersion + "/ " + treeHash+ " | (cd .. && tar xf -)"

	if e := helper.Local(cmd, 3, "local", env); e != nil {
		return e
	}

	msg = "\r\n完成\r\n"
	helper.SendSetup("local", 3, msg)
	helper.WsWriteMsg( "local", 3, msg, "")

	if det.PostCode != "" {
		msg = util.HumanNowTime() + " 检出后任务...   "
		helper.SendSetup("local", 4, msg)
		helper.WsWriteMsg("local", 4, msg, "")
		cmd = "cd " + appRoot + "/"+ GlobalVersion + " && " + det.PostCode
		if e :=  helper.Local(cmd, 4, "local", env); e != nil {
			return e
		}
	}

	msg = util.HumanNowTime() + " 执行打包...   "
	helper.SendSetup("local", 5, msg)
	helper.WsWriteMsg("local",5,  msg, "")

	rule := make(map[string]string)
	if e := json.Unmarshal([]byte(det.FilterRule), &rule); e != nil {
		msg := "载入打包规则出错:" + e.Error() + "\r\n"
		helper.WsWriteMsg("local", 5, msg, "error")
		helper.SendError("local", msg)
		return e
	}

	exclude := ""
	contain	:= GlobalVersion
	if len(rule) > 1 {
		if rule["type"] == "exclude" {
			for _, v := range strings.Split(rule["data"],"\n") {
				exclude = exclude + " --exclude=" + v
			}
		} else {
			contain = ""
			for _, v := range strings.Split(rule["data"],"\n") {
				contain = contain + " " + GlobalVersion + "/" + v
			}
		}
	}

	cmd = "cd " + appRoot + " && tar zcf " + GlobalVersion + ".tar.gz " + exclude + " " + contain
	if e := helper.Local(cmd, 5, "local", env); e != nil {
		return e
	}
	helper.SendSetup("local", 6, "完成 ")
	helper.WsWriteMsg("local",6,  "完成 ", "")

	// 主机发布
	hosts := strings.Split(det.HostIds,",")
	for _, v := range hosts {
		e := CommonDeployHost(v, helper, det)
		if e != nil {
			return  e
		}
	}

	return  nil
}

func CommonUndoHost(hid string, helper help.Helper, det models.DeployExtend) error {
	msg := util.HumanNowTime() + " 回滚数据准备...   "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid,1,  msg, "")

	msg = " 完成 "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid, 1, msg, "")

	var host models.Host
	models.DB.Model(&models.Host{}).
		Where("id = ?", hid).
		Where("status = 1").
		Find(&host)

	if host.ID == 0 {
		helper.WsWriteMsg(hid, 1, "No such host \r\n", "error")
		helper.SendError(hid, "No such host \r\n")
		return fmt.Errorf("No such host")
	}

	clientConfig, _ := util.ReturnClientConfig(host.Username, "")
	hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
	Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
	if e != nil {
		helper.WsWriteMsg(hid, 1, e.Error() + "\r\n", "error")
		helper.SendError(hid, e.Error() + "\r\n")
		return e
	}

	// pre host
	if det.PreDeploy != "" {
		msg = util.HumanNowTime() + " 发布前任务...   "
		helper.SendSetup(hid,2, msg)
		helper.WsWriteMsg(hid, 2, msg, "")
		cmd := "cd " + det.DstRepo + " && " + det.PreDeploy
		if e := helper.Remote(hid,2, Scli, cmd); e != nil {
			return e
		}
	}

	// deploy
	msg = util.HumanNowTime() + " 执行发布...   "
	helper.SendSetup(hid, 3, msg)
	helper.WsWriteMsg(hid, 3, msg, "")
	cmd := "rm -rf " + det.DstDir + " && ln -sfn " + det.DstRepo + "/" + GlobalVersion + " " + det.DstDir
	if e := helper.Remote(hid,3, Scli, cmd); e != nil {
		return e
	}
	msg = util.HumanNowTime() + " 完成  "
	helper.SendSetup(hid, 3, msg)
	helper.WsWriteMsg(hid, 3, msg, "")

	// post host
	if det.PostDeploy != "" {
		msg = util.HumanNowTime() + " 发布后任务...   "
		helper.SendSetup(hid,4, msg)
		helper.WsWriteMsg(hid, 4, msg, "")
		cmd = "cd " + det.DstRepo + " && " + det.PostDeploy
		if e := helper.Remote(hid,4, Scli, cmd); e != nil {
			return e
		}
	}

	msg = "\r\n" + util.HumanNowTime() + " *** 发布成功 ***"
	helper.SendSetup(hid, 5, msg)
	helper.WsWriteMsg(hid,5, msg, "")

	return nil
}

func CommonDeployHost(hid string, helper help.Helper, det models.DeployExtend) error {
	msg := util.HumanNowTime() + " 数据准备...   "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid,1,  msg, "")

	var host models.Host
	models.DB.Model(&models.Host{}).
		Where("id = ?", hid).
		Where("status = 1").
		Find(&host)

	if host.ID == 0 {
		helper.WsWriteMsg(hid, 1, "No such host \r\n", "error")
		helper.SendError(hid, "No such host \r\n")
		return fmt.Errorf("No such host")
	}

	clientConfig, _ := util.ReturnClientConfig(host.Username, "")
	hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
	Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
	if e != nil {
		helper.WsWriteMsg(hid, 1, e.Error() + "\r\n", "error")
		helper.SendError(hid, e.Error() + "\r\n")
		return e
	}

	cmd := "mkdir -p " + det.DstRepo + "  && [ -e " + det.DstDir + " ] && [ ! -L " + det.DstDir + " ]"
	_,e = util.ExecuteCmdRemote(cmd, Scli)
	if e == nil {
		msg = "检测到该主机的发布目录 " + det.DstDir + " 已存在，为了数据安全请自行备份后删除该目录，发布平台 将会创建并接管该目录。"
		helper.WsWriteMsg(hid, 1, msg, "error")
		helper.SendError(hid, msg)
		return fmt.Errorf(msg)
	}

	// 保留备份
	index	:= strings.Join(strings.Split(GlobalVersion, "_")[0:2],"_") + "_*"
	cmd 	= "cd " + det.DstRepo + " && rm -rf "+ GlobalVersion + " && ls -1tr -d " + index + " 2> /dev/null | head -n -" + strconv.Itoa(det.Versions) + " | xargs -d '\\n' rm -rf --"
	if e := helper.Remote(hid, 1, Scli, cmd); e != nil {
		return e
	}

	// transfer files
	sourceDir 		:= util.ReturnAppGitRoot(det.Aid) + "/"
	sourceTarFile 	:= GlobalVersion + ".tar.gz"
	SftpClient, e 	:= util.GetSftpClient(Scli)
	e 	= util.PutFile(SftpClient, sourceDir+sourceTarFile, det.DstRepo)
	if e != nil {
		helper.WsWriteMsg(hid, 1, e.Error() + "\r\n", "error")
		helper.SendError(hid, e.Error() + "\r\n")
		return e
	}

	cmd	= "cd " + det.DstRepo + " && tar xf " + sourceTarFile + " && rm -f " + index + "*.tar.gz"
	if e := helper.Remote(hid, 1, Scli, cmd); e != nil {
		return e
	}
	msg = util.HumanNowTime() + " 完成 "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid, 1, msg, "")

	// pre host
	if det.PreDeploy != "" {
		msg = util.HumanNowTime() + " 发布前任务...   "
		helper.SendSetup(hid,2, msg)
		helper.WsWriteMsg(hid, 2, msg, "")
		cmd = "cd " + det.DstRepo + " && " + det.PreDeploy
		if e := helper.Remote(hid,2, Scli, cmd); e != nil {
			return e
		}
	}

	// deploy
	msg = util.HumanNowTime() + " 执行发布...   "
	helper.SendSetup(hid, 3, msg)
	helper.WsWriteMsg(hid, 3, msg, "")
	cmd = "rm -rf " + det.DstDir + " && ln -sfn " + det.DstRepo + "/" + GlobalVersion + " " + det.DstDir
	if e := helper.Remote(hid,3, Scli, cmd); e != nil {
		return e
	}
	msg = util.HumanNowTime() + " 完成  "
	helper.SendSetup(hid, 3, msg)
	helper.WsWriteMsg(hid, 3, msg, "")

	// post host
	if det.PostDeploy != "" {
		msg = util.HumanNowTime() + " 发布后任务...  "
		helper.SendSetup(hid,4, msg)
		helper.WsWriteMsg(hid, 4, msg, "")
		cmd = "cd " + det.DstRepo + " && " + det.PostDeploy
		if e := helper.Remote(hid,4, Scli, cmd); e != nil {
			return e
		}
	}

	msg = "\r\n" + util.HumanNowTime() + " *** 发布成功 ***"
	helper.SendSetup(hid, 5, msg)
	helper.WsWriteMsg(hid,5, msg, "")

	return nil
}

func CustomDeploy(id string, deploy models.AppDeploy, det models.DeployExtend, helper help.Helper) error {
	// 用户自定义环境变量
	env := make(map[string]string)
	// 用户自定义发布方式回滚 相当于重新走一遍发布流程。
	step := 2
	if det.PreCode != "" {
		tmpArr := strings.Split(det.PreCode,"|")
		for _,v := range tmpArr {
			var cdp CustomDeployParse
			json.Unmarshal([]byte(v), &cdp)

			msg := util.HumanNowTime() + " 开始执行 ： " + cdp.Title
			helper.SendSetup("local", step, msg)
			helper.WsWriteMsg("local", step, msg, "" )

			cmd := "cd /tmp && " + cdp.Data
			if e := helper.Local(cmd, step, "local", env); e != nil {
				return e
			}
			step++
		}
		helper.SendSetup("local", step, "完成 ")
		helper.WsWriteMsg("local",step,  "完成 ", "")
	}

	if det.PreDeploy != "" {
		// 主机发布
		hosts := strings.Split(det.HostIds,",")
		for _, v := range hosts {
			e := CustomDeployHost(v, helper, det)
			if e != nil {
				return  e
			}
		}
	} else {
		msg := util.HumanNowTime() + " ** 发布成功 **"
		helper.SendSetup("local", step, msg)
		helper.WsWriteMsg("local",step, msg, "")
	}


	return  nil
}

func CustomDeployHost(hid string, helper help.Helper, det models.DeployExtend) error {
	msg := util.HumanNowTime() + " 数据准备...   "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid,1,  msg, "")

	var host models.Host
	models.DB.Model(&models.Host{}).
		Where("id = ?", hid).
		Where("status = 1").
		Find(&host)

	if host.ID == 0 {
		helper.WsWriteMsg(hid, 1, "No such host \r\n", "error")
		helper.SendError(hid, "No such host \r\n")
		return fmt.Errorf("No such host")
	}

	clientConfig, _ := util.ReturnClientConfig(host.Username, "")
	hostIp 			:= host.Addres + ":" + strconv.Itoa(host.Port)
	Scli, e 		:= util.GetSshClient(hostIp, clientConfig)
	if e != nil {
		helper.WsWriteMsg(hid, 1, e.Error() + "\r\n", "error")
		helper.SendError(hid, e.Error() + "\r\n")
		return e
	}

	msg = util.HumanNowTime() + " 完成 "
	helper.SendSetup(hid, 1, msg)
	helper.WsWriteMsg(hid, 1, msg, "")

	step 	:= 2
	tmpArr 	:= strings.Split(det.PreDeploy,"|")
	for _,v := range tmpArr {
		var cdp CustomDeployParse
		json.Unmarshal([]byte(v), &cdp)

		msg := util.HumanNowTime() + " 开始执行 ： " + cdp.Title
		helper.SendSetup(hid, step, msg)
		helper.WsWriteMsg(hid, step, msg, "" )

		cmd := "cd /tmp && " + cdp.Data
		if e := helper.Remote(hid, step, Scli, cmd); e != nil {
			return e
		}
		step++
	}
	msg = util.HumanNowTime() + " ** 发布成功 **"
	helper.SendSetup(hid, step, msg)
	helper.WsWriteMsg(hid,step, msg, "")

	return nil
}