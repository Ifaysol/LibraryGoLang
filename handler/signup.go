package handler

import (
	//"fmt"
	"log"
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"

	"golang.org/x/crypto/bcrypt"
)

type SignupFormData struct {
	SignupForm SignupForm
	Errors map[string]string
}


type SignupForm struct {
	ID              int       `db:"id"`
	Username        string    `db:"username"`
	Password        string	  `db:"password"`
	ConfirmPassword string
	FirstName       string    `db:"first_name"`
	LastName        string    `db:"last_name"`
	Email           string    `db:"email"`
	IsVerified      bool      `db:"is_verified"`
}

func (s SignupForm) validate() error {
	return validation.ValidateStruct(&s,
	validation.Field(&s.FirstName,
		validation.Required.Error("The first name field is required."),
		),
		validation.Field(&s.LastName,
			validation.Required.Error("The last name field is required."),
		),
		validation.Field(&s.Email,
			validation.Required.Error("The email field is required."),
			is.Email,
		),
		validation.Field(&s.Username,
			validation.Required.Error("The username field is required."),
		),
		validation.Field(&s.Password,
			validation.Required.Error("The password field is required."),
			validation.Length(6, 18).Error("The password field must be between 6 to 18 characters."),
		),
		validation.Field(&s.ConfirmPassword,
			validation.Required.Error("The confirm password field is required."),
			validation.Length(6, 18).Error("The password field must be between 6 to 18 characters."),
		),
	)
}

func(h *Handler) signup(rw http.ResponseWriter, r *http.Request) {
	formData :=  SignupFormData{}

	h.loadSignupForm(rw, formData)
}

func(h *Handler) register(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatal(err)
	}

	var form SignupForm
	if err := h.decoder.Decode(&form, r.PostForm); err != nil {
		log.Fatal(err)
	}

	if form.Password != form.ConfirmPassword {
		formData := SignupFormData {
			SignupForm: form,
			Errors: map[string]string{"Password" : "The password does not match with confirm password"},
		}
		h.loadSignupForm(rw, formData)
	}

	if err := form.validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			var vErrs map[string]string
			for key, value := range vErrors {
				vErrs[strings.Title(key)] = value.Error()
			}

			formData := SignupFormData {
				SignupForm: form,
				Errors: vErrs,
			}

			h.loadSignupForm(rw, formData)
		}
	}

	const insertUser = `INSERT INTO users(first_name, last_name, email, username, password) VALUES($1, $2, $3, $4, $5);`

	pass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	res := h.db.MustExec(insertUser, form.FirstName, form.LastName, form.Email, form.Username, string(pass))
	
	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return	
	}

	session, err := h.sess.Get(r, sessionName)
	if err != nil {
		log.Fatal(err)
	}

	session.AddFlash("Registration successfully, please login now.")
	if err := session.Save(r, rw); err != nil {
		log.Fatal(err)
	}

	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
	//fmt.Println(1)
}

func(h *Handler) loadSignupForm(rw http.ResponseWriter, formData SignupFormData) {
	if err := h.templates.ExecuteTemplate(rw, "signup.html", formData); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
}