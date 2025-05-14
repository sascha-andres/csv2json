package file

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/sascha-andres/csv2json"
)

const calculatedFieldsPathName = "calculated_fields"

// CreateCalculatedFields to add a calculated fields to the project
func (s Storer) CreateCalculatedFields(projectID string, variables map[string]csv2json.CalculatedField) error {
	existing, err := s.loadCalculatedFields(projectID)
	if err != nil {
		return err
	}
	if existing == nil {
		existing = make(map[string]csv2json.CalculatedField)
	}
	for k, v := range variables {
		existing[k] = v
	}
	return s.saveCalculatedFields(projectID, existing)
}

// GetCalculatedFields to load all calculated fields
func (s Storer) GetCalculatedFields(projectID string) (map[string]csv2json.CalculatedField, error) {
	return s.loadCalculatedFields(projectID)
}

// RemoveCalculatedFields to remove one or more calculated fields
func (s Storer) RemoveCalculatedFields(projectID string, calculatedFields []string) error {
	existing, err := s.loadCalculatedFields(projectID)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil
	}
	for _, k := range calculatedFields {
		delete(existing, k)
	}
	return s.saveCalculatedFields(projectID, existing)
}

// ClearCalculatedFields to clear the project from extra variables
func (s Storer) ClearCalculatedFields(projectID string) error {
	return s.bucket.Delete(context.Background(), getCalculatedFieldsPathForProject(projectID))
}

// loadCalculatedFields deserializes mappings from file storage
func (s Storer) loadCalculatedFields(projectID string) (map[string]csv2json.CalculatedField, error) {
	path := getCalculatedFieldsPathForProject(projectID)
	ctx := context.Background()
	known, err := s.bucket.Exists(ctx, path)
	if err != nil {
		return nil, err
	}
	if !known {
		return nil, nil
	}
	reader, err := s.bucket.NewReader(ctx, path, nil)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = reader.WriteTo(&buf)
	if err != nil {
		return nil, err
	}
	var calculatedFields map[string]csv2json.CalculatedField
	err = json.Unmarshal(buf.Bytes(), &calculatedFields)
	if err != nil {
		return nil, err
	}
	return calculatedFields, nil
}

// saveCalculatedFields persists to storage
func (s Storer) saveCalculatedFields(id string, calculatedFields map[string]csv2json.CalculatedField) error {
	path := getCalculatedFieldsPathForProject(id)
	ctx := context.Background()
	w, err := s.bucket.NewWriter(ctx, path, nil)
	if err != nil {
		return err
	}
	data, err := json.Marshal(calculatedFields)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return w.Close()
}

// getCalculatedFieldsPathForProject returns the file storage path for mappings
func getCalculatedFieldsPathForProject(projectId string) string {
	return fmt.Sprintf("admin/projects/%s/%s.json", projectId, calculatedFieldsPathName)
}
