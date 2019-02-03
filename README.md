# flow oAuth
![Travis CI build](https://api.travis-ci.org/go-flow/oauth.svg?branch=master)
[![GoDoc](https://godoc.org/github.com/go-flow/oauth?status.svg)](https://godoc.org/github.com/go-flow/oauth)
[![GoReport](https://goreportcard.com/badge/github.com/go-flow/oauth)](https://goreportcard.com/report/github.com/go-flow/oauth)
![GitHub contributors](https://img.shields.io/github/contributors/go-flow/oauth.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Flow oAuth is created to manage social oAuth authentication without problems.
the package is inspired by Gocialite 

## Installation

To install it, just run `go get github.com/go-flow/oauth` and include it in your app: `import "github.com/go-flow/oauth"`.

## Available drivers

- Amazon
- Asana
- Bitbucket
- Facebook
- Foursquare
- Github
- Google
- LinkedIn
- Slack

## Create new driver

Please see [Contributing page](https://github.com/go-flow/oauth/blob/master/CONTRIBUTING.md) to learn how to create new driver and test it.

## Set scopes

**Note**: Flow oAuth set some default scopes for the user profile, for example for *Facebook* it specify `email` and for *Google* `profile, email`.  
When you use the following method, you don't have to rewrite them. 

Use the `Scopes([]string)` method of your `oauth` instance. Example:

```go
oauth.Scopes([]string{"public_repo"})
```

## Set driver

Use the `Driver(string)` method of your `oauth` instance. Example:

```go
oauth.Driver("facebook")
```

The driver name will be the provider name in lowercase.

## How to use it

**Note**: All oAuth methods are chainable.

Declare a "global" variable outside your `main` func:

```go
import (
	...
)

var oAuth = oauth.NewDispatcher()

func main() {
```

Then create a route to use as redirect bridge, for example `/auth/github`. With this route, the user will be redirected to the provider oAuth login. In this case we use Flow framework. You have to specify the provider with the `Driver()` method.
Then, with `Scopes()`, you can set a list of scopes as slice of strings. It's optional.  
Finally, with `Redirect()` you can obtain the redirect URL. In this method you have to pass three parameters:

1. Client ID
1. Client Secret
1. Redirect URL

```go

func main() {
	
    app := flow.New()

	app.GET("/auth/github", redirectHandler)

	if err := app.Serve(); err != nil && err != http.ErrServerClosed {
		app.Logger.Error(err.Error())
	}
}

// Redirect to correct oAuth URL
func redirectHandler(ctx *flow.Context) {
	authURL, err := oAuth.New().
		Driver("github"). // Set provider
		Scopes([]string{"public_repo"}). // Set optional scope(s)
		Redirect( // 
			"xxxxxxxxxxxxxx", // Client ID
			"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", // Client Secret
			"http://localhost:3000/auth/github/callback", // Redirect URL
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		ctx.ServeError(400, err)
		return
	}

	// Redirect with authURL
	ctx.Redirect(http.StatusFound, authURL) // Redirect with 302 HTTP code
}
```

Now create a callback handler route, where we'll receive the content from the provider.  
In order to validate the oAuth and retrieve the data, you have to invoke the `Handle()` method with two query parameters: `state` and `code`. In your URL, they will look like this: `http://localhost:3000/auth/github/callback?state=xxxxxxxx&code=xxxxxxxx`.  
The `Handle()` method returns the user info, the token and error if there's one or `nil`.  
If there are no errors, in the `user` variable you will find the logged in user information and in the `token` one, the token info (it's a [oauth2.Token struct](https://godoc.org/golang.org/x/oauth2#Token)). The data of the user - which is a [models.User struct](https://github.com/go-flow/oauth/blob/master/structs/user.go) - are the following:

- ID
- FirstName
- LastName
- FullName
- Email
- Avatar (URL)
- Raw (the full JSON returned by the provider)

Note that they can be empty.

```go
func main() {
	app := flow.New()

	app.GET("/auth/github", redirectHandler)
	app.GET("/auth/github/callback", callbackHandler)

	if err := app.Serve(); err != nil && err != http.ErrServerClosed {
		app.Logger.Error(err.Error())
	}
}

// Redirect to correct oAuth URL
// Handle callback of provider
func callbackHandler(ctx *flow.Context) {
	// Retrieve query params for code and state
	code := ctx.Query("code")
	state := ctx.Query("state")

	// Handle callback and check for errors
	user, token, err := oAuth.Handle(state, code)
	if err != nil {
		ctx.ServeError(400, err)
		return
	}

	// Print in terminal user information
	fmt.Printf("%#v", token)
	fmt.Printf("%#v", user)

	// If no errors, show provider name
	ctx.Response.Write([]byte("Hi, " + user.FullName))
}
```