package util

import (
	"api/pkg/setting"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetGitLocalPath() string {
	return setting.AppSetting.GitLocalPath
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

	sshAuth, err := ssh.NewPublicKeysFromFile(
		"git",
		setting.AppSetting.GitSshKey, "")

	_, err 		= git.PlainClone(path, false, &git.CloneOptions{
		URL:               url,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		Auth:              sshAuth,
	})

	if err != nil {
		return err
	}

	return nil
}

func GitPull(aid int,url string) error {
	w, err := ReturnGitWorkDir(aid, url)
	if err != nil {
		return err
	}

	sshAuth, err := ssh.NewPublicKeysFromFile(
		"git",
		setting.AppSetting.GitSshKey, "")

	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: sshAuth,})

	if err != nil {
		return  err
	}

	return nil
}



func GitCheckoutByBranch(aid int, url, branch string) error {
	GitPull(aid, url)

	w, err := ReturnGitWorkDir(aid, url)
	if err != nil {
		return err
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName("refs/heads/" + branch),
	})
	if err != nil {
		return err
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

func ReturnGitBranch(aid int, url string) ([]string, error) {
	e 	:= GitPull(aid, url)
	if e!= nil && e.Error() != "already up-to-date" {
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
	err := GitCheckoutByBranch(aid, url, branch)
	if err != nil {
		return nil, err
	}

	cmd 		:= exec.Command("git", "log", "--pretty=format:%h (%cr)  %ce %s", "-10")
	cmd.Dir 	=  ReturnGitLocalPath(aid, url)
	out, err 	:= cmd.CombinedOutput()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
		return nil, err
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
	fmt.Println(LongHash)

	e 	= w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(LongHash),
	})

	if e!= nil {
		return e
	}

	return nil
}