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
	mux.Get("/about", dynamicMiddleware.ThenFunc(app.about))

	mux.Get("/post/:id", dynamicMiddleware.ThenFunc(app.showPost))

	// Registration.
	mux.Get("/user/register", dynamicMiddleware.ThenFunc(app.registerUserForm))
	mux.Post("/user/register", dynamicMiddleware.ThenFunc(app.registerUser))

	// User authentication.
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	/*
	 * Admin section.
	 */
	mux.Get("/admin", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.dashboard))

	mux.Get("/admin/posts", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.dashboardAllPosts))
	mux.Get("/admin/post/create", dynamicMiddleware.ThenFunc(app.dashboardCreatePostForm))
	mux.Post("/admin/post/create", dynamicMiddleware.ThenFunc(app.dashboardCreatePost))
	mux.Get("/admin/post/:id", dynamicMiddleware.ThenFunc(app.dashboardUpdatePostForm))
	mux.Post("/admin/post/:id", dynamicMiddleware.ThenFunc(app.dashboardUpdatePost))

	// mux.Get("/admin/page/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPageForm))
	// mux.Post("/admin/page/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createPage))

	mux.Get("/admin/profile", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userProfile))
	mux.Get("/admin/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePasswordForm))
	mux.Post("/admin/change-password", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.changePassword))

	// Static files.
	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
