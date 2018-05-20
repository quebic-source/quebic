//    Copyright 2018 Tharanga Nilupul Thennakoon
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package httphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/types"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

const jwtSecret = "SWaBjrk552D2LnF6"

//AuthContext context key for access auth object
const AuthContext = "auth"

//AuthHandler handler
func (httphandler *Httphandler) AuthHandler(router *mux.Router) {

	db := httphandler.db

	router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {

		authDTO := &types.AuthDTO{}
		err := processRequestAuth(r, authDTO)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		//find saved user from db
		savedUser := types.User{}
		err = dao.GetByID(db, &types.User{Username: authDTO.Username}, func(savedUserINBytes []byte) error {

			if savedUserINBytes != nil {
				json.Unmarshal(savedUserINBytes, &savedUser)
			}

			return nil
		})
		if err != nil {
			makeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		//check for password
		if savedUser.Password != authDTO.Password {
			makeErrorResponse(w,
				http.StatusUnauthorized,
				fmt.Errorf("authentication failed. credentials not match"))
			return
		}

		//jwt token creating
		claims := jwt.MapClaims{
			"username":  savedUser.Username,
			"firstname": savedUser.Firstname,
		}
		tokenString, err := common.CreateJWTToken(claims, jwtSecret)
		if err != nil {
			makeErrorResponse(w,
				http.StatusInternalServerError,
				fmt.Errorf("token creation failed. err : %v", err))
			return
		}

		writeResponse(w, types.JWTToken{Token: tokenString}, 200)

	}).Methods("POST")

	router.HandleFunc("/auth/current", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		currentAuthContext := context.Get(r, AuthContext)
		//mapstructure.Decode(currentAuthContext.(jwt.MapClaims), &A)
		writeResponse(w, currentAuthContext.(jwt.MapClaims), 200)

	})).Methods("POST")

}

func validateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		authorizationHeader := req.Header.Get("authorization")

		if authorizationHeader != "" {

			bearerToken := strings.Split(authorizationHeader, " ")

			if len(bearerToken) == 2 {

				//jwt token parse
				token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unable to parse jwt token")
					}
					return []byte(jwtSecret), nil
				})
				if err != nil {
					log.Printf("jwt token parse failed. err : %v", err)
					makeErrorResponse(w,
						http.StatusUnauthorized,
						fmt.Errorf("invalid authorization token"))
					return

				}

				//if valid token
				if token.Valid {

					context.Set(req, AuthContext, token.Claims)
					next(w, req)

				} else {

					makeErrorResponse(w,
						http.StatusUnauthorized,
						fmt.Errorf("invalid authorization token"))

				}
			}
		} else {
			makeErrorResponse(w,
				http.StatusBadRequest,
				fmt.Errorf("authorization header is required"))
		}
	})
}

func processRequestAuth(r *http.Request, authDTO *types.AuthDTO) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return makeError("unable to read request %v", err)
	}

	err = json.Unmarshal(body, authDTO)
	if err != nil {
		return makeError("unable to parse json request to authDTO %v", err)
	}

	return nil

}
