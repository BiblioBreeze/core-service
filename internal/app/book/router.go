package book

import "github.com/go-chi/chi/v5"

func (s *service) routes(r chi.Router) {
	r.Route("/book", func(r chi.Router) {
		r.Post("/", s.create)
		r.Get("/{id}", s.get)
		r.Get("/list", s.list)
	})
}
