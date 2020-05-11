package models

// 环境分类表
type ConfigEnv struct {
	Model
	Name      	string
	Desc  		string
}

type AppType struct {
	Model
	Name      	string
	Desc  		string
}

type App struct {
	Model
	Tid				int
	EnvId			int
	Name      		string
	Active			int
	DeployType		int
	EnableSync		int
	Desc  			string
}

type DeployExtend struct {
	Dtid            int
	Aid				int
	Tag				string
	TemplateName    string
	EnableCheck		int
	HostIds         string
	RepoUrl         string
	Versions        int
	PreCode  		string
	PostCode  		string
	PreDeploy       string
	PostDeploy      string
	DstDir          string
	DstRepo			string
}

type AppSyncValue struct {
	Model
	Aid				int
	Name      		string
	Value      		string
	Desc  			string
}
