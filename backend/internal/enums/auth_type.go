// internal/enums/auth_type.go
package enums

type AuthType string

const (
	AuthTypeTelegram AuthType = "telegram"
	AuthTypeDiscord  AuthType = "discord"
)

func (AuthType) Values() []AuthType {
	return []AuthType{
		AuthTypeTelegram,
		AuthTypeDiscord,
	}
}
