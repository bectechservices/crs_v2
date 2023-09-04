package actions

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	b64 "encoding/base64"

	"github.com/gobuffalo/buffalo"
)

func ShowPasswordManager(c buffalo.Context) error {

	c.Set("password", c.Param("password"))
	return c.Render(http.StatusOK, r.HTML("passwordmanager.html"))
}

func SavePassword(c buffalo.Context) error {
	request := AccountSetupRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	password := request.Password
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	err = ioutil.WriteFile(filepath.Join(exPath, "becencrpt.txt"), []byte(b64.StdEncoding.EncodeToString(EncryptPassword([]byte(password)))), 0644)
	if err != nil {
		log.Fatal(err)
	}

	return c.Redirect(http.StatusFound, "/")
}
