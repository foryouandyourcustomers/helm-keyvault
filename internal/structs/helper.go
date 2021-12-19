package structs

import (
	"errors"
	"fmt"
	"strings"
)

func splitPath(p string) ([]string, error) {
	paths := strings.Split(p, "/")
	// path needs to contain at least scheme, host, type, and name, version is optional
	if len(paths) < 3 {
		println("yarp")
		return []string{}, errors.New(fmt.Sprintf("Invalid keyvault path '%s'", p))
	}

	return paths, nil
}
