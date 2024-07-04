package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"Service/internal/models"
	"Service/internal/validator"

	"github.com/gorilla/mux"
)

type dayCreateForm struct {
	Calories            int    `form:"calories"`
	Proteins            int    `form:"proteins"`
	Fats                int    `form:"fats"`
	Carbs               int    `form:"carbs"`
	Fibre               int    `form:"fibre"`
	Date                string `form:"date"`
	validator.Validator `form:"-"`
}

type userSignUpForm struct {
	IsLoginForm         bool   `form:"-"`
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userSignInForm struct {
	IsLoginForm         bool   `form:"-"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	days, err := app.days.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Days = days

	app.render(w, http.StatusOK, "home.tmpl.html", data)
}

func (app *application) dayCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = dayCreateForm{}
	app.render(w, http.StatusOK, "createDay.tmpl.html", data)
}

func (app *application) dayCreatePost(w http.ResponseWriter, r *http.Request) {
	var form dayCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.MinMaxInteger(form.Calories, 0, 100000), "calories", "This field contains incorrect value")
	form.CheckField(validator.MinMaxInteger(form.Proteins, 0, 100000), "proteins", "This field contains incorrect value")
	form.CheckField(validator.MinMaxInteger(form.Fats, 0, 100000), "fats", "This field contains incorrect value")
	form.CheckField(validator.MinMaxInteger(form.Carbs, 0, 100000), "carbs", "This field contains incorrect value")
	form.CheckField(validator.MinMaxInteger(form.Fibre, 0, 100000), "fibre", "This field contains incorrect value")
	form.CheckField(validator.CheckDate(form.Date, "2006-01-02"), "date", "This field contains incorrect date")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "createDay.tmpl.html", data)
		return
	}

	id, err := app.days.Insert(form.Calories, form.Proteins, form.Fats, form.Carbs, form.Fibre, form.Date)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Day successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/days/view/%d", id), http.StatusSeeOther)
}

func (app *application) dayView(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)["id"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	day, err := app.days.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}

		return
	}

	data := app.newTemplateData(r)
	data.Day = day

	// Use the new render helper.
	app.render(w, http.StatusOK, "view.tmpl.html", data)
}

/* -------------------------------------------------------------------------- */
/*                                signIn signUp                               */
/* -------------------------------------------------------------------------- */

func (app *application) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.IsWithNav = false
	data.Form = userSignUpForm{IsLoginForm: false}
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	var form userSignUpForm
	form.IsLoginForm = false
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using our helper functions.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		data.IsWithNav = false
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			data.IsWithNav = false
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignInForm{IsLoginForm: true}
	data.IsWithFooter = true
	data.IsWithNav = false
	data.IsWithHeader = true
	app.render(w, http.StatusOK, "login.tmpl.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userSignInForm
	form.IsLoginForm = true

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.IsWithNav = false
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			data.IsWithNav = false
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl.html", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the RenewToken() method on the current session to change the session
	// ID. It's good practice to generate a new session ID when the
	// authentication state or privilege levels changes for the user (e.g. login
	// and logout operations).
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/days/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Use the RenewToken() method on the current session to change the session
	// ID again.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Remove the authenticatedUserID from the session data so that the user is
	// 'logged out'.
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) accountView(w http.ResponseWriter, r *http.Request) {
	userID := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	user, err := app.users.Get(userID)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.User = user

	app.render(w, http.StatusOK, "account.tmpl.html", data)
}

func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
