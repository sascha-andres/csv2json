package types

import (
	"github.com/sascha-andres/csv2json"
	"github.com/sascha-andres/csv2json/admin"
	"github.com/sascha-andres/csv2json/storer"
)

type (
	// Storer defines all methods to access persistence
	Storer interface {
		// CreateProject is used to create a project
		CreateProject(p storer.Project) error

		// RemoveProject removes project data (incl all run data)
		RemoveProject(id string) error

		// ListProjects returns all known projects
		ListProjects() ([]storer.Project, error)

		// CreateMappings allows adding or updating mappings for a project
		CreateMappings(projectID string, columns map[string]csv2json.ColumnConfiguration) (map[string]admin.ActionTaken, error)

		// RemoveMappings removes all mappings provided from storage
		RemoveMappings(projectID string, columns []string) error

		// GetMappings reads all mappings and returns those
		GetMappings(projectID string) (map[string]csv2json.ColumnConfiguration, error)

		// ClearMappings is called when project ist removed
		ClearMappings(projectID string) error

		// CreateExtraVariables to add a static variable to the project
		CreateExtraVariables(projectID string, variables map[string]string) error

		// GetExtraVariables to load all extra variables
		GetExtraVariables(projectID string) (map[string]string, error)

		// RemoveExtraVariables to remove one or more extra variables
		RemoveExtraVariables(projectID string, extraVariables []string) error

		// ClearExtraVariables to clear the project from extra variables
		ClearExtraVariables(projectID string) error

		// CreateCalculatedFields to add a calculated fields to the project
		CreateCalculatedFields(projectID string, variables map[string]csv2json.CalculatedField) error

		// GetCalculatedFields to load all calculated fields
		GetCalculatedFields(projectID string) (map[string]csv2json.CalculatedField, error)

		// RemoveCalculatedFields to remove one or more calculated fields
		RemoveCalculatedFields(projectID string, calculatedFields []string) error

		// ClearCalculatedFields to clear the project from extra variables
		ClearCalculatedFields(projectID string) error
	}
)
