package db

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/shintaro123/ucwork-go/internal/repository"
)

type datastoreDB struct {
	client *datastore.Client
}

var _ repository.MemberDatabase = &datastoreDB{}

func NewDatastoreDB(client *datastore.Client) (repository.MemberDatabase, error) {
	ctx := context.Background()
	// Verify that we can communicate and authenticate with the datastore service.
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	return &datastoreDB{
		client: client,
	}, nil
}

func (db *datastoreDB) AddMember(member *repository.Member) (id int64, err error){
	ctx := context.Background()
	k := datastore.IncompleteKey("Member", nil)
	k, err = db.client.Put(ctx, k, member)
	if err != nil {
		return 0, fmt.Errorf("datastoredb: could not put Member: %v", err)
	}
	return k.ID, nil
}

func (db *datastoreDB) ListMembers() ([]*repository.Member, error) {
	ctx := context.Background()
	members := make([]*repository.Member, 0)
	q := datastore.NewQuery("Member").
		Order("Name")

	keys, err := db.client.GetAll(ctx, q, &members)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list books: %v", err)
	}

	for i, k := range keys {
		members[i].ID = k.ID
	}

	return members, nil
}
