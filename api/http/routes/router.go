package routes

import (
	"budgetting/api/http/handler"
	"budgetting/config"
	jwtMiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	adapter "github.com/gwatts/gin-adapter"
	"net/http"
	"net/url"
	"time"
)

type Router struct {
	AuthMiddlewareFunc gin.HandlerFunc
	Config             *config.Config
}

func NewRouter(config *config.Config) *Router {
	return &Router{Config: config}
}

func (r *Router) SetupRoutes(e *gin.Engine, h *handler.Handler) {
	r.SetupAuthMiddleware()
	api := e.Group("/api/v1")

	r.SetupHealthRoutes(api, h)
	r.SetupBankingRoutes(api, h)
	r.SetupProfileRoutes(api, h)
}

func (r *Router) SetupAuthMiddleware() {
	issuerURL, _ := url.Parse(r.Config.Auth0IssuerUrl)
	audience := r.Config.Auth0Audience

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)
	customClaims := func() validator.CustomClaims {
		return &handler.CustomClaims{}
	}

	jwtValidator, _ := validator.New(provider.KeyFunc,
		validator.RS256,
		issuerURL.String()+"/",
		[]string{audience},
		validator.WithCustomClaims(customClaims),
	)

	authMiddleware := jwtMiddleware.New(jwtValidator.ValidateToken,
		jwtMiddleware.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			w.WriteHeader(http.StatusUnauthorized)
		}))

	r.AuthMiddlewareFunc = adapter.Wrap(authMiddleware.CheckJWT)
}
