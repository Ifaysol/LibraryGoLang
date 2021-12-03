package handler

import (
	//"go/token"
	"net/http"

	"text/template"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

const sessionName = "library-session"

type Book struct {
	ID            int    `db:"id"`
	CategoryID    int    `db:"catid"`
	BookName      string `db:"bookname"`
	Image         string `db:"image"`
	IsAvailable   bool   `db:"is_available"`
	Category_Name string
}

func (b *Book) validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.BookName,
			validation.Required.Error("The book name must be needed"),
			validation.Length(3, 0),
		),
	)
}

type Handler struct {
	templates *template.Template
	db        *sqlx.DB
	decoder   *schema.Decoder
	sess      *sessions.CookieStore
}

func New(db *sqlx.DB, decoder *schema.Decoder, sess *sessions.CookieStore) *mux.Router {
	h := &Handler{
		db:      db,
		decoder: decoder,
		sess:    sess,
	}

	h.parseTemplate()

	r := mux.NewRouter()

    r.HandleFunc("/", h.Home)
	r.HandleFunc("/login",  h.login)
	r.HandleFunc("/login/check",  h.loginCheak)
	r.HandleFunc("/signup", h.signup)
	r.HandleFunc("/signup/check", h.register)
	r.HandleFunc("/logout", h.logout)
	s := r.NewRoute().Subrouter()

	
	s.HandleFunc("/books/create", h.createBook)
	s.HandleFunc("/books/store", h.storeBook)
	s.HandleFunc("/books/{book:[0-9]+}/available", h.availableBook)
	s.HandleFunc("/books/{book:[0-9]+}/edit", h.editBook)
	s.HandleFunc("/books/{book:[0-9]+}/update", h.updateBook)
	s.HandleFunc("/books/{book:[0-9]+}/delete", h.deleteBook)
	s.HandleFunc("/books/search", h.searchBook)
	s.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))




	s.HandleFunc("/categories", h.Categoryhome)
	s.HandleFunc("/categories/create", h.createCategory)
	s.HandleFunc("/categories/store", h.storeCategory)
	s.HandleFunc("/categories/{category}/edit", h.editCategory)
	s.HandleFunc("/categories/{category}/update", h.updateCategory)
	//r.HandleFunc("/books/{category}/delete", h.deleteCategory)
	s.HandleFunc("/categories/search", h.searchCategory)

	s.HandleFunc("/booking/{id}/create", h.createBooking)
	s.HandleFunc("/booking/store", h.storeBooking)
	s.Use(h.authMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if err := h.templates.ExecuteTemplate(rw, "404.html", nil); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})
	return r
}

func (h *Handler) parseTemplate() {
	h.templates = template.Must(template.ParseFiles(
		"templates/create-book.html",
		"templates/list-book.html",
		"templates/edit-book.html",
		"templates/create-category.html",
		"templates/list-category.html",
		"templates/edit-category.html",
		"templates/404.html",
		"templates/create-booking.html",
		"templates/login.html",
		"templates/signup.html",
	))
}

func(h *Handler) authMiddleware (next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		session, _ := h.sess.Get(r, sessionName)
		// if err != nil {
		// 	http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
		// 	return
		// }

		// Check if user is authenticated
		ok := session.Values["authUserId"]
		if ok == nil {
			http.Redirect(rw, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(rw, r)
	})
}
