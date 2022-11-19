package main

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"jonppenny.co.uk/webapp/pkg/forms"
	"jonppenny.co.uk/webapp/pkg/models"
)

// Website Section.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	l, err := app.readInt(qs, "limit", 12)
	if err != nil {
		app.serverError(w, err)
		return
	}

	p, err := app.readInt(qs, "page", 1)
	if err != nil {
		app.serverError(w, err)
		return
	}

	ps, err := app.posts.Latest(l, (p*l)-l)
	if err != nil {
		app.serverError(w, err)
		return
	}

	t, err := app.posts.Total()
	if err != nil {
		app.serverError(w, err)
	}

	pg := &models.Pagination{
		CurrentPage: p,
		TotalPages:  int(math.Ceil(float64(t / l))),
		NextPage:    p + 1,
		PrevPage:    p - 1,
	}

	app.render(w, r, "home.page.tmpl", &templateData{Posts: ps, Pagination: pg}, false, false)
}

func (app *application) showPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "post.page.tmpl", &templateData{Post: s}, false, false)
}

func (app *application) showPage(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Query().Get(":slug")
	if slug == "" {
		app.notFound(w)
		return
	}

	p, err := app.pages.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "page.page.tmpl", &templateData{Page: p}, false, false)
}

// Register section.
func (app *application) registerUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.tmpl", &templateData{
		Form: forms.New(nil),
	}, false, true)
}

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("username", "email", "password")
	form.MaxLength("username", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 8)

	if !form.Valid() {
		app.render(w, r, "register.page.tmpl", &templateData{Form: form}, false, true)
		return
	}

	err = app.users.Insert(form.Get("username"), form.Get("email"), form.Get("password"), "administrator")
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "register.page.tmpl", &templateData{Form: form}, false, true)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Credentials section.
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{Form: forms.New(nil)}, false, true)
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add(
				"generic",
				"Email or Password is incorrect",
			)
			app.render(w, r, "login.page.tmpl", &templateData{Form: form}, false, true)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", id)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You have been logged out successfully.")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Admin Section.
func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "dashboard.page.tmpl", nil, true, false)
}

func (app *application) dashboardAllPosts(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	l, err := app.readInt(qs, "limit", 12)
	if err != nil {
		app.serverError(w, err)
		return
	}

	p, err := app.readInt(qs, "page", 1)
	if err != nil {
		app.serverError(w, err)
		return
	}

	ps, err := app.posts.Latest(l, (p*l)-l)
	if err != nil {
		app.serverError(w, err)
		return
	}

	t, err := app.posts.Total()
	if err != nil {
		app.serverError(w, err)
	}

	pg := &models.Pagination{
		CurrentPage: p,
		TotalPages:  int(math.Ceil(float64(t / l))),
		NextPage:    p + 1,
		PrevPage:    p - 1,
	}

	app.render(w, r, "posts.page.tmpl", &templateData{Posts: ps, Pagination: pg}, true, false)
}

func (app *application) dashboardCreatePostForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create-post.page.tmpl", &templateData{Form: forms.New(nil)}, true, false)
}

func (app *application) dashboardCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "excerpt")
	form.MaxLength("title", 100)
	form.MaxLength("excerpt", 255)
	form.PermittedValues("status", "published", "draft")

	if !form.Valid() {
		app.render(w, r, "create-post.page.tmpl", &templateData{Form: form}, true, false)
		return
	}

	f, err := app.uploadFile(r, "image")
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.posts.Insert(form.Get("title"), form.Get("content"), form.Get("status"), f)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Post created successfully.")

	http.Redirect(w, r, fmt.Sprintf("/admin/post/%d", id), http.StatusSeeOther)
}

func (app *application) dashboardUpdatePostForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	p, err := app.posts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "update-post.page.tmpl", &templateData{Form: forms.New(nil), Post: p}, true, false)
}

func (app *application) dashboardUpdatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content")
	form.MaxLength("title", 100)
	form.PermittedValues("status", "published", "draft")

	if !form.Valid() {
		app.render(w, r, "update-post.page.tmpl", &templateData{Form: form}, true, false)
		return
	}

	/*f, err := app.uploadFile(r, "image")
	if err != nil {
		app.serverError(w, err)
		return
	}*/

	pid, _ := strconv.Atoi(form.Get("post_id"))
	err = app.posts.Update(pid, form.Get("title"), form.Get("content"), form.Get("status"), "", form.Get("excerpt"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Post updated successfully.")

	http.Redirect(w, r, fmt.Sprintf("/admin/post/%d", pid), http.StatusSeeOther)
}

func (app *application) userProfile(w http.ResponseWriter, r *http.Request) {
	userID := app.session.GetInt(r, "authenticatedUserID")

	user, err := app.users.Get(userID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "profile.page.tmpl", &templateData{User: user}, true, false)
}

func (app *application) changePasswordForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "password.page.tmpl", &templateData{Form: forms.New(nil)}, true, false)
}

func (app *application) changePassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("current", "new", "confirm")
	form.MinLength("new", 8)
	if form.Get("new") != form.Get("confirm") {
		form.Errors.Add("confirm", "Passwords do not match")
	}

	if !form.Valid() {
		app.render(w, r, "password.page.tmpl", &templateData{Form: form}, true, false)
		return
	}

	userID := app.session.GetInt(r, "authenticatedUserID")

	err = app.users.Update(userID, form.Get("current"), form.Get("new"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Current Password is incorrect")
			app.render(w, r, "password.page.tmpl", &templateData{Form: form}, true, false)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.session.Put(r, "flash", "Your password has been updated.")

	http.Redirect(w, r, "/admin/profile", http.StatusSeeOther)
}

func (app *application) dashboardAllPages(w http.ResponseWriter, r *http.Request) {
	p, err := app.pages.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "pages.page.tmpl", &templateData{Pages: p}, true, false)
}

func (app *application) dashboardCreatePageForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create-page.page.tmpl", &templateData{Form: forms.New(nil)}, true, false)
}

func (app *application) dashboardCreatePage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "slug", "content")
	form.MaxLength("title", 100)
	form.PermittedValues("status", "published", "draft")

	if !form.Valid() {
		app.render(w, r, "create-page.page.tmpl", &templateData{Form: form}, true, false)
		return
	}

	id, err := app.pages.Insert(form.Get("title"), form.Get("content"), form.Get("status"), form.Get("slug"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Page created successfully.")

	http.Redirect(w, r, fmt.Sprintf("/admin/page/%d", id), http.StatusSeeOther)
}

func (app *application) dashboardUpdatePageForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	p, err := app.pages.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "update-page.page.tmpl", &templateData{Form: forms.New(nil), Page: p}, true, false)
}

func (app *application) dashboardUpdatePage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "slug", "content")
	form.MaxLength("title", 100)
	form.PermittedValues("status", "published", "draft")

	if !form.Valid() {
		app.render(w, r, "update-page.page.tmpl", &templateData{Form: form}, true, false)
		return
	}

	pid, _ := strconv.Atoi(form.Get("page_id"))
	err = app.pages.Update(pid, form.Get("title"), form.Get("content"), form.Get("status"), form.Get("slug"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "flash", "Page updated successfully.")

	http.Redirect(w, r, fmt.Sprintf("/admin/page/%d", pid), http.StatusSeeOther)
}
