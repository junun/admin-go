package util

import (
	"api/models"
	"api/pkg/logging"
	"api/pkg/setting"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetGitLocalPath() string {
	return setting.AppSetting.GitLocalPath
}

func ReturnAppGitRoot(aid int) string {
	dir, _ 	:= os.Getwd()
	path 	:= dir + "/" + GetGitLocalPath() + strconv.Itoa(aid)

	return path
}

func ReturnAuthConfig() (*ssh.PublicKeys, error) {
	// 获取私钥
	var set models.Settings
	models.DB.Model(&models.Settings{}).
		Where("name = ? ", "private_key").
		Find(&set)
	keystr := ""
	if set.ID == 0 {
		key, e := rsa.GenerateKey(rand.Reader, 2048)
		if e != nil {
			logging.Error("Private key cannot be created.", e.Error())
			return nil, e
		}
		pubkey 			:= &key.PublicKey
		privateKey, _ 	:= DumpPrivateKeyBuffer(key)
		publicKey,  _ 	:= DumpPublicKeyBuffer(pubkey)
		publicKeyStr, _ := LoadPublicKeyToAuthorizedFormat(publicKey)

		setPrivateKey 	:= models.Settings{Name: "private_key", Value: privateKey, Desc: "私钥"}
		setPublicKey 	:= models.Settings{Name: "public_key", Value: publicKeyStr, Desc: "公钥"}

		if e:= models.DB.Create(&setPrivateKey).Error; e!=nil{
			return nil, e
		}

		if e:=  models.DB.Create(&setPublicKey).Error; e!=nil{
			return nil, e
		}

		keystr = privateKey
	} else  {
		keystr = set.Value
	}

	sshAuth, e := ssh.NewPublicKeys(
		"git",
		[]byte(keystr),
		"")

	if e != nil {
		return nil, e
	}

	return  sshAuth, nil
}

func ReturnGitLocalPath(aid int, url string) string {
	dir, _ 		:= os.Getwd()
	ss 			:= strings.Split(url,"/")
	appName 	:= strings.Split(ss[len(ss)-1], ".git")[0]

	path := dir + "/" + GetGitLocalPath() + strconv.Itoa(aid) +  "/" + appName

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	return  path
}

func ReturnGitPlainOpen(aid int, url string) (*git.Repository, error) {
	path 	:= ReturnGitLocalPath(aid, url)
	r, err 	:= git.PlainOpen(path)
	if err 	!= nil {
		GitCLone(aid, url)
		r, err	= git.PlainOpen(path)
		if err 	!= nil {
			return nil, err
		}
	}

	return r, nil
}

func ReturnGitWorkDir(aid int, url string) (*git.Worktree, error) {
	r, err 	:= ReturnGitPlainOpen(aid, url)
	if err 	!= nil {
		return nil, err
	}

	// Get the working directory for the repository
	w, err 	:= r.Worktree()
	if err != nil {
		return nil, err
	}

	return w, nil
}

func GitCLone(aid int, url string) error {
	path := ReturnGitLocalPath(aid, url)

	//sshAuth, err := ssh.NewPublicKeysFromFile(
	//	"git",
	//	setting.AppSetting.GitSshKey, "")

	sshAuth, e := ReturnAuthConfig()
	if e != nil {
		return e
	}

	_, e 		= git.PlainClone(path, false, &git.CloneOptions{
		URL:               	url,
		RecurseSubmodules: 	git.DefaultSubmoduleRecursionDepth,
		Auth:              	sshAuth,
		NoCheckout:			true,
	})

	if e != nil {
		return e
	}

	return nil
}

func GitPull(aid int, url string) error {
	w, e := ReturnGitWorkDir(aid, url)
	if e != nil {
		return e
	}

	//sshAuth, e := ssh.NewPublicKeysFromFile(
	//	"git",
	//	setting.AppSetting.GitSshKey, "")

	sshAuth, e := ReturnAuthConfig()
	if e != nil {
		return e
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	e = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: sshAuth,
		Force: true,
	})

	if e != nil && e.Error() != "already up-to-date" && e.Error() != "non-fast-forward update" {
		return  e
	}

	return nil
}

func GitCheckoutByBranch(aid int, url, branch string) error {
	GitPull(aid, url)

	w, e := ReturnGitWorkDir(aid, url)
	if e != nil {
		return e
	}

	branchStr := fmt.Sprintf("refs/heads/%s", branch)
	b := plumbing.ReferenceName(branchStr)

	e = w.Checkout(&git.CheckoutOptions{
		Create: false,
		Force: true,
		Keep:true,
		Branch: b,
	})

	if e != nil {
		// got an error  - try to create it
		err := w.Checkout(&git.CheckoutOptions{Create: true, Force: false, Branch: b} )
		CheckIfError(err)
	}

	return nil
}

func remoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	}, refs), nil
}

func ReturnGitTagByCommand(aid int, url string) ([]string, error) {
	r, e :=ReturnGitPlainOpen(aid, url)
	if e != nil {
		return nil, e
	}
	h, _ := r.Head()

	// go-git pull cannot update tags
	cmd 		:= exec.Command("git", "pull",  "origin", h.Strings()[0] , "--tags")
	cmd.Dir 	=  ReturnGitLocalPath(aid, url)
	_, e 		= cmd.CombinedOutput()

	if e != nil && e.Error() != "already up-to-date." && e.Error() != "non-fast-forward update" {
		return nil, e
	}

	cmd 		= exec.Command("git", "tag", "--sort=-creatordate", "-n")
	cmd.Dir 	= ReturnGitLocalPath(aid, url)
	out, e 		:= cmd.CombinedOutput()

	if e != nil {
		return nil, e
	}

	if e != nil {
		logging.Error("cmd.Run() failed with %s\n", e)
		return nil, e
	}

	s := string(out)
	arr := strings.Split(s,"\n")

	var max int
	if len(arr) - 1 >= 10 {
		max = 10
	} else {
		max = len(arr)-1
	}

	return arr[0:max], nil
}

func ReturnGitTag(aid int, url string) ([]string, error) {
	// go-git pull cannot update tags
	cmd 		:= exec.Command("git", "pull")
	cmd.Dir 	=  ReturnGitLocalPath(aid, url)
	_, e 		:= cmd.CombinedOutput()
	if e != nil {
		return nil, e
	}

	var list []string
	r, e 	:= ReturnGitPlainOpen(aid, url)
	if e 	!= nil {
		return nil, e
	}

	tagrefs, e := r.Tags()
	if e!= nil {
		return nil, e
	}

	e = tagrefs.ForEach(func(t *plumbing.Reference) error {
		list = append(list, strings.Split(t.Strings()[0],"/")[len(t.Strings())])
		return nil
	})

	return  list, nil
}

func ReturnGitBranch(aid int, url string) ([]string, error) {
	e 	:= GitPull(aid, url)
	if e!= nil {
		return nil, e
	}

	var list []string
	r, e 	:= ReturnGitPlainOpen(aid, url)
	if e 	!= nil {
		return nil, e
	}

	bs, e := remoteBranches(r.Storer)
	e 	= bs.ForEach(func(b *plumbing.Reference) error {
		s 	:=  strings.Split(b.Strings()[0], "/")
		list = append(list, s[len(s)-1])
		return nil
	})

	return  list, nil
}

func GetGitLastTenCommitByBranch(aid int, url, branch string) ([]string, error) {
	e 	:= GitCheckoutByBranch(aid, url, branch)

	if e!= nil {
		return nil, e
	}

	cmd 		:= exec.Command("git", "log", "--pretty=format:%h (%cr)  %ce %s", "-10")
	cmd.Dir 	=  ReturnGitLocalPath(aid, url)
	out, e 		:= cmd.CombinedOutput()

	if e != nil {
		logging.Error("cmd.Run() failed with %s\n", e)
		return nil, e
	}

	s := string(out)

	return strings.Split(s,"\n"), nil
}

func GetGitCommitByCommand(aid int, url string)  {
	path 	:= ReturnGitLocalPath(aid, url)
	r, err	:= git.PlainOpen(path)
	if err != nil {
		fmt.Println(err)
	}

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()

	// ... retrieves the commit history
	since := time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)
	//until := time.Date(2019, 7, 30, 0, 0, 0, 0, time.UTC)
	cIter, err := r.Log(
					&git.LogOptions{
						From: ref.Hash(),
						All: true,
						Since: &since,
					})

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c.Hash, c.Author.Email, c.Author.When, c.Message)
		return nil
	})
}

func GitCheckoutByCommit(aid int, url, commit string) error {
	e 		:= GitPull(aid, url)

	if e 	!= nil && e.Error() != "already up-to-date" {
		return e
	}

	w, e	:= ReturnGitWorkDir(aid, url)
	if e != nil {
		return e
	}

	// 通过缩短的commit hash查找commit hash，用checkout
	cmd 		:= exec.Command("git", "log", "-1", "--pretty=format:%H", commit)
	cmd.Dir 	=  ReturnGitLocalPath(aid, url)
	out, e 		:= cmd.CombinedOutput()

	if e 	!= nil {
		return  e
	}

	s 	:= strings.Split(string(out),"\n")
	LongHash :=	s[0]

	e 	= w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(LongHash),
	})

	if e!= nil {
		return e
	}

	return nil
}

//func FetchVersions(aid int, url string)  {
//
//}