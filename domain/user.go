package domain

type RoleType string

const (
	Role_Student RoleType = "student"
	Role_Teacher RoleType = "teacher"
)

func (r RoleType) IsValid() bool {
	return r == Role_Teacher || r == Role_Student
}

type User struct {
	ID       uint
	Username string
	Email    string
	Password string
	Role     RoleType
}
