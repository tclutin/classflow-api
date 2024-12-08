package errors

import "errors"

var (
	// ErrUserNotFound UserService
	ErrUserNotFound = errors.New("user not found")

	// ErrUserAlreadyExists UserService
	ErrUserAlreadyExists = errors.New("user already exists")

	// ErrWrongPassword AuthService
	ErrWrongPassword = errors.New("wrong password")

	// ErrProgramNotFound EduService
	ErrProgramNotFound = errors.New("program not found")

	// ErrFacultyNotFound EduService
	ErrFacultyNotFound = errors.New("faculty not found")

	// ErrTypeOfSubjectNotFound EduService
	ErrTypeOfSubjectNotFound = errors.New("type of subject not found")

	// ErrBuildingNotFound EduService
	ErrBuildingNotFound = errors.New("building not found")

	// ErrGroupNotFound GroupService
	ErrGroupNotFound = errors.New("group not found")

	// ErrGroupAlreadyExists GroupService
	ErrGroupAlreadyExists = errors.New("group already exists with this shortname")

	// ErrFacultyProgramIdMismatch GroupService
	ErrFacultyProgramIdMismatch = errors.New("faculty and program id does not match")

	// ErrAlreadyInGroup GroupService
	ErrAlreadyInGroup = errors.New("you are already in a group")

	// ErrGroupAlreadyHasSchedule GroupService
	ErrGroupAlreadyHasSchedule = errors.New("group already has schedule")

	//ErrMemberNotFound GroupService
	ErrMemberNotFound = errors.New("member not found")
)
