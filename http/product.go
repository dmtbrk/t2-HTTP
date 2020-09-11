package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ortymid/t2-http/market"
)

// ProductHandler forwards product requests to the business logic.
type ProductHandler struct {
	market market.Interface
}

func (h *ProductHandler) RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/", h.List).Methods(http.MethodGet)
	r.HandleFunc("/", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/{id}", h.Detail).Methods(http.MethodGet)
	r.HandleFunc("/{id}", h.Edit).Methods(http.MethodPut)
	r.HandleFunc("/{id}", h.Delete).Methods(http.MethodDelete)
}

// List handles requests for all products.
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.market.Products()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	resp := productListReponse(products)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

// Detail handles requests for the specific product detail.
func (h *ProductHandler) Detail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString, ok := vars["id"]
	if !ok {
		writeError(w, http.StatusBadRequest, errors.New("id not specified"))
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("id is not an integer"))
		return
	}

	product, err := h.market.Product(id)
	if err != nil {
		err = fmt.Errorf("getting product: %w", err)
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if product == nil {
		writeError(w, http.StatusInternalServerError, errors.New("something went wrong"))
		return
	}

	resp := productDetailReponse(*product)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

// Create handles requests for creation of new products.
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(KeyUserID).(string)
	if !ok {
		writeError(w, http.StatusForbidden, errors.New("authorization required"))
		return
	}

	data := struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		err = fmt.Errorf("decoding create product request: %w", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	product := &market.Product{Name: data.Name, Price: data.Price, Seller: userID}
	product, err = h.market.AddProduct(product, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, &market.ErrPermission{}) {
			status = http.StatusBadRequest
		}
		writeError(w, status, err)
		return
	}
	if product == nil {
		writeError(w, http.StatusInternalServerError, errors.New("something went wrong"))
		return
	}

	resp := productCreateReponse(*product)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

// Edit handles product edit requests.
func (h *ProductHandler) Edit(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(KeyUserID).(string)
	if !ok {
		writeError(w, http.StatusForbidden, errors.New("authorization required"))
		return
	}

	id, err := getVarProductID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	product := &market.Product{ID: id}

	data := struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		err = fmt.Errorf("decoding create product request: %w", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	product.Name = data.Name
	product.Price = data.Price

	product, err = h.market.ReplaceProduct(product, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, &market.ErrPermission{}) || errors.Is(err, market.ErrProductNotFound) {
			status = http.StatusBadRequest
		}
		writeError(w, status, err)
		return
	}
	if product == nil {
		writeError(w, http.StatusInternalServerError, errors.New("something went wrong"))
		return
	}

	resp := productEditReponse(*product)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
}

// Delete handles product delete requests.
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(KeyUserID).(string)
	if !ok {
		writeError(w, http.StatusForbidden, errors.New("authorization: token required"))
		return
	}

	id, err := getVarProductID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	// product, err := h.market.Product(id)
	// if err != nil {
	// 	writeError(w, http.StatusInternalServerError, err)
	// 	return
	// }
	// if product == nil {
	// 	writeError(w, http.StatusInternalServerError, errors.New("something went wrong"))
	// 	return
	// }

	err = h.market.DeleteProduct(id, userID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, &market.ErrPermission{}) || errors.Is(err, market.ErrProductNotFound) {
			status = http.StatusBadRequest
		}
		writeError(w, status, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getVarProductID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idString, ok := vars["id"]
	if !ok {
		return 0, errors.New("id not specified")
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, errors.New("id is not an integer")
	}
	return id, nil
}

type productListReponse []*market.Product

func (r productListReponse) MarshalJSON() ([]byte, error) {
	type respProduct struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Seller string `json:"seller"`
	}

	respProducts := make([]respProduct, len(r))
	for i, p := range r {
		respProducts[i] = respProduct{
			ID:     p.ID,
			Name:   p.Name,
			Price:  p.Price,
			Seller: p.Seller,
		}
	}

	return json.Marshal(respProducts)
}

type productDetailReponse market.Product

func (r productDetailReponse) MarshalJSON() ([]byte, error) {
	type respProduct struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Seller string `json:"seller"`
	}

	return json.Marshal(respProduct(r))
}

type productCreateReponse market.Product

func (r productCreateReponse) MarshalJSON() ([]byte, error) {
	type respProduct struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Seller string `json:"seller"`
	}

	return json.Marshal(respProduct(r))
}

type productEditReponse market.Product

func (r productEditReponse) MarshalJSON() ([]byte, error) {
	type respProduct struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Price  int    `json:"price"`
		Seller string `json:"seller"`
	}

	return json.Marshal(respProduct(r))
}
