package models

import (
	"log"
)

const (
	// redis 角色菜单 key 前缀
	RoleMenuListKey = "role_menu_list_"

	// redis 角色权限 key 前缀
	RoleRermsSetKey = "role_perms_set_"

	// 锁定项目，一个项目同时只能允许一个执行提单或者发布
	GitAppOnWorking	= "git_app_"

	// all perms key
	AllPermsKey 	= "all_perms_key"
)

// 角色
type Role struct {
	Model
	Name      	string
	Desc  		string
}

// 用户
type User struct {
	Model
	Rid				int
	Name      		string
	Nickname 		string
	PasswordHash 	string `json:"-"`
	Email			string
	Mobile			string
	IsSupper  		int
	IsActive		int
	AccessToken 	string
	TokenExpired 	int64
}

// 菜单权限
type MenuPermissions struct {
	Model
	Pid				int
	Name      		string
	Type			int
	Permission		string
	Url				string
	Icon			string
	Desc			string
	Children    	[]*MenuPermissions 	`json:"children"`
}

//角色权限
type RolePermissionRel struct {
	Model
	Rid				int
	Pid				int
}

type History struct {
	Model
	Title     string
	Year      int
	Month     int
	Day       int
	ImageUrl  string
	Desc  	  string
}

type Banner struct {
	Model
	Bid       	int
	Title     	string
	ImageUrl  	string
	Icon		string
	Status		int
	CreateTime	int64
	UpdateTime 	int64
}

func (u *User) ReturnPermissions() []string {
	var res []string
	if u.IsSupper != 1 {
		rows, err := DB.Table("menu_permissions").
			Select("menu_permissions.permission").
			Joins("left join role_permission_rel on menu_permissions.id = role_permission_rel.pid").
			Where("role_permission_rel.rid = ?", u.Rid).
			Rows()

		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var name string
			if e := rows.Scan(&name); e != nil {
				panic(e)
			}
			res = append(res, name)
		}
	}

	return res
}

func SetRolePermToSet(key string, rid int) {
	var mps []MenuPermissions

	DB.Table("menu_permissions").
		Select("menu_permissions.permission").
		Joins("left join role_permission_rel on menu_permissions.id = role_permission_rel.pid").
		Where("role_permission_rel.rid = ?", rid).
		Find(&mps)

	for _, v := range mps {
		e := SetValBySetKey(key, v.Permission)

		if e != nil {
			log.Fatal(e)
		}
	}
}

