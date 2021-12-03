package handler

import (
	// "fmt"
	// "fmt"
	// "fmt"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Booking struct {
	ID         int    `db:"id"`
	User_id    int    `db:"user_id"`
	Book_id    int    `db:"book_id"`
	S_time string
	E_time   string 
	Start_time     time.Time  `db:"start_time"`
	End_time     time.Time  `db:"end_time"`
	Book_name  string
}

type BookingForm struct {
	Bookings Booking
	ID       string
	Errors   map[string]string
}

func (h *Handler) createBooking(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	fmt.Println(id)
	Errors := map[string]string{}
	book_id := id
	h.loadCreatedBookingForm(rw, book_id, Errors)
}
func (h *Handler) loadCreatedBookingForm(rw http.ResponseWriter,data string, errs map[string]string) {

	form := BookingForm{
		ID: data,
	}
	// fmt.Println(form)
	if err := h.templates.ExecuteTemplate(rw, "create-booking.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) storeBooking(rw http.ResponseWriter, r *http.Request) {
	
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var booking Booking
	if err := h.decoder.Decode(&booking, r.PostForm); err != nil {

		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
    b_id:=booking.Book_id
	const insertBooking = `INSERT INTO bookings (user_id,book_id,start_time,end_time) VALUES ($1,$2,$3,$4)`

	res := h.db.MustExec(insertBooking, booking.User_id, booking.Book_id, booking.S_time, booking.E_time)
	
	const completedBook = `UPDATE books SET is_available = true WHERE id = $1`
		 h.db.MustExec( completedBook, b_id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	//fmt.Println(res)

	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)

}