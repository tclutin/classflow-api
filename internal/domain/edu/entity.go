package edu

type Faculty struct {
	FacultyID uint64
	Name      string
}

type Program struct {
	ProgramID uint64
	FacultyID uint64
	Name      string
}
