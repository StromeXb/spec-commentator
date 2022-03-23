package gateway

import (
	"context"
	"net/http"
	"os"
	"spec-commentor/internal/openapi"
	"spec-commentor/ms/gateway/generated"
	gMdlwr "spec-commentor/pkg/middleware"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"), //"https://localhost:8080/api/v1/google-callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"), 
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"), 
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"},
		Endpoint:     google.Endpoint,
	}
	// TODO randomize
	randState = "random"
)

func (s Server) AuthLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(randState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s Server) AuthCallback(w http.ResponseWriter, r *http.Request, params generated.AuthCallbackParams) {
	if r.FormValue("state") != randState {
		s.logger.Error().Msg("State is not valid")
		return
	}
	ctx := context.Background()
	token, err := googleOauthConfig.Exchange(ctx, r.FormValue("code"))
	if err != nil {
		s.logger.Error().Msgf("Could not get token %s\n", err.Error())
		return
	}
	response := generated.AuthToken{
		Token: token.AccessToken,
	}
	openapi.Resp(w, http.StatusOK, response)
}

func (s Server) DeleteComment(w http.ResponseWriter, r *http.Request, commentId int) {
	openapi.Resp(w, http.StatusOK, "empty")
}

func (s Server) PostComment(w http.ResponseWriter, r *http.Request, commentId int) {
	openapi.Resp(w, http.StatusOK, "empty")
}

// Метод редактирования комментария
// (PATCH /comments/{comment_id})
func (s Server) PatchComment(w http.ResponseWriter, r *http.Request, commentId int) {
	openapi.Resp(w, http.StatusOK, "empty")
}

// Получить список спек
// (GET /specs)
func (s Server) GetSpecList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cred := ctx.Value(gMdlwr.Userinfo)
	openapi.Resp(w, http.StatusOK, cred)
}

// Метод получения спеки с комментариями
// (GET /specs/{spec_id})
func (s Server) GetSpecWithComments(w http.ResponseWriter, r *http.Request, specId int) {
	openapi.Resp(w, http.StatusOK, "empty")
}
