package conversation

// Our main User struct
type User struct {
	uuid string
	fbid string
	name string
}

// Getters
func (user *User) Uuid() string { return user.uuid }
func (user *User) Fbid() string { return user.fbid }
func (user *User) Name() string { return user.name }
