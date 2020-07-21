package models

type About struct {
	SystemInfo 		string
	Golangversion 	string
	GinVersion		string
}

type Settings struct {
	Model
	Name        string
	Value		string
	Desc     	string
}

type SettingRobot struct {
	Model
	Name        string
	Webhook		string
	Secret		string
	Keyword		string
	Type 		int
	Status 		int
	Desc     	string
}