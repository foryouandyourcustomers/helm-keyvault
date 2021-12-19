package structs

import (
	"errors"
	"fmt"
	"strings"
)

func splitPath(p string) ([]string, error) {
	paths := strings.Split(p, "/") // we assume we always have a leading /!
	// path needs to contain at least type and name, version is optional
	if len(paths) < 3 {
		return nil, errors.New(fmt.Sprintf("Invalid keyvault path '%s'", p))
	}
	if len(paths) > 4 {
		return nil, errors.New(fmt.Sprintf("Invalid keyvault path '%s'", p))
	}
	return paths, nil
}
