package models

import "time"

type DomainInfo struct {
	Model
	Name      		string
	Status			int
	IsCert			int
	CertName 		string
	DomainEndTime 	time.Time
	CertEndTime 	time.Time
	Desc  			string
}
