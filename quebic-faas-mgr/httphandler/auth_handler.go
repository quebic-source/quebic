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
	"log"
	"net/http"
	"quebic-faas/auth"
	"quebic-faas/common"
	"quebic-faas/quebic-faas-mgr/config"
	"quebic-faas/quebic-faas-mgr/dao"
	"quebic-faas/types"

	bolt "github.com/coreos/bbolt"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

const authContext = auth.AuthContext

//AuthHandler handler
func (httphandler *Httphandler) AuthHandler(router *mux.Router) {

	db := httphandler.db
	authConfig := httphandler.config.Auth

	router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {

		authDTO := &types.AuthDTO{}
		err := processRequest(r, authDTO)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		//find saved user from db
		savedUser := types.User{}
		err = dao.GetByID(db, &types.User{Username: authDTO.Username}, func(savedUserINBytes []byte) error {

			if savedUserINBytes == nil {
				return fmt.Errorf("User not found")
			}

			json.Unmarshal(savedUserINBytes, &savedUser)

			return nil
		})
		if err != nil {
			log.Printf("unable to load saved user. error : %v", err)
			makeErrorResponse(w, http.StatusNotFound, err)
			return
		}

		//check for password
		if savedUser.Password != authDTO.Password {
			makeErrorResponse(w,
				http.StatusUnauthorized,
				fmt.Errorf(auth.ErrorMessageInvalidCredentials))
			return
		}

		//jwt token creating
		claims := jwt.MapClaims{
			auth.ClaimsUsername:  savedUser.Username,
			auth.ClaimsFirstname: savedUser.Firstname,
			auth.ClaimsRole:      savedUser.Role,
		}

		tokenString, err := auth.CreateJWTToken(claims, authConfig.JWTSecret)
		if err != nil {
			log.Printf("token creation failed. err : %v", err)
			makeErrorResponse(w,
				http.StatusInternalServerError,
				fmt.Errorf(auth.ErrorMessageInternalError))
			return
		}

		writeResponse(w, types.JWTToken{Token: tokenString}, 200)

	}).Methods("POST")

	router.HandleFunc("/auth/current", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		currentAuthContext := context.Get(r, authContext)
		//authUser := &types.User{}
		//mapstructure.Decode(currentAuthContext.(jwt.MapClaims), authUser)
		writeResponse(w, currentAuthContext, 200)

	}, auth.RoleAny, authConfig)).Methods("GET")

	//add user
	router.HandleFunc("/auth/users", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		user := &types.User{}
		err := processRequest(r, user)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		errors := validateUser(db, user, true)
		if errors != nil {
			status := http.StatusBadRequest
			writeResponse(w, types.ErrorResponse{Cause: common.ErrorValidationFailed, Message: errors, Status: status}, status)
			return
		}

		err = dao.AddUser(db, user)
		if err != nil {
			log.Printf("user save failed. err : %v", err)
			makeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("User save failed"))
			return
		}

		user.Password = ""

		writeResponse(w, user, 200)

	}, auth.RoleAny, authConfig)).Methods("POST")

	//update user
	router.HandleFunc("/auth/users", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		user := types.User{}
		err := processRequest(r, &user)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		authUserName := getAuthUserName(r)

		savedUser, err := getUser(db, authUserName)
		if err != nil {
			log.Printf("unable to load saved user. error : %v", err)
			makeErrorResponse(w, http.StatusInternalServerError, fmt.Errorf(auth.ErrorMessageInternalError))
			return
		}

		savedUser.Firstname = user.Firstname

		err = dao.SaveUser(db, savedUser)
		if err != nil {
			log.Printf("user save failed. err : %v", err)
			makeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("User save failed"))
			return
		}

		writeResponse(w, user, 200)

	}, auth.RoleAny, authConfig)).Methods("PUT")

	//update user
	router.HandleFunc("/auth/users/change-password", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		user := types.User{}
		err := processRequest(r, &user)
		if err != nil {
			makeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		authUserName := getAuthUserName(r)

		savedUser, err := getUser(db, authUserName)
		if err != nil {
			log.Printf("unable to load saved user. error : %v", err)
			makeErrorResponse(w, http.StatusInternalServerError, fmt.Errorf(auth.ErrorMessageInternalError))
			return
		}

		savedUser.Password = user.Password

		err = dao.SaveUser(db, savedUser)
		if err != nil {
			log.Printf("user save failed. err : %v", err)
			makeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("User save failed"))
			return
		}

		writeResponse(w, user, 200)

	}, auth.RoleAny, authConfig)).Methods("POST")

	//get all users
	router.HandleFunc("/auth/users", validateMiddleware(func(w http.ResponseWriter, r *http.Request) {

		var users []types.User
		err := dao.GetAll(db, &types.User{}, func(k, v []byte) error {

			user := types.User{}
			json.Unmarshal(v, &user)

			user.Password = ""

			users = append(users, user)
			return nil
		})

		if err != nil {
			log.Printf("users get failed. err : %v", err)
			makeErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("users get failed"))
			return
		}

		if users == nil {
			var emptyStr [0]string
			writeResponse(w, emptyStr, http.StatusOK)
		} else {
			writeResponse(w, users, http.StatusOK)
		}

	}, auth.RoleAdmin, authConfig)).Methods("GET")

}

func validateMiddleware(next http.HandlerFunc, role string, authConfig config.AuthConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		authorizationHeader := req.Header.Get("authorization")

		if authorizationHeader != "" {

			//jwt token parse
			token, err := jwt.Parse(authorizationHeader, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unable to parse jwt token")
				}
				return []byte(authConfig.JWTSecret), nil
			})
			if err != nil {
				log.Printf("jwt token parse failed. err : %v", err)
				makeErrorResponse(w,
					http.StatusUnauthorized,
					fmt.Errorf(auth.ErrorMessageInvalidToken))
				return

			}

			//if valid token
			if token.Valid {

				context.Set(req, authContext, token.Claims)
				next(w, req)

			} else {

				makeErrorResponse(w,
					http.StatusUnauthorized,
					fmt.Errorf(auth.ErrorMessageAuthFailed))

			}
		} else {
			makeErrorResponse(w,
				http.StatusBadRequest,
				fmt.Errorf(auth.ErrorMessageAuthHeaderNotFound))
		}
	})
}

func getAuthUserName(r *http.Request) string {
	currentAuthContext := context.Get(r, authContext).(jwt.MapClaims)
	return currentAuthContext[auth.ClaimsUsername].(string)
}

func getUser(db *bolt.DB, authUserName string) (*types.User, error) {
	savedUser := &types.User{}
	err := dao.GetByID(db, &types.User{Username: authUserName}, func(savedUserINBytes []byte) error {

		if savedUserINBytes != nil {
			json.Unmarshal(savedUserINBytes, savedUser)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return savedUser, nil

}

func validateUser(db *bolt.DB, user *types.User, isCreate bool) []string {

	var errors []string

	if isCreate {

		if user.Username == "" {
			errors = append(errors, "username should not be empty")
		}

		if user.Password == "" {
			errors = append(errors, "password should not be empty")
		}

		if user.Role == "" {
			errors = append(errors, "role should not be empty")
		}

		if user.Role != "" {
			if user.Role != auth.RoleAdmin && user.Role != auth.RoleDeveloper && user.Role != auth.RoleTester {
				errors = append(errors, "invalide role")
			}
		}

		if user.Username != "" {
			u, _ := getUser(db, user.Username)
			if u.Username != "" {
				errors = append(errors, "user allready exists")
			}
		}

	} else {
		//TODO
	}

	return errors

}
