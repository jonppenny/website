package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))

	mux.Get("/post/:id", dynamicMiddleware.ThenFunc(app.showPost))

	// Registration (disable for now).
	// mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	// mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))

	// User authentication.
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	// Admin section.
	mux.Get("/admin/post/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPostForm))
	mux.Post("/admin/post/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPost))

	// mux.Get("/admin/page/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPageForm))
	// mux.Post("/admin/page/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPage))

	mux.Get("/admin/user/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userProfile))
	mux.Get("/admin/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePasswordForm))
	mux.Post("/admin/user/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePassword))

	// Static files.
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
