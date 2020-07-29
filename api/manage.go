package main

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"api/models"
	"api/pkg/util"
)


var (
	h bool
	c string
	Passwd string
	Check string
)

func init() {
	flag.BoolVar(&h, "h", false, "this help")

	flag.StringVar(&c, "c", "", "create_admin : 创建管理员账户, enable_admin : 启用管理员账户")

	// 改变默认的 Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: progarm [-h] [-c do some work] 

Options:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	if len(os.Args) == 1 {
		flag.Usage()
		return
	}

	switch  {
		case c == "create_admin":
			CreateAdmin()

		case c == "enable_admin":
			EnableAdmin()

		default:
			flag.Usage()
	}
}

func CreateAdmin() {
	var user models.User

	//检查 admin 用户是否存在
	err := models.DB.Where("name = ?", "admin").First(&user).Error
	if err == gorm.ErrRecordNotFound {
		fmt.Printf("Please enter password for admin : ")
		fmt.Scanln(&Passwd)
		// 新增用户
		user.Name 				= "admin"
		user.PasswordHash, _ 	= util.HashPassword(Passwd)
		user.Nickname 			= "admin"
		user.IsSupper			= 1
		user.IsActive 			= 1
		user.TwoFactor			= 0

		if e := models.DB.Create(&user).Error; e != nil {
			panic(e)
		}
	}

	if user.Name == "admin" {
		fmt.Printf("已存在管理员账户admin，需要重置密码[y|n]？ : ")
		fmt.Scanln(&Check)

		if Check == "y" {
			fmt.Printf("Please enter password for admin : ")
			fmt.Scanln(&Passwd)
			passwdhash, _ := util.HashPassword(Passwd)
			e := models.DB.Model(&user).Update("password_hash", passwdhash).Error
			if e != nil {
				panic(e)
			}
		}
	}
}

func EnableAdmin() {
	var user models.User
	models.DB.Where("name = ?", "admin").First(&user)

	user.IsActive = 1
	e := models.DB.Save(&user).Error

	if e != nil {
		panic(e)
	}
}
