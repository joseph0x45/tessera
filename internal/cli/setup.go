package cli

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/joseph0x45/goutils"
	"golang.org/x/crypto/bcrypt"
)

func ensureConfig(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaultHash, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		contents := fmt.Sprintf(`PORT=8080
DASHBOARD_USER=admin
DASHBOARD_PASSWORD_HASH=%s
`,
			defaultHash,
		)
		return os.WriteFile(path, []byte(contents), 0600)
	}
	return nil
}

func setupEnv() int {
	err := ensureConfig(goutils.GetAppConfigFilePath())
	if err != nil {
		log.Println(err)
	}
	return 0
}

func SetAdminPassword(args []string) int {
	flagSet := flag.NewFlagSet("setAdminPassword", flag.ContinueOnError)
	adminPassword := flagSet.String("admin-password", "", "Value of the admin password")
	flagSet.Parse(args)
	if *adminPassword == "" {
		log.Println("admin-password is required")
		return 1
	}
	file, err := os.OpenFile(goutils.GetAppConfigFilePath(), os.O_RDWR, 0600)
	if err != nil {
		log.Println("Error while opening config file:", err.Error())
		return 1
	}
	defer file.Close()
	var conf []byte
	if conf, err = io.ReadAll(file); err != nil {
		log.Println("Error while reading config file:", err.Error())
		return 1
	}

	parsed := goutils.ParseSimpleDotenv(string(conf))
	hash, _ := goutils.HashPassword(*adminPassword)
	parsed.SetKey("DASHBOARD_PASSWORD_HASH", hash)
	updatedConf := parsed.Write()
	if _, err := file.WriteAt([]byte(updatedConf), 0); err != nil {
		log.Println("Error while updating config file:", err.Error())
		return 1
	}
	log.Println("Admin password updated successfully!")
	return 0
}
