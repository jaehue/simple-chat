package oauth2

import (
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/gomniauth/common"
)

// AuthorizationHeader returns the key, value pair to insert into an authorized request.
func AuthorizationHeader(creds *common.Credentials) (key, value string) {
	return "Authorization", "Bearer " + creds.Get(OAuth2KeyAccessToken).Str("Invalid")
}
