package handler

import "net/http"

type ListCategory struct {
	Categories []Category
	SearchQuery string
}

func(h *Handler) Categoryhome(rw http.ResponseWriter, r *http.Request) {
	categories := []Category{}
    h.db.Select(&categories, "SELECT * FROM categories")
	lc := ListCategory{
		Categories: categories,
		SearchQuery: "",
	}
	if err :=h.templates.ExecuteTemplate(rw, "list-category.html", lc); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	} 
}
