package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

type LoginForm struct {
	Username string
	Password string
	Errors map[string]string
	Message string
}

func (l *LoginForm) validate() error {
	return validation.ValidateStruct(&l,
	validation.Field(&l.Username,
		validation.Required.Error("The username field is required."),
		),
		validation.Field(&l.Password,
			validation.Required.Error("The password field is required."),
			validation.Length(6, 18).Error("The password field must be between 6 to 18 characters."),
			),
	)
}

func(h *Handler) login(rw http.ResponseWriter, r *http.Request) {
	session, err := h.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	form :=  LoginForm{}


	if flashes := session.Flashes(); len(flashes) > 0 {
		if val, ok := flashes[0].(string); ok {
			form.Message = val
		}
	}
	
	fmt.Println(form.Message)
	if err := h.templates.ExecuteTemplate(rw, "login.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
}

func(h *Handler) loginCheak(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("done")
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	var form LoginForm
	if err := h.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatal(err)
	}

	if err := form.validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			var vErrs map[string]string
			for key, value := range vErrors {
				vErrs[strings.Title(key)] = value.Error()
			}

			form.Errors = vErrs
			
			h.loadLoginForm(rw, form)
			return
		}
	}

	userQuery := `SELECT * FROM users WHERE username = $1`
	var user SignupForm
	h.db.Get(&user, userQuery, form.Username)

	if user.ID == 0 {
		form.Errors = map[string]string{"Username" : "Invalid username given"}
		h.loadLoginForm(rw, form)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		form.Errors = map[string]string{"Username" : "Invalid username given"}
		h.loadLoginForm(rw, form)
		return
	}

	session, err := h.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	session.Options.HttpOnly = true

	session.Values["authUserId"] = user.ID
	if err := session.Save(r, rw); err != nil {
		log.Fatal(err)
	}

	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
}

func (h *Handler) logout(rw http.ResponseWriter, r *http.Request) {
	session, err := h.sess.Get(r, sessionName)
	if err!=nil {
		log.Fatal(err)
	}

	session.Values["authUserId"] = false
	session.Save(r, rw)

	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
}

func(h *Handler) loadLoginForm(rw http.ResponseWriter, form LoginForm) {
	if err := h.templates.ExecuteTemplate(rw, "login.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
}

