package actions

import (
	b64 "encoding/base64"
	"fmt"
	"log"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/jinzhu/gorm"
	"github.com/robfig/cron"

	//import the mssql driver
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

//DatabaseConnection holds the global database connection
var DatabaseConnection *gorm.DB

//OpenDatabaseConnection opens a database connection
func OpenDatabaseConnection() {
	var err error
	port := envy.Get("CRS_DB_PORT", "10501")
	encryptedPassword, err := b64.StdEncoding.DecodeString(string(ReadPasswordFromFile()))
	if err != nil {
		log.Println(err)
	}
	password := DecryptPassword(encryptedPassword)
	DatabaseConnection, err = gorm.Open("mssql", fmt.Sprintf("sqlserver://crsadmin:%s@localhost:%s?database=CRS_DB", password, port))
	if err != nil {
		panic(err)
	}
}

//AuthID returns the id of the authenticated user
func AuthID(c buffalo.Context) int {
	return AuthUser(c).ID
}

//AuthUser returns the authenticated user
func AuthUser(c buffalo.Context) User {
	return c.Value("auth_user").(User)
}

func StartCronScheduler() *cron.Cron {
	c := cron.New()
	_ = c.AddFunc("0 0 10 20 * *", ScheduleEmailsToBeSent)
	c.Start()
	return c
}
