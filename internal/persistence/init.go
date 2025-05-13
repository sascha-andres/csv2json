package persistence

import (
	"fmt"
	"strings"

	"github.com/sascha-andres/csv2json/internal/persistence/file"
	"github.com/sascha-andres/csv2json/internal/types"
)

// CreateFunc describes how a Storer can be instantiated
type CreateFunc func(string) (types.Storer, error)

// repository is used to track Storer implementations
var repository = make(map[string]CreateFunc)

// init used to fill repository
func init() {
	repository["file"] = file.Create
}

// GetStorer looks up the storer based on dsn and returns an instance
func GetStorer(dsn string) (types.Storer, error) {
	splittedDsn := strings.SplitN(dsn, ":", 2)
	initFunc, ok := repository[splittedDsn[0]]
	if !ok {
		return nil, fmt.Errorf("unknown storer %s", splittedDsn[0])
	}
	return initFunc(splittedDsn[1])
}
