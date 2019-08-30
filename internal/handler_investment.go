package internal

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) GetInvestmentByRoomID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	roomID := param.ByName("roomID")

	query := fmt.Sprintf("SELECT id, room_id,product_id,total_payment FROM investments WHERE room_id=$1")

	rows, err := h.DB.Query(query, roomID)
	if err != nil {
		log.Println(err)
		return
	}
	var invests []Investment
	for rows.Next() {
		invest := Investment{}
		err := rows.Scan(&invest.ID,&invest.ProductId,&invest.RoomId,&invest.TotalPayment)
		if err != nil {
			log.Println()
			return
		}
		invests = append(invests, invest)
	}
	bytes, err := json.Marshal(invests)
	if err != nil {
		log.Println()
		return
	}
	renderJSON(w, bytes, http.StatusOK)
}

func (h *Handler) AddNewInvest(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// read json body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Fail to read body"
			`), http.StatusBadRequest)
		return
	}
	var invest Investment
	err = json.Unmarshal(body, &invest)
	if err != nil {
		log.Println(err)
		return
	}
	query := fmt.Sprintf("INSERT INTO investment (id,room_id,product_id,total_payment) VALUES (%d,'%s',%d) ",invest.ID, invest.RoomId, invest.ProductId,invest.TotalPayment)
	_, err = h.DB.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Insert Product Successfully"
	}
	`), http.StatusOK)
}

func (h *Handler) EditProductByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Fail to read body"
			`), http.StatusBadRequest)
		return
	}
	var product Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		log.Println(err)
		return
	}
	// executing insert query
	query := fmt.Sprintf("UPDATE products SET product_name=$1, price=$2 WHERE id=$3 ")
	_, err = h.DB.Query(query,product.ProductName,product.Price,param.ByName("productID"))
	if err != nil {
		log.Println(err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Insert Product Successfully"
	}
	`), http.StatusOK)
}
func (h *Handler) DeleteProductByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	productID := param.ByName("productID")
	query := fmt.Sprintf("DELETE FROM products WHERE id=%s", productID)
	_, err := h.DB.Exec(query)
	if err != nil {
		log.Println(err)
		return
	}
	renderJSON(w, []byte(`
	{
		status:"success",
		message:"Product Deleted Successfully"
	}
	`), http.StatusOK)
}

