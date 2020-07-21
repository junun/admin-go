package util

import (
	"api/models"
	"api/pkg/logging"
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	Scli *ssh.Client
	SftpClient *sftp.Client
)
type Charset string

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

// Conn wraps a net.Conn, and sets a deadline for every read
// and write operation.
type Conn struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *Conn) Read(b []byte) (int, error) {
	err := c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

func (c *Conn) Write(b []byte) (int, error) {
	err := c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	if err != nil {
		return 0, err
	}
	return c.Conn.Write(b)
}

func SSHDialTimeoutClient(network, addr string, config *ssh.ClientConfig, timeout time.Duration) (*ssh.Client, error) {
	conn, err := net.DialTimeout(network, addr, timeout)
	if err != nil {
		return nil, err
	}

	timeoutConn := &Conn{conn, timeout, timeout}
	c, chans, reqs, err := ssh.NewClientConn(timeoutConn, addr, config)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(c, chans, reqs)

	//this sends keepalive packets every 2 seconds
	//there's no useful response from these, so we can just abort if there's an error
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for range t.C {
			_, _, err := client.Conn.SendRequest("keepalive@golang.org", true, nil)
			if err != nil {
				return
			}
		}
	}()
	return client, nil
}

func ReturnClientConfig(username string, password string) (*ssh.ClientConfig, error) {
	var clientConfig *ssh.ClientConfig
	if password != "" {
		clientConfig = &ssh.ClientConfig{
			User:    username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,  //Timeout is the maximum amount of time for the TCP connection to establish.
		}
	} else {
		// 获取私钥
		var set models.Settings
		models.DB.Model(&models.Settings{}).
			Where("name = ? ", "private_key").
			Find(&set)

		// Create the Signer for this private key.
		signer, err := ssh.ParsePrivateKey([]byte(set.Value))
		if err != nil {
			logging.Error("unable to parse private key: %v", err)
			return  nil, err
		}

		clientConfig = &ssh.ClientConfig{
			User: 	username,
			Auth: 	[]ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
			Timeout: 10 * time.Second,  //Timeout is the maximum amount of time for the TCP connection to establish.
		}
	}

	return clientConfig, nil
}


func GetSshClient(hostname string, clientConfig *ssh.ClientConfig ) (*ssh.Client, error ){
	if Scli != nil {
		return Scli, nil
	}

	Scli, err := ssh.Dial("tcp", hostname, clientConfig)
	if err != nil {
		//log.Fatalf("unable to connect: %v", err)
		return  nil, err
	}

	return  Scli, nil
}

func GetSftpClient(sshClient *ssh.Client) (*sftp.Client, error) {
	// create sftp client
	SftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	return SftpClient, nil
}

func executeRuntimeCmd(cmd, hostname string, cli *ssh.Client) {
	session, _ := cli.NewSession()

	defer session.Close()

	cmdArgs := strings.Fields(cmd)

	command := exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	stdout, _ := command.StdoutPipe()
	command.Start()


	oneByte := make([]byte, 100)
	//num := 1

	for {
		_, err := stdout.Read(oneByte)
		if err != nil {
			fmt.Printf(err.Error())
			break
		}
		r := bufio.NewReader(stdout)
		line, _, _ := r.ReadLine()
		fmt.Println(string(line))
		//num = num + 1
		//if num > 3 {
		//	os.Exit(0)
		//}
	}

	command.Wait()
}


func PutFile(sftpClient *sftp.Client, localFilePath string, remotePath string) error {
	srcFile, e := os.Open(localFilePath)
	if e != nil {
		return e
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)
	dstFile, e := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if e != nil {
		return e
	}
	defer dstFile.Close()

	ff, e 	:= ioutil.ReadAll(srcFile)
	if e != nil {
		return e
	}
	_, e	= dstFile.Write(ff)
	if e != nil {
		return e
	}

	return nil
}

func PutDirectory(sftpClient *sftp.Client, localPath string, remotePath string) error {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		return err
	}

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			err := sftpClient.MkdirAll(remoteFilePath)
			if err != nil {
				return err
			}
			PutDirectory(sftpClient, localFilePath, remoteFilePath)
		} else {
			PutFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
		}
	}

	return nil
}



func split(r rune) bool {
	return r == '\n' || r == '\r'
}

func ValidHosh(addres string, port int, username string, password string)  bool {
	clientConfig, _ := ReturnClientConfig(username, password)
	hostIp 			:= addres + ":" + strconv.Itoa(port)

	Scli, e := GetSshClient(hostIp, clientConfig)
	//Scli, err := SSHDialTimeoutClient("tcp", hostIp, clientConfig, 3 * time.Second)

	if e != nil {
		logging.Error("connect host err: %v", e)
		return false
	}
	defer Scli.Close()

	if password != "" {
		keystr 	:= ""
		// 获取公钥
		var set models.Settings
		models.DB.Model(&models.Settings{}).
			Where("name = ? ", "public_key").
			Find(&set)

		if set.ID == 0 {
			key, e := rsa.GenerateKey(rand.Reader, 2048)
			if e != nil {
				logging.Error("Private key cannot be created.", e.Error())
				return false
			}
			pubkey := &key.PublicKey
			privateKey, _ 	:= DumpPrivateKeyBuffer(key)
			publicKey,  _ 	:= DumpPublicKeyBuffer(pubkey)

			setPrivateKey 	:= models.Settings{Name: "private_key", Value: privateKey, Desc: "私钥"}
			setPublicKey 	:= models.Settings{Name: "public_key", Value: publicKey, Desc: "公钥"}
			if e:= models.DB.Create(&setPrivateKey).Error; e!=nil{
				return false
			}
			if e:=  models.DB.Create(&setPublicKey).Error; e!=nil{
				return false
			}

			keystr = publicKey
		} else  {
			keystr = set.Value
		}

		command := "mkdir -p -m 700 ~/.ssh " +
			"&& echo '%v' >> ~/.ssh/authorized_keys " +
			"&& chmod 600 ~/.ssh/authorized_keys"

		keystr, _ = LoadPublicKeyToAuthorizedFormat(keystr)
		_, e := ExecuteCmdRemote(fmt.Sprintf(command, keystr), Scli)
		if e != nil {
			logging.Error("add public key error: %v", e)
			return false
		}
	} else {
		_, e := ExecuteCmdRemote("ping -c 1 127.0.0.1", Scli)
		if e != nil {
			logging.Error("auth fail : %v", e)
			return false
		}
	}



	//dir, _ := os.Getwd()
	//path := dir + "/" + GetIdRsaPath()
	//cmd, err := ioutil.ReadFile(path + "test")
	//
	//m := strings.FieldsFunc(string(cmd), split)
	//for _,c := range m {
	//	if c != "" {
	//		go fmt.Println(executeCmd(c, addres, Scli))
	//	}
	//}

	//fmt.Println(executeCmd("ls -lah" ,addres, Scli))
	//time.Sleep(1 * time.Second)
	//
	//var localFilePath 	= "/Users/angus/Documents/git/go_spug/api/runtime/upload/images"
	//var remoteDir 		= "/tmp/"
	//SftpClient,_ := GetSftpClient(Scli)
	////PutFile(SftpClient,localFilePath, remoteDir)
	//
	//PutDirectory(SftpClient,localFilePath, remoteDir )

	//fmt.Println("time now")
	//fmt.Println(executeCmd(string(cmd),addres, Scli))
	//fmt.Println(executeCmd("cat /proc/cpuinfo \n" +
	//	"who",addres, Scli))
	return true
}

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	defer func() {
	}()
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ExecRuntimeCmdToWs(cmd, path string, ws  *websocket.Conn) error {
	outr, outw, e 	:= os.Pipe()
	if e 	!= nil {
		return e
	}
	defer outr.Close()
	defer outw.Close()

	cmdArgs 	:= strings.Fields(cmd)
	command 	:= exec.Command(cmdArgs[0], cmdArgs[1:len(cmdArgs)]...)
	command.Dir	= path

	stdout, _ 	:= command.StdoutPipe()
	command.Stderr = command.Stdout
	//done := make(chan struct{})
	if e = command.Start(); e != nil {
		return e
	}


	// Create a scanner which scans stdout in a line-by-line fashion
	scanner := bufio.NewScanner(stdout)


	//scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		e 	= ws.WriteMessage(websocket.TextMessage, []byte(m))
		if e != nil {
			ws.Close()
			break
		}
	}

	if e = command.Wait(); e != nil {
		return e
	}

	return nil
}

func ExecRuntimeCmd(command string) (error, string) {
	var message string
	outr, outw, e 	:= os.Pipe()
	if e 	!= nil {
		return e, ""
	}
	defer outr.Close()
	defer outw.Close()

	cmd := exec.Command("/bin/sh", "-c", command)

	stdout, _ 	:= cmd.StdoutPipe()
	cmd.Stderr 	= cmd.Stdout
	if e = cmd.Start(); e != nil {
		return e, ""
	}

	// Create a scanner which scans stdout in a line-by-line fashion
	scanner := bufio.NewScanner(stdout)

	//scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		if message != "" {
			message = message + "\n" + m
		}
		message = message  + m
	}

	return  nil, message
}

func ExecCmdBySshToWs(cmd string, cli *ssh.Client, ws *websocket.Conn)  {
	session, _ := cli.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	e := session.Run(cmd)
	if e != nil {
		WsWriteMessage(e.Error(), ws)
	} else {
		WsWriteMessage("执行成功！", ws)
	}
}

func ExecuteCmdRemote(cmd string, cli *ssh.Client) (string, error) {
	session, _ := cli.NewSession()
	defer session.Close()

	command := "set -e\n" + cmd

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stdoutBuf

	e 	:= session.Run(command)
	if e != nil {
		return stdoutBuf.String(), e
	}
	return stdoutBuf.String(), nil
}