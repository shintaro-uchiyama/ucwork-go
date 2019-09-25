package internal

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/shintaro123/ucwork-go/internal/db"
	"github.com/shintaro123/ucwork-go/internal/repository"
	"log"
)

var (
	DB repository.MemberDatabase
)

func init(){
	var err error
	DB, err = configureDatastore("ucwork-ai-000002")
	if err != nil {
		log.Fatal(err)
	}
}

func configureDatastore(projectID string) (repository.MemberDatabase, error){
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return db.NewDatastoreDB(client)
}

