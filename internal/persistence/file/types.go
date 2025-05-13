package file

import (
	"gocloud.dev/blob"
)

type (
	Storer struct {
		// dsn to storage location
		dsn string

		// bucket is instance of gocloud.dev/blob stuff
		bucket *blob.Bucket
	}
)
