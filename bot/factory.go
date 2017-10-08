package bot

import "github.com/aziule/conversation-management/platform"

type FacebookBot struct {
	*bot
}

func NewFacebookBot() *FacebookBot {
	return &FacebookBot{
		bot: &bot{
			uuid: "Abc",
			platform: platform.FACEBOOK,
			entrypoint: nil,
			webhooks: nil,
		},
	}
}