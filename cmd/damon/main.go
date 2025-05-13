package main

import (
	"log"

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
	for _, p := range l {
		log.Printf("%+v", p)
		err = s.RemoveProject(p.Id)
		if err != nil {
			log.Fatal(err)
		}
	}
}
