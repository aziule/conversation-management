package facebook

// Our FacebookUser, implementing the conversation.User interface
type FacebookUser struct {
	uuid string
	fbid string
	name string
}

// Getters
func (facebookUser *FacebookUser) Uuid() string { return facebookUser.uuid }
func (facebookUser *FacebookUser) Fbid() string { return facebookUser.fbid }
func (facebookUser *FacebookUser) Name() string { return facebookUser.name }
