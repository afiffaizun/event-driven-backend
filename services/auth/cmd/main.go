r.Route("/api/v1", func(r chi.Router) {
	r.Post("/auth/login", authHandler.Login)
	r.Post("/auth/refresh", authHandler.Refresh)
	r.Post("/auth/logout", authHandler.Logout)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(cfg.JWTSecret))
		r.Get("/protected", func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte("protected content"))
		})
	})
})
