package help

import (
	"api/models"
	"api/pkg/logging"
	"api/pkg/util"
	"bufio"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/ssh"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type Helper struct {
	Rdbkey 		string
	Wsconn   	*websocket.Conn
}

type Msg struct {
	Key   	string
	Status 	string
	Step	string
	Data 	string
}

func (h *Helper) WsWriteMsg(key string, step int, message, status string)  {
	var msg  Msg
	msg.Step = strconv.Itoa(step)
	msg.Key  = key
	msg.Status = "info"
	if status != "" {
		msg.Status = status
	}

	message  = message + "\r\n"
	msg.Data = message
	e := h.Wsconn.WriteMessage(websocket.TextMessage, []byte(util.JSONMarshalToString(msg)))
	if e != nil {
		logging.Error(e)
	}
}

func (h *Helper) Send(message Msg)  {
	models.Rdb.LPush(h.Rdbkey, util.JSONMarshalToString(message))
}

func (h *Helper) SendInfo(key, message string)  {
	var msg Msg
	msg.Key		= key
	msg.Status	= "info"
	msg.Data 	= message + "\r\n"
	h.Send(msg)
}

func (h *Helper) SendError(key, message string)  {
	message = "\r\n" + message
	var msg Msg
	msg.Key		= key
	msg.Status	= "error"
	msg.Data 	= message
	h.Send(msg)
}

func (h *Helper) SendSetup(key string, step int, data string)  {
	var msg Msg
	msg.Key		= key
	msg.Step	= strconv.Itoa(step)
	msg.Data 	= data
	h.Send(msg)
}

func (h *Helper) Local(command string, step int, key string, env map[string]string) error {
	command	= "set -e\n" + command

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Env = os.Environ()
	if len(env) >0 {
		for k, v := range env {
			str := k + "=" + v
			cmd.Env = append(cmd.Env, str)
		}
	}

	out, e := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	e = cmd.Start()

	// Create a scanner which scans stdout in a line-by-line fashion
	scanner := bufio.NewScanner(out)

	for scanner.Scan() {
		m := scanner.Text()
		h.SendInfo("local", m)
		h.WsWriteMsg(key, step, m, "")
	}

	if e = cmd.Wait(); e != nil {
		if exiterr, ok := e.(*exec.ExitError); ok {
			exitcode := -1
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitcode = status.ExitStatus()
			}

			h.WsWriteMsg(key, step, "exit code : " +  strconv.Itoa(exitcode), "error")
			h.SendError("local", "exit code : " +  strconv.Itoa(exitcode))

			return e
		}
	}

	return nil
}

func (h *Helper) Remote(key string, step int, cli *ssh.Client, command string) error {
	command	= "set -e\n" + command

	 out, e	:= util.ExecuteCmdRemote(command, cli)
	 if e != nil {
	 	h.WsWriteMsg(key, step, e.Error(), "error")
	 	h.SendError(key, e.Error())
	 	return e
	 }

	 h.SendInfo(key, out)
	 h.WsWriteMsg(key, step, out, "")

	 return nil
}
