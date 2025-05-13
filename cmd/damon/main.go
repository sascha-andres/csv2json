package main

import (
	"log"

	"github.com/sascha-andres/csv2json"
	"github.com/sascha-andres/csv2json/internal/persistence"
	"github.com/sascha-andres/csv2json/storer"
)

func main() {
	s, err := persistence.GetStorer("file:file:///Users/andres/tmp/csv2json")
	if err != nil {
		log.Fatal(err)
	}
	err = s.CreateProject(storer.Project{Id: "123", Name: "test"})
	if err != nil {
		log.Fatal(err)
	}
	l, err := s.ListProjects()
	if err != nil {
		log.Fatal(err)
	}
	mapping := make(map[string]csv2json.ColumnConfiguration)
	mapping["name"] = csv2json.ColumnConfiguration{Property: "property", Type: "string"}
	mapping["age"] = csv2json.ColumnConfiguration{Property: "property", Type: "int"}
	a, err := s.CreateMappings(l[0].Id, mapping)
	if err != nil {
		log.Fatal(err)
	}
	for i := range a {
		log.Printf("%+v", a[i])
	}
	err = s.RemoveMappings(l[0].Id, []string{"name"})
	if err != nil {
		log.Fatal(err)
	}
	a, err = s.CreateMappings(l[0].Id, mapping)
	if err != nil {
		log.Fatal(err)
	}
	for i := range a {
		log.Printf("%+v", a[i])
	}
	readMappings, err := s.GetMappings(l[0].Id)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", readMappings)
	err = s.CreateExtraVariables(l[0].Id, map[string]string{"name": "Andres"})
	if err != nil {
		log.Fatal(err)
	}
	readExtraVariables, err := s.GetExtraVariables(l[0].Id)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", readExtraVariables)
	for _, p := range l {
		log.Printf("%+v", p)
		err = s.RemoveProject(p.Id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
