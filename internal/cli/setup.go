package cli

import (
	"fmt"
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
