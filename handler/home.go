package handler

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)


type ListBook struct {
	Books []Book
	Offset int 
	Limit int 
	Total int 
	Pagination Pagination
	CurrentPage int
	SearchQuery string
}

type Item struct {
	URL string
	PageNumber int
}

type Pagination struct {
	Items []Item
	NextPageURL string
	PreviousPageURL string
}




func(h *Handler) Home(rw http.ResponseWriter, r *http.Request) {

	page := r.URL.Query().Get("page")

	p, _ := strconv.Atoi(page)
	// if err != nil {
	// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	//fmt.Println(page)
	books := []Book{}
	offset := 0
	limit := 2

	if p > 0 {
		offset = limit * p - limit
	}
	
	total := 0
	h.db.Get(&total, `SELECT count(*) FROM books`) 
    h.db.Select(&books, "SELECT * FROM books offset $1 limit $2", offset, limit)

	totalPage := int(math.Ceil(float64(total)/float64(limit)))

	pagination := Pagination{}

	items := make([]Item, totalPage)

	for i := 0 ; i < totalPage; i++ {
		items[i] = Item{
			URL: fmt.Sprintf("http://localhost:4000?page=%d", i+1),
			PageNumber: i + 1,
		}
		if i + 1 == p {
			if i != 0 {
				pagination.PreviousPageURL = fmt.Sprintf("http://localhost:4000?page=%d", i)
			}
			if i +1 != totalPage {
				pagination.NextPageURL = fmt.Sprintf("http://localhost:4000?page=%d", i+2)
			}
		}
	}

	for key, value := range books {
		const getCategory = `SELECT categoryname FROM categories WHERE id=$1`
		var category Category
		h.db.Get(&category, getCategory, value.CategoryID)
		books[key].Category_Name = category.CategoryName
	}
	
	lb := ListBook{
		Books: books,
		Offset: offset,
		Limit: limit,
		Total: total,
		Pagination: pagination,
		CurrentPage: p,
		SearchQuery: "",
	}

	
	if err :=h.templates.ExecuteTemplate(rw, "list-book.html", lb  ); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	
	
}
		

			
			
