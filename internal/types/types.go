package types

import "github.com/sascha-andres/csv2json/storer"

type (
	// Storer defines all methods to access persistence
	Storer interface {
		// CreateProject is used to create a project
		CreateProject(p storer.Project) error

		// RemoveProject removes project data (incl all run data)
		RemoveProject(id string) error

		// ListProjects returns all known projects
		ListProjects() ([]storer.Project, error)
	}
)
