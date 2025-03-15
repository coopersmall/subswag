package chatcompletion

import "errors"

type Message struct {
	Role    Role
	Content string
}

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

func Roles() []Role {
	return []Role{
		RoleUser,
		RoleAssistant,
		RoleSystem,
	}
}

func (r Role) String() string {
	return string(r)
}

func (r Role) Validate() error {
	roles := Roles()
	for _, role := range roles {
		if role == r {
			return nil
		}
	}
	return errors.New("invalid role")
}
