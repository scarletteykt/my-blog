package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	mw "github.com/scraletteykt/my-blog/internal/middleware"
	"github.com/scraletteykt/my-blog/internal/middleware/auth"
	"github.com/scraletteykt/my-blog/internal/user"
	"net/http"
)

type API struct {
	users *user.Service
}

func New(users *user.Service) *API {
	return &API{
		users: users,
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
			Secret:   "deadbeef",
			Services: h.services,
		}),
	)...)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", h.SignUp) // POST /api/auth/sign-up
			r.Post("/sign-in", h.signIn) // POST /api/auth/sign-in
		})
		r.Route("/posts", func(r chi.Router) {
			r.Get("/", listPosts)   // GET /api/posts
			r.Post("/", createPost) // POST /api/posts
			r.Route("/{postID}", func(r chi.Router) {
				r.Get("/", getPost)       // GET /api/posts/123
				r.Put("/", updatePost)    // PUT /api/posts/123
				r.Delete("/", deletePost) // DELETE /api/posts/123
			})
		})
		r.Route("/tag", func(r chi.Router) {
			r.Post("/", createTag) // POST /api/tag
			r.Route("/{postSlug:[a-z-]+}}", func(r chi.Router) {
				r.Get("/", getTag)       // GET /api/tag/development
				r.Put("/", updateTag)    // PUT /api/tag/development
				r.Delete("/", deleteTag) // DELETE /api/tag/development
			})
		})
	})

	return r
}
