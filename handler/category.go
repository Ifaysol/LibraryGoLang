package handler

import (
	"net/http"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"

)

type Category struct {
	ID           int    `db:"id"`
	CategoryName string `db:"categoryname"`
}

func (c *Category) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.CategoryName, validation.Required.Error("Must insert a category name"), validation.Length(3, 0)),
	)
}

type FData struct {
	Category Category
	Errors   map[string]string
}

type searchCategory struct {
	Category []Category
}

func (h *Handler) createCategory(rw http.ResponseWriter, r *http.Request) {
	vErrs := map[string]string{}
	category := Category{}
	h.loadCreatedCategoryForm(rw, category, vErrs)

}

func (h *Handler) storeCategory(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	var category Category
	if err := h.decoder.Decode(&category, r.PostForm); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := category.Validate(); err != nil {
		vErrors, ok := err.(validation.Errors)
		if ok {
			vErr := make(map[string]string)
			for key, value := range vErrors {
				vErr[strings.Title(key)] = value.Error()

			}

			h.loadCreatedCategoryForm(rw, category, vErr)
			return

		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return

	}

	const insertCategory = `INSERT INTO categories(categoryname) VALUES($1);`
	res := h.db.MustExec(insertCategory, category.CategoryName)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(rw, r, "/categories", http.StatusTemporaryRedirect)

}

func (h *Handler) editCategory(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["category"]
	//fmt.Println(id)
	

	if id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getCategory = `SELECT * FROM categories WHERE id = $1`
	var category Category
	h.db.Get(&category, getCategory, id)

	if category.ID == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	h.loadEditCategoryForm(rw, category, map[string]string{})

}

func (h *Handler) updateCategory(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["category"]
    //fmt.Println(id)
	if id == "" {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	const getCategory = `SELECT * FROM categories WHERE id = $1`
	var category Category
	h.db.Get(&category, getCategory, id)

	

	if category.ID == 0 {
		http.Error(rw, "Invalid URL", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 



	categoryname := r.FormValue("CategoryName")

	category.CategoryName = categoryname


	if categoryname == "" {
		vErrs := map[string]string{
			"CategoryName" : "The category name field is required",
		}
		h.loadEditCategoryForm(rw, category, vErrs)
		return
	}

	if len(categoryname) < 3 {
		vErrs := map[string]string{
			"CategoryName" : "The category name field must be greater than or equals 3",
		}
		h.loadEditCategoryForm(rw, category, vErrs)
		return
	} 
	
	const availableCategory = `UPDATE categories SET categoryname = $2 WHERE id = $1`
	res := h.db.MustExec(availableCategory, id, categoryname)

	if ok, err := res.RowsAffected(); err != nil || ok == 0 {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return	
	}

	http.Redirect(rw, r, "/categories", http.StatusTemporaryRedirect)

}


func (h *Handler) searchCategory(rw http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("done")

	res := r.FormValue("search")
	//fmt.Println(res)

	const getValue = `SELECT * FROM categories WHERE categoryname ILIKE '%%' || $1 || '%%'`
	var b []Category
	h.db.Select(&b, getValue, res)

	lt := ListCategory{
		Categories: b,
		SearchQuery: r.FormValue("search"),
	}
	//fmt.Println(lt)

	if err := h.templates.ExecuteTemplate(rw, "list-category.html", lt); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadCreatedCategoryForm(rw http.ResponseWriter, category Category, errs map[string]string) {
	form := FData{
		Category: category,
		Errors:   errs,
	}

	if err := h.templates.ExecuteTemplate(rw, "create-category.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadEditCategoryForm(rw http.ResponseWriter,category Category, errs map[string]string) {
	form := FData {
		Category: category,
		Errors: errs,	
	}
	//fmt.Println(form)

	if err := h.templates.ExecuteTemplate(rw, "edit-category.html", form); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
}
