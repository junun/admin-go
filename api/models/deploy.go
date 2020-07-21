package models

import "time"

type AppDeploy struct {
	Model
	Tid				int
	GitType			string
	Name      		string
	TagBranch 		string
	Commit 		    string
	IsPass			int
	Version			string
	Reason        	string
	Desc 			string
	Status			int
	Operator		int
	Review          int
	Deploy          int
	UpdateTime 		time.Time
}
