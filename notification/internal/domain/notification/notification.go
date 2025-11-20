package notification

type Channel int

const (
	ChannelEmail Channel = iota
	ChanelPush
	ChannelCall

	Reject
)

type UserContact struct {
	Channel Channel
	Value   string
	Meta    map[string]string
}

type UserContacts struct {
	UserID   string
	Contacts []*UserContact
}

type Input struct {
	UserContacts *UserContacts
}

type Output struct {
}
