package facebook

import "time"

type message struct {
	SenderId string
	RecipientId string
	Timestamp time.Time
	Mid string
}

type textMessage struct {
	message
	Text string
}

type quickReplyMessage struct {
	textMessage
	QuickReplyPayload string
}
