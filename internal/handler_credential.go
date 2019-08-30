package internal

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	// read json body
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		renderJSON(w, []byte(`
			message: "Fail to read body"
			`), http.StatusBadRequest)
		return
	}
	// parse json body
	var credReal Credential
	err = json.Unmarshal(body, &credReal)
	if err != nil {
		log.Printf("[internal][Login] fail to convert json into array :%+v\n",
			err)
		return
	}
	query := fmt.Sprintf("SELECT username,password FROM credentials WHERE username=$1")

	rows, err := h.DB.Query(query, credReal.Username)
	if err != nil {
		log.Printf("[internal][Login] fail to select user user_id:%s :%+v\n",
			credReal.Username, err)
		return
	}
	var cred Credential
	for rows.Next() {
		err := rows.Scan(&cred.Username, &cred.Password)
		if err != nil {
			log.Println(err)
			return
		}
	}
	pwdMatch := comparePasswords(cred.Password,[]byte(credReal.Password))
	if pwdMatch == true{
		renderJSON(w,[]byte(`
	{
		status:"success",
		message:"Login Success"
	}
	`), http.StatusOK)
	}
}

func createToken()

func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}