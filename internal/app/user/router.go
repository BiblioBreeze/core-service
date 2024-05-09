package user

import "github.com/go-chi/chi/v5"

func (s *service) routes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", s.register)
		r.Post("/login", s.login)

		r.With(s.tokenService.AuthMiddleware()).Route("/", func(r chi.Router) {
			r.Get("/", s.get)
			r.Put("/", s.update)
		})
	})
}
