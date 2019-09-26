package endpoint

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/shintaro123/ucwork-go/internal"
	"github.com/shintaro123/ucwork-go/internal/model/request"
	"github.com/shintaro123/ucwork-go/internal/repository"
	"net/http"
)

type Member struct {
	Name string
}

type Members []Member

func listHandler(w http.ResponseWriter, r *http.Request) *appError {
	members, err := internal.DB.ListMembers()
	if err != nil {
		return appErrorFormat(err, "%s", err)
	}

	response, jsonError := json.Marshal(members)
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

func createHandler(w http.ResponseWriter, r *http.Request) *appError {
	// json decode
	decoder := json.NewDecoder(r.Body)
	var memberRequest request.MemberRequest
	err := decoder.Decode(&memberRequest)
	if err != nil {
		return appErrorFormat(err, "decode error: %s", err)
	}

	// object convert
	member, err := memberFromJson(&memberRequest)
	if err != nil {
		return appErrorFormat(err, "convert error: %s", err)
	}

	// save member to db
	id, err := internal.DB.AddMember(member)
	if err != nil {
		return appErrorFormat(err, "add db error: %s", err)
	}

	// create response
	response, jsonError := json.Marshal(member)
	if jsonError != nil {
		return appErrorFormat(jsonError, "%s", jsonError)
	}
	_, writeError := w.Write(response)
	if writeError != nil {
		return appErrorFormat(writeError, "%s", writeError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/members/"+string(id))
	w.WriteHeader(201)
	return nil
}

func memberFromJson(memberRequest *request.MemberRequest) (*repository.Member, error) {
	member := &repository.Member{
		Name: memberRequest.Name,
	}
	return member, nil
}

func updateHandler(w http.ResponseWriter, r *http.Request) *appError {
	response, jsonError := json.Marshal(Members{
		Member{
			Name: "updated Name " + mux.Vars(r)["id"],
		},
	})
	if jsonError != nil {
		return appErrorFormat(jsonError, "%s", jsonError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, writeError := w.Write(response)
	if writeError != nil {
		return appErrorFormat(writeError, "%s", writeError)
	}
	return nil
}

func deleteHandler(w http.ResponseWriter, r *http.Request) *appError {
	requestId := mux.Vars(r)["id"]
	if requestId == "2" {
		return appErrorFormat(errors.New("invalid reques id"), "invalid id:  %s", requestId)
	}
	w.Header().Set("Content-Type", "application/json")
	return nil
}

