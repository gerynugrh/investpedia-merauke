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

func (h *Handler) GetInvestmentByID(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	productID, err := strconv.ParseInt(param.ByName("productID"), 10, 64)
	if err != nil {
		log.Println(err)
		renderJSON(w, []byte(`
		{
			"message":"Gaboleh nakal :)"
		}
		`), http.StatusBadRequest)
		return
	}
	query := fmt.Sprintf("SELECT id, product_name,price FROM products WHERE id=$1")

	rows, err := h.DB.Query(query, productID)
	if err != nil {
		log.Printf("[internal][GetProductByID] fail to select product product_id:%s :%+v\n",
			param.ByName("productID"), err)
		return
	}
	var products []Product
	for rows.Next() {
		product := Product{}
		err := rows.Scan(&product.ID,&product.ProductName, &product.Price)
		if err != nil {
			log.Println()
			return
		}
		products = append(products, product)
	}
	bytes, err := json.Marshal(products)
	if err != nil {
		log.Println()
		return
	}
	renderJSON(w, bytes, http.StatusOK)
}

func (h *Handler) InsertProduct(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// read json body
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
		log.Printf("[internal][InsertProduct] fail to convert json into array :%+v\n",
			err)
		return
	}
	query := fmt.Sprintf("INSERT INTO products (id,product_name,price) VALUES (%d,'%s',%d) ",product.ID, product.ProductName, product.Price)
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

