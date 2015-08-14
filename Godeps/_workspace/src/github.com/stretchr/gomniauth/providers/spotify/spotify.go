package spotify

import (
	"net/http"

	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/gomniauth"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/gomniauth/common"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/gomniauth/oauth2"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/objx"
)

const (
	spotifyDefaultScope    string = "user-read-email"
	spotifyName            string = "spotify"
	spotifyDisplayName     string = "Spotify"
	spotifyAuthURL         string = "https://accounts.spotify.com/authorize"
	spotifyTokenURL        string = "https://accounts.spotify.com/api/token"
	spotifyEndpointProfile string = "https://api.spotify.com/v1/me"
)

// SpotifyProvider implements the Provider interface and provides Spotify
// OAuth2 communication capabilities.
type SpotifyProvider struct {
	config         *common.Config
	tripperFactory common.TripperFactory
}

func New(clientId, clientSecret, redirectUrl string) *SpotifyProvider {

	p := new(SpotifyProvider)
	p.config = &common.Config{Map: objx.MSI(
		oauth2.OAuth2KeyAuthURL, spotifyAuthURL,
		oauth2.OAuth2KeyTokenURL, spotifyTokenURL,
		oauth2.OAuth2KeyClientID, clientId,
		oauth2.OAuth2KeySecret, clientSecret,
		oauth2.OAuth2KeyRedirectUrl, redirectUrl,
		oauth2.OAuth2KeyScope, spotifyDefaultScope,
		oauth2.OAuth2KeyAccessType, oauth2.OAuth2AccessTypeOnline,
		oauth2.OAuth2KeyApprovalPrompt, oauth2.OAuth2ApprovalPromptAuto,
		oauth2.OAuth2KeyResponseType, oauth2.OAuth2KeyCode)}
	return p
}

// TripperFactory gets an OAuth2TripperFactory
func (provider *SpotifyProvider) TripperFactory() common.TripperFactory {

	if provider.tripperFactory == nil {
		provider.tripperFactory = new(oauth2.OAuth2TripperFactory)
	}

	return provider.tripperFactory
}

// PublicData gets a public readable view of this provider.
func (provider *SpotifyProvider) PublicData(options map[string]interface{}) (interface{}, error) {
	return gomniauth.ProviderPublicData(provider, options)
}

// Name is the unique name for this provider.
func (provider *SpotifyProvider) Name() string {
	return spotifyName
}

// DisplayName is the human readable name for this provider.
func (provider *SpotifyProvider) DisplayName() string {
	return spotifyDisplayName
}

// GetBeginAuthURL gets the URL that the client must visit in order
// to begin the authentication process.
//
// The state argument contains anything you wish to have sent back to your
// callback endpoint.
// The options argument takes any options used to configure the auth request
// sent to the provider. In the case of OAuth2, the options map can contain:
//   1. A "scope" key providing the desired scope(s). It will be merged with the default scope.
func (provider *SpotifyProvider) GetBeginAuthURL(state *common.State, options objx.Map) (string, error) {
	if options != nil {
		scope := oauth2.MergeScopes(options.Get(oauth2.OAuth2KeyScope).Str(), spotifyDefaultScope)
		provider.config.Set(oauth2.OAuth2KeyScope, scope)
	}
	return oauth2.GetBeginAuthURLWithBase(provider.config.Get(oauth2.OAuth2KeyAuthURL).Str(), state, provider.config)
}

// Get makes an authenticated request and returns the data in the
// response as a data map.
func (provider *SpotifyProvider) Get(creds *common.Credentials, endpoint string) (objx.Map, error) {
	return oauth2.Get(provider, creds, endpoint)
}

// GetUser uses the specified common.Credentials to access the users profile
// from the remote provider, and builds the appropriate User object.
func (provider *SpotifyProvider) GetUser(creds *common.Credentials) (common.User, error) {

	profileData, err := provider.Get(creds, spotifyEndpointProfile)

	if err != nil {
		return nil, err
	}

	// build user
	user := NewUser(profileData, creds, provider)

	return user, nil
}

// CompleteAuth takes a map of arguments that are used to
// complete the authorisation process, completes it, and returns
// the appropriate Credentials.
func (provider *SpotifyProvider) CompleteAuth(data objx.Map) (*common.Credentials, error) {
	return oauth2.CompleteAuth(provider.TripperFactory(), data, provider.config, provider)
}

// GetClient returns an authenticated http.Client that can be used to make requests to
// protected spotify resources
func (provider *SpotifyProvider) GetClient(creds *common.Credentials) (*http.Client, error) {
	return oauth2.GetClient(provider.TripperFactory(), creds, provider)
}
