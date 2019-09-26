package endpoint

import (
	"encoding/json"
	"github.com/shintaro123/ucwork-go/internal"
	"github.com/shintaro123/ucwork-go/internal/model/request"
	"github.com/shintaro123/ucwork-go/internal/repository"
	"net/http"
)

func listOrdersHandler(w http.ResponseWriter, r *http.Request) *appError {
	orders, err := internal.DBSql.ListOrders()
	if err != nil {
		return appErrorFormat(err, "%s", err)
	}

	response, jsonError := json.Marshal(orders)
	if jsonError != nil {
		return appErrorFormat(jsonError, "%s", jsonError)
	}

	_, writeError := w.Write(response)
	if writeError != nil {
		return appErrorFormat(writeError, "%s", writeError)
	}
	w.Header().Set("Content-Type", "application/json")
	return nil
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) *appError {
	// json decode
	decoder := json.NewDecoder(r.Body)
	var orderRequest request.OrderRequest
	err := decoder.Decode(&orderRequest)
	if err != nil {
		return appErrorFormat(err, "decode error: %s", err)
	}

	// object convert
	order, err := orderFromJson(&orderRequest)
	if err != nil {
		return appErrorFormat(err, "convert error: %s", err)
	}

	// save order to db
	id, err := internal.DBSql.AddOrder(order)
	if err != nil {
		return appErrorFormat(err, "add db error: %s", err)
	}

	// create response
	response, jsonError := json.Marshal(order)
	if jsonError != nil {
		return appErrorFormat(jsonError, "%s", jsonError)
	}
	_, writeError := w.Write(response)
	if writeError != nil {
		return appErrorFormat(writeError, "%s", writeError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/orders/"+string(id))
	w.WriteHeader(201)
	return nil
}

func orderFromJson(orderRequest *request.OrderRequest) (*repository.Order, error) {
	order := &repository.Order{
		Name: orderRequest.Name,
	}
	return order, nil
}

