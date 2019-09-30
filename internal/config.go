package internal

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/storage"
	"context"
	"github.com/shintaro123/ucwork-go/internal/db"
	"github.com/shintaro123/ucwork-go/internal/repository"
	"log"
	"os"
)

var (
	DB repository.MemberDatabase
	DBSql repository.OrderDatabase

	StorageBucket     *storage.BucketHandle
	StorageBucketName string

)

func init(){
	var err error
	// Cloud Datastoreの初期設定
	DB, err = configureDatastore("ucwork-ai-000002")
	if err != nil {
		log.Fatal(err)
	}

	// Cloud SQLの初期設定
	DBSql, err = configureCloudSQL(cloudSQLConfig{
		Username: "root",
		Password: "root",
		Instance: "ucwork-ai-000002:asia-northeast1:ucwork",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Cloud Storageの初期設定
	StorageBucketName = "<your-storage-bucket>"
	StorageBucket, err = configureStorage(StorageBucketName)
}

func configureDatastore(projectID string) (repository.MemberDatabase, error){
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return db.NewDatastoreDB(client)
}

type cloudSQLConfig struct {
	Username, Password, Instance string
}

func configureCloudSQL(config cloudSQLConfig) (repository.OrderDatabase, error) {
	if os.Getenv("GAE_INSTANCE") != "" {
		// Running in production.
		return db.NewMySQLDB(db.MySQLConfig{
			Username:   config.Username,
			Password:   config.Password,
			UnixSocket: "/cloudsql/" + config.Instance,
		})
	}

	// Running locally.
	return db.NewMySQLDB(db.MySQLConfig{
		Username: config.Username,
		Password: config.Password,
		Host:     "localhost",
		Port:     3306,
	})
}

func configureStorage(bucketID string) (*storage.BucketHandle, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return client.Bucket(bucketID), nil
}

