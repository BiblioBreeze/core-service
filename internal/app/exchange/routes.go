package exchange

import "github.com/go-chi/chi/v5"

func (s *service) routes(r chi.Router) {
	r.Route("/exchange", func(r chi.Router) {
		r.Post("/", s.create)
		r.Get("/list", s.list)
	})
}
