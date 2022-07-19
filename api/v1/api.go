package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/scraletteykt/my-blog/internal/config"
	mw "github.com/scraletteykt/my-blog/internal/middleware"
	"github.com/scraletteykt/my-blog/internal/middleware/auth"
	"github.com/scraletteykt/my-blog/internal/service"
	"github.com/scraletteykt/my-blog/pkg/logger"
)

type API struct {
	cfg   *config.Config
	users *service.UsersService
	posts *service.PostsService
	tags  *service.TagsService
	log   logger.Logger
}

func NewAPI(cfg *config.Config, users *service.UsersService, posts *service.PostsService, tags *service.TagsService, log logger.Logger) *API {
	return &API{
		cfg:   cfg,
		users: users,
		posts: posts,
		tags:  tags,
		log:   log,
	}
}

func (a *API) Router() chi.Router {
	r := chi.NewRouter()

	mwAuth := auth.New(&auth.Config{Secret: a.cfg.Auth.Secret}, a.users)
	r.Use(mwAuth.Handler)
	r.Use(mw.Middleware()...)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", a.SignUp)
			r.Post("/sign-in", a.SignIn)
		})
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", a.GetPosts)
			r.Post("/", a.CreatePost)
			r.Route("/{postID}", func(r chi.Router) {
				r.Get("/", a.GetPostByID)
				r.Put("/", a.UpdatePost)
				r.Delete("/", a.DeletePost)
			})
		})
		r.Route("/tags", func(r chi.Router) {
			r.Get("/", a.GetTags)
			r.Post("/", a.CreateTag)
			r.Route("/{tagID}", func(r chi.Router) {
				r.Get("/", a.GetTagByID)
				r.Get("/posts", a.GetPostsByTag)
				r.Put("/", a.UpdateTag)
				r.Delete("/", a.DeleteTag)
			})
		})
	})

	return r
}
