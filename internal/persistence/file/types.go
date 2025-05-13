package file

import (
	"context"
	"encoding/json"
	"fmt"

	"gocloud.dev/blob"

	"github.com/sascha-andres/csv2json/storer"
)

type (
	Storer struct {
		// dsn to storage location
		dsn string

		// bucket is instance of gocloud.dev/blob stuff
		bucket *blob.Bucket
	}
)

func (s Storer) CreateProject(p storer.Project) error {
	w, err := s.bucket.NewWriter(context.Background(), fmt.Sprintf("projects/%s", p.Id), nil)
	if err != nil {
		return err
	}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return w.Close()
}
