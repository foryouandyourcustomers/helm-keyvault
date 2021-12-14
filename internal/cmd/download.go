package cmd

import (
	"fmt"
)

// DownloadSecret - Download and decode secret to be used as downloader plugin
func Download(uri string) error {

	su, err := newUri(uri[len("keyvault+"):])
	if err != nil {
		return err
	}

	value, err := su.download()
	if err != nil {
		return err
	}
	fmt.Print(value)
	return nil
}
