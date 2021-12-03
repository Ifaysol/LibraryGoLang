package handler

import (
	"io/ioutil"

	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	//"strconv"
)

type FormData struct {
	Book Book
	Category []Category
	Errors map[string] string
}

func (h *Handler) createBook(rw http.ResponseWriter, r *http.Request) {
	category := []Category{}
	h.db.Select(&category, "SELECT * FROM categories")
	vErrs := map[string]string{}
	book := Book{}
	h.loadCreatedBookForm(rw, book, category, vErrs)

}

func (h *Handler) storeBook(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
	category := []Category{}
	h.db.Get(&category, "SELECT * FROM categories")
   
	var book Book
	if err := h.decoder.Decode(&book, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	file, _, err := r.FormFile("Image")
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
    defer file.Close()
   
	img := "upload-*.png"
    tempFile, err := ioutil.TempFile("assets/images", img)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
    defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
    }
	
    tempFile.Write(fileBytes)
	a := tempFile.Name()

	if err := book.validate(); err != nil {
		
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErrs := make(map[string]string)
			for key, value := range  vErrors {
				vErrs[key] = value.Error()
			}
			//fmt.Println(vErrs)
			h.loadCreatedBookForm(rw, book, category, vErrs)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	
	
		const insertBook = `INSERT INTO books(bookName, catid, image, is_available) VALUES($1, $2, $3, $4);`
		res := h.db.MustExec(insertBook, book.BookName, book.CategoryID, a, false)
		
		if ok, err := res.RowsAffected(); err != nil || ok == 0 {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		    return	
		}
	
		http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
		
}

func (h *Handler) availableBook(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["book"]
	if id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}
	
	const availableBook = `UPDATE books SET is_available = true WHERE id = $1`
	res := h.db.MustExec(availableBook, id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return	
	}


	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)

}

func (h *Handler) editBook(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//fmt.Println(vars)
	id := vars["book"]
	

	if id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getBook = `SELECT * FROM books WHERE id = $1`
	var book Book
	h.db.Get(&book, getBook, id)

	category := []Category{}
	h.db.Select(&category, "SELECT * FROM categories")
	// fmt.Println(category)

	if book.ID == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	h.loadEditBookForm(rw, book, category, map[string]string{})

}

func (h *Handler) updateBook(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["book"]
    //fmt.Println(id)
	if id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getBook = `SELECT * FROM books WHERE id = $1`
	var book Book
	h.db.Get(&book, getBook, id)

	category := []Category{}
	h.db.Select(&category, "SELECT * FROM categories")

	

	if book.ID == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 



	bookname := r.FormValue("BookName")

	book.BookName = bookname


	if bookname == "" {
		vErrs := map[string]string{
			"BookName" : "The book name field is required",
		}
		h.loadEditBookForm(rw, book, category, vErrs)
		return
	}

	if len(bookname) < 3 {
		vErrs := map[string]string{
			"BookName" : "The book name field must be greater than or equals 3",
		}
		h.loadEditBookForm(rw, book, category, vErrs)
		return
	} 
	
	const availableBook = `UPDATE books SET bookname = $2 WHERE id = $1`
	res := h.db.MustExec(availableBook, id, bookname)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return	
	}

	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)

}

func (h *Handler) deleteBook(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["book"]

	if id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getBook = `SELECT * FROM books WHERE id = $1`
	var book Book
	h.db.Get(&book, getBook, id)

	

	if book.ID == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const deleteBook = `DELETE FROM books WHERE id = $1`
	res := h.db.MustExec(deleteBook, id)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return	
	}

	http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)

}

func (h *Handler) searchBook(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
	// fmt.Println("done")

	res:= r.FormValue("search")
	// fmt.Println(res)

	const getValue= `SELECT * FROM books WHERE bookname ILIKE '%%' || $1 || '%%'`
	var b []Book
	h.db.Select(&b,getValue,res)

	for key, value := range b {
		const getCategory = `SELECT categoryname FROM categories WHERE id=$1`
		var category Category
		h.db.Get(&category, getCategory, value.CategoryID)
		b[key].Category_Name = category.CategoryName
	}

	lt:= ListBook{
		Books: b,
		SearchQuery: r.FormValue("search"),
	}
     //fmt.Println(lt)
	if err := h.templates.ExecuteTemplate(rw, "list-book.html", lt); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
}

func (h *Handler) loadCreatedBookForm(rw http.ResponseWriter, book Book, category []Category, errs map[string]string) {
	form := FormData {
		Book : book,
		Category: category,
		Errors: errs,	
	}

	if err := h.templates.ExecuteTemplate(rw, "create-book.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
}
	

func (h *Handler) loadEditBookForm(rw http.ResponseWriter, book Book, category []Category, errs map[string]string) {
	form := FormData {
		Book : book,
		Category: category,
		Errors: errs,	
	}
	//fmt.Println(form)

	if err := h.templates.ExecuteTemplate(rw, "edit-book.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
}
		



