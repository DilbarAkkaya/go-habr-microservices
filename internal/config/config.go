package config

import "os"

func GetDBConnStr()string {
	if val:=os.Getenv("DATABASE_URL"); val!=""{
		return val
	}
	return "host=localhost port=5432 user=user password=password dbname=habr_db sslmode=disable"
}