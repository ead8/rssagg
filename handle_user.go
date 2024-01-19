package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ead8/rssagg/auth"
	"github.com/ead8/rssagg/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *apiConfig)handleCreateUser(w http.ResponseWriter,r *http.Request){

	type parameters struct{
		Name string `json:"name"`
	}	
	decoder:=json.NewDecoder(r.Body)
	params:=parameters{}
	err:=decoder.Decode(&params)
	if err!=nil{
		respondWithError(w,400,fmt.Sprintf("Error Parsing JSON: %v",err))
		return
	}
	usr,err:=apiConfig.DB.GetUser(r.Context(),database.GetUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
	})

	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	repondWithJSON(w,201,databaseUserToUser(usr))
}


func (apiConfig *apiConfig)handleGetUser(w http.ResponseWriter,r *http.Request){

	apiKey,err:= auth.GetAPIKey(r.Header)
	if err!=nil{
		respondWithError(w,403,fmt.Sprintf("Auth error: %v",err))
		return
	}
	user,err:=apiConfig.DB.GetUserByAPIKey(r.Context(),apiKey)
	if err!=nil{
		respondWithError(w,403,fmt.Sprintf("Couldn't find user: %v",err))
		return
	}
	repondWithJSON(w,200,databaseUserToUser(user))	
}