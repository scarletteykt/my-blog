package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	mw "github.com/scraletteykt/my-blog/internal/middleware"
	"github.com/scraletteykt/my-blog/internal/middleware/auth"
	"github.com/scraletteykt/my-blog/internal/post"
	"github.com/scraletteykt/my-blog/internal/tag"
	"github.com/scraletteykt/my-blog/internal/user"
	"net/http"
)

type API struct {
	users *user.Users
	posts *post.Posts
	tags  *tag.Tags
}

func New(users *user.Users, posts *post.Posts, tags *tag.Tags) *API {
	return &API{
		users: users,
		posts: posts,
		tags:  tags,
	}
}

func (a *API) Router() chi.Router {
	r := chi.NewRouter()

	r.Use(mw.Middleware(
		mw.WithCORSOptions(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete,
			},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		}),
		mw.WithAuthOptions(auth.Options{
			Secret: "deadbeef",
			Users:  a.users,
		}),
	)...)

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
			r.Post("/", a.CreateTag)
			r.Put("/{tagID}", a.UpdateTag)
			r.Delete("/{tagID}", a.DeleteTag)
		})
	})

	return r
}
