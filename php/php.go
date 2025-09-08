package php

import (
	"fmt"
	"os"
	"regexp"
)

func SetConfig(version string, key string, value string) error {
	file := fmt.Sprintf("/opt/homebrew/etc/php/%s/php.ini", version)

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(fmt.Sprintf("^;?%s.*$", key))
	newContent := re.ReplaceAllString(string(content), fmt.Sprintf("%s = %s", key, value))

	err = os.WriteFile(file, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
