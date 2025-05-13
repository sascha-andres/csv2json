package types

import "github.com/sascha-andres/csv2json/storer"

type (
	// Storer defines all methods to access persistence
	Storer interface {
		// CreateProject is used to create a project
		CreateProject(p storer.Project) error
	}
)
