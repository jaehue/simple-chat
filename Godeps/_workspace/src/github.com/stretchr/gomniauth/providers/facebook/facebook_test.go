package facebook

import (
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/gomniauth/common"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/gomniauth/oauth2"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/gomniauth/test"
	"github.com/jaehue/simple-chat/Godeps/_workspace/src/github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestGitHubImplementrsProvider(t *testing.T) {

	var provider common.Provider
	provider = new(FacebookProvider)

	assert.NotNil(t, provider)

}

func TestGetUser(t *testing.T) {

	g := New("clientID", "secret", "http://myapp.com/")
	creds := &common.Credentials{Map: objx.MSI()}

	testTripperFactory := new(test.TestTripperFactory)
	testTripper := new(test.TestTripper)
	testTripperFactory.On("NewTripper", mock.Anything, g).Return(testTripper, nil)
	testResponse := new(http.Response)
	testResponse.Header = make(http.Header)
	testResponse.Header.Set("Content-Type", "application/json")
	testResponse.StatusCode = 200
	testResponse.Body = ioutil.NopCloser(strings.NewReader(`{
  "id": "631819186",
  "name": "their-name",
  "first_name": "Mat",
  "last_name": "Ryer",
  "link": "https://www.facebook.com/matryer",
  "username": "loginname",
  "bio": "http://www.stretchr.com/",
  "gender": "male",
  "email": "email@address.com",
  "timezone": -6,
  "locale": "en_GB",
  "verified": true,
  "updated_time": "2013-10-03T19:55:28+0000",
  "picture": {
    "url": "http://www.myface.com/"
  }
  }`))
	testTripper.On("RoundTrip", mock.Anything).Return(testResponse, nil)

	g.tripperFactory = testTripperFactory

	user, err := g.GetUser(creds)

	if assert.NoError(t, err) && assert.NotNil(t, user) {

		assert.Equal(t, user.Name(), "their-name")
		assert.Equal(t, user.AuthCode(), "") // doesn't come from github
		assert.Equal(t, user.Nickname(), "loginname")
		assert.Equal(t, user.Email(), "email@address.com")
		assert.Equal(t, user.AvatarURL(), "http://www.myface.com/")
		assert.Equal(t, user.Data()["link"], "https://www.facebook.com/matryer")

		facebookCreds := user.ProviderCredentials()[facebookName]
		if assert.NotNil(t, facebookCreds) {
			assert.Equal(t, "631819186", facebookCreds.Get(common.CredentialsKeyID).Str())
		}

	}

}

func TestNewfacebook(t *testing.T) {

	g := New("clientID", "secret", "http://myapp.com/")

	if assert.NotNil(t, g) {

		// check config
		if assert.NotNil(t, g.config) {

			assert.Equal(t, "clientID", g.config.Get(oauth2.OAuth2KeyClientID).Data())
			assert.Equal(t, "secret", g.config.Get(oauth2.OAuth2KeySecret).Data())
			assert.Equal(t, "http://myapp.com/", g.config.Get(oauth2.OAuth2KeyRedirectUrl).Data())
			assert.Equal(t, facebookDefaultScope, g.config.Get(oauth2.OAuth2KeyScope).Data())

			assert.Equal(t, facebookAuthURL, g.config.Get(oauth2.OAuth2KeyAuthURL).Data())
			assert.Equal(t, facebookTokenURL, g.config.Get(oauth2.OAuth2KeyTokenURL).Data())

		}

	}

}

func TestfacebookTripperFactory(t *testing.T) {

	g := New("clientID", "secret", "http://myapp.com/")
	g.tripperFactory = nil

	f := g.TripperFactory()

	if assert.NotNil(t, f) {
		assert.Equal(t, f, g.tripperFactory)
	}

}

func TestfacebookName(t *testing.T) {
	g := New("clientID", "secret", "http://myapp.com/")
	assert.Equal(t, facebookName, g.Name())
}

func TestfacebookGetBeginAuthURL(t *testing.T) {

	common.SetSecurityKey("ABC123")

	state := &common.State{Map: objx.MSI("after", "http://www.stretchr.com/")}

	g := New("clientID", "secret", "http://myapp.com/")

	url, err := g.GetBeginAuthURL(state, nil)

	if assert.NoError(t, err) {
		assert.Contains(t, url, "client_id=clientID")
		assert.Contains(t, url, "redirect_uri=http%3A%2F%2Fmyapp.com%2F")
		assert.Contains(t, url, "scope="+facebookDefaultScope)
		assert.Contains(t, url, "access_type="+oauth2.OAuth2AccessTypeOnline)
		assert.Contains(t, url, "approval_prompt="+oauth2.OAuth2ApprovalPromptAuto)
	}

	state = &common.State{Map: objx.MSI("after", "http://www.stretchr.com/")}

	g = New("clientID", "secret", "http://myapp.com/")

	url, err = g.GetBeginAuthURL(state, objx.MSI(oauth2.OAuth2KeyScope, "avatar"))

	if assert.NoError(t, err) {
		assert.Contains(t, url, "client_id=clientID")
		assert.Contains(t, url, "redirect_uri=http%3A%2F%2Fmyapp.com%2F")
		assert.Contains(t, url, "scope=avatar+"+facebookDefaultScope)
		assert.Contains(t, url, "access_type="+oauth2.OAuth2AccessTypeOnline)
		assert.Contains(t, url, "approval_prompt="+oauth2.OAuth2ApprovalPromptAuto)
	}

}
