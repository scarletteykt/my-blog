package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/scraletteykt/my-blog/internal/config"
	mw "github.com/scraletteykt/my-blog/internal/middleware"
	"github.com/scraletteykt/my-blog/internal/middleware/auth"
	"github.com/scraletteykt/my-blog/internal/post"
	"github.com/scraletteykt/my-blog/internal/tag"
	"github.com/scraletteykt/my-blog/internal/user"
)

type API struct {
	cfg   *config.Config
	users *user.Users
	posts *post.Posts
	tags  *tag.Tags
}

func New(cfg *config.Config, users *user.Users, posts *post.Posts, tags *tag.Tags) *API {
	return &API{
		cfg:   cfg,
		users: users,
		posts: posts,
		tags:  tags,
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
				r.Get("/posts", a.GetPostsByTag)
				r.Put("/", a.UpdateTag)
				r.Delete("/", a.DeleteTag)
			})
		})
	})

	return r
}
