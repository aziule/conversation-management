package platform

type Platform string

const (
	FACEBOOK Platform = "facebook"
)

// All available platforms
func Platforms() []Platform {
	return []Platform{
		FACEBOOK,
	}
}