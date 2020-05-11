package models

import "time"

type AppDeploy struct {
	Model
	Tid				int
	Name      		string
	RepoBranch 		string
	RepoCommit 		string
	Status			int
	Operator		int
	Review          int
	Deploy          int
	UpdateTime 		time.Time
}
