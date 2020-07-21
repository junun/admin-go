package models

// 主机类型
type HostRole struct {
	Model
	Name      	string
	Desc  		string
}

// 主机业务关联
type HostApp struct {
	Model
	Aid   		int
	Hid     	int
	Status		int
	Desc  		string
}

// 主机信息
type Host struct {
	Model
	Rid				int
	EnvId 			int
	ZoneId			int
	Status 			int
	Enable   		int
	Name      		string
	Addres 			string
	Port			int
	Username		string
	Desc  			string
	Operator		int
}
