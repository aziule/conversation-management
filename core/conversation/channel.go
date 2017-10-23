package conversation

type ChannelType string

const (
	OneToOne ChannelType = "onetoone"
	OneToMany ChannelType = "onetomany"
)

type Channel interface {
	Users() []User
	Type() ChannelType
}
