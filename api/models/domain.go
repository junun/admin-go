package models

import "time"

type DomainInfo struct {
	Model
	Name      		string
	Status			int
	StartTime       time.Time
	EndTime 		time.Time
	Channel			string
	Desc  			string
}


type CertificateInfo struct {
	Model
	Did				int
	Name      		string
	Status			int
	StartTime       time.Time
	EndTime 		time.Time
	Channel			string
	Desc  			string
}

