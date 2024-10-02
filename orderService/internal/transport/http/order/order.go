package order

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {

	var body CreateOrderJSONBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	order, err := body.ToOrder()
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte(err.Error()))
		return
	}

	id, err := s.orderService.CreateOrder(r.Context(), order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"id": %d}`, id)))
}

func (s *Server) FindOrderByOrderId(w http.ResponseWriter, r *http.Request, id string) {
	//TODO implement me
	panic("implement me")
}
