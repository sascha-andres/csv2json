package file

import (
	"context"

	"gocloud.dev/blob"

	"github.com/sascha-andres/csv2json/internal/types"
)

// Create returns a Storer instance
func Create(dsn string) (types.Storer, error) {
	s := Storer{dsn: dsn}
	bucket, err := blob.OpenBucket(context.Background(), s.dsn)
	if err != nil {
		return nil, err
	}
	s.bucket = bucket
	return s, nil
}
