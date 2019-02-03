package oauth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var oAuthTest OAuth

func TestScopes(t *testing.T) {
	oAuthTest.Scopes([]string{"email"})
	assert.Equal(t, oAuthTest.scopes, []string{"email"})
	assert.NotEqual(t, oAuthTest.scopes, []string{})

	oAuthTest.
		Driver("google").
		Scopes([]string{"calendar.readonly"})
	assert.Equal(t, oAuthTest.scopes, []string{"profile", "email", "calendar.readonly"})
	assert.NotEqual(t, oAuthTest.scopes, []string{"profile", "email"})
	assert.NotEqual(t, oAuthTest.scopes, []string{})
}
func TestConf(t *testing.T) {
	assert := assert.New(t)

	oAuthTest.
		Driver("github").
		Redirect(
			"foo",
			"bar",
			"http://example.com/auth/callback",
		)

	assert.Equal(oAuthTest.conf.ClientID, "foo")
	assert.NotEqual(oAuthTest.conf.ClientID, "")
	assert.NotNil(oAuthTest.conf.ClientID)

	assert.Equal(oAuthTest.conf.ClientSecret, "bar")
	assert.NotEqual(oAuthTest.conf.ClientSecret, "")
	assert.NotNil(oAuthTest.conf.ClientSecret)

	assert.Equal(oAuthTest.conf.RedirectURL, "http://example.com/auth/callback")
	assert.NotEqual(oAuthTest.conf.RedirectURL, "")
	assert.NotNil(oAuthTest.conf.RedirectURL)
}
func TestDriver(t *testing.T) {
	var err error

	_, err = oAuthTest.Driver("unknown").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"http://example.com/auth/callback",
		)
	assert.NotNil(t, err)

	_, err = oAuthTest.Driver("github").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"http://example.com/auth/callback",
		)
	assert.Nil(t, err)
}
func TestRedirectURL(t *testing.T) {
	var err error

	_, err = oAuthTest.Driver("github").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"/auth/callback",
		)
	assert.NotNil(t, err)

	_, err = oAuthTest.Driver("github").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"http://example.com/auth/callback",
		)
	assert.Nil(t, err)
}

func TestState(t *testing.T) {
	var err error

	err = oAuthTest.Driver("github").
		Handle("fakeState", "foo")
	assert.NotNil(t, err)
}

func TestExchange(t *testing.T) {
	var err error

	// Generate a state
	oAuthTest.
		Driver("github").
		Redirect(
			"xxxxxxxx",
			"xxxxxxxxxxxxxxxxxxxxxxxx",
			"http://example.com/auth/callback",
		)

	err = oAuthTest.Handle(oAuthTest.state, "foo")
	assert.NotNil(t, err)
}
