package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/zaurbekzet/todo-server/internal/server"
	"github.com/zaurbekzet/todo-server/internal/store"
)

var port = flag.String("p", "9000", "the port the server will listen on")

func main() {
	flag.Parse()

	srv := server.New(store.New())

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Route("/todo", func(r chi.Router) {
		r.Post("/", srv.CreateTodo)
		r.Get("/", srv.GetAll)

		r.Route("/{id:\\d+}", func(r chi.Router) {
			r.Use(idCtx)
			r.Get("/", srv.GetTodo)
			r.Put("/", srv.ToggleTodo)
			r.Delete("/", srv.DeleteTodo)
		})
	})

	log.Fatal(http.ListenAndServe(":"+*port, r))
}

func idCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), server.IDCtxKey{}, uint32(id))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
