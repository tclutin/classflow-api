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

type TypeOfSubject struct {
	TypeOfSubjectID uint64
	Name            string
}

type Building struct {
	BuildingID uint64
	Name       string
	Latitude   float64
	Longitude  float64
	Address    string
}
