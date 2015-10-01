package main

import (
    "os"
    "fmt"
    "log"
    "strings"
    "io/ioutil"
    "net/http"
    "encoding/json"
	"github.com/go-martini/martini"
	gooauth2 "github.com/golang/oauth2"
	"github.com/martini-contrib/oauth2"
	"github.com/martini-contrib/sessions"
)

type User struct {
    User_id     string
    User_name   string
    Email       string
}
type Vcap struct {
    Uris        []string
}

func oauth_get(url string, tokens oauth2.Tokens) []byte {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("Authorization", "Bearer " + tokens.Access())

    client := http.Client{}
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }

    body, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
    return body
}

func main() {
    var v Vcap
    err := json.Unmarshal([]byte(os.Getenv("VCAP_APPLICATION")), &v)
    if err != nil {
        log.Fatal(err)
    }

	m := martini.Classic()

    app_uri := v.Uris[0]
    cfhost := strings.SplitAfterN(app_uri, ".", 2)[1]
	oauthOpts := &gooauth2.Options{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  "https://" + app_uri + "/oauth2callback",
		Scopes:       []string{""},
	}

	cf := oauth2.NewOAuth2Provider(
        oauthOpts,
        "https://login." + cfhost + "/oauth/authorize",
		"http://uaa." + cfhost + "/oauth/token",
    )

	m.Get("/oauth2error", func() string {
        log.Fatal("oAuth error")
        return "oAuth error :("
	})

	m.Handlers(
		sessions.Sessions("session", sessions.NewCookieStore([]byte("secret123"))),
		cf,
		oauth2.LoginRequired,
		martini.Logger(),
		martini.Static("public"),
	)

	m.Get("/", func(tokens oauth2.Tokens) string {

		if tokens.IsExpired() {
			return "Not logged in, or the access token is expired"
		}

        // oAuth user information from the UAA
        body := oauth_get("http://uaa." + cfhost + "/userinfo", tokens)
        var u User
        err := json.Unmarshal(body, &u)
        if err != nil {
            log.Fatal(err)
        }

        // Example actual API call to get the list of spaces the user belongs to
        space_body := oauth_get("http://api." + cfhost + "/v2/users/" + u.User_id + "/spaces", tokens)

        return fmt.Sprintf("User ID: %s\n\nSpaces info: %s", u.User_id, space_body)
	})

	m.Run()
}
