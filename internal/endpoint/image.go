package endpoint

import (
	"context"
	"cloud.google.com/go/storage"
	"encoding/json"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"github.com/shintaro123/ucwork-go/internal"
	"github.com/shintaro123/ucwork-go/internal/model/request"
	"github.com/shintaro123/ucwork-go/internal/repository"
	"github.com/gofrs/uuid"
	"io"
	"net/http"
	"path"
)

func createImageHandler(w http.ResponseWriter, r *http.Request) (string, *appError) {
	// json decode
	decoder := json.NewDecoder(r.Body)
	var imageRequest request.ImageRequest
	err := decoder.Decode(&imageRequest)
	if err != nil {
		return "", appErrorFormat(err, "decode error: %s", err)
	}

	if internal.StorageBucket == nil {
		err := errors.New("storage bucket is missing - check config.go")
		return "", appErrorFormat(err, "error", err)
	}

	// random filename, retaining existing extension.
	name := uuid.Must(uuid.NewV4()).String() + path.Ext(fh.Filename)

	ctx := context.Background()
	w := internal.StorageBucket.Object(name).NewWriter(ctx)

	// Warning: storage.AllUsers gives public read access to anyone.
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"

	// save image to db
	id, err := internal.DBSql.AddImage(image)
	if err != nil {
		return appErrorFormat(err, "add db error: %s", err)
	}

	// create response
	response, jsonError := json.Marshal(image)
	if jsonError != nil {
		return appErrorFormat(jsonError, "%s", jsonError)
	}
	_, writeError := w.Write(response)
	if writeError != nil {
		return appErrorFormat(writeError, "%s", writeError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/images/"+string(id))
	w.WriteHeader(201)
	return nil
}

func uploadFileFromForm(r *http.Request) (url string, err error) {
	f, fh, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	if internal.StorageBucket == nil {
		return "", errors.New("storage bucket is missing - check config.go")
	}

	// random filename, retaining existing extension.
	name := uuid.Must(uuid.NewV4()).String() + path.Ext(fh.Filename)

	ctx := context.Background()
	w := internal.StorageBucket.Object(name).NewWriter(ctx)

	// Warning: storage.AllUsers gives public read access to anyone.
	w.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}
	w.ContentType = fh.Header.Get("Content-Type")

	// Entries are immutable, be aggressive about caching (1 day).
	w.CacheControl = "public, max-age=86400"

	if _, err := io.Copy(w, f); err != nil {
		return "", err
	}
	if err := w.Close(); err != nil {
		return "", err
	}

	const publicURL = "https://storage.googleapis.com/%s/%s"
	return fmt.Sprintf(publicURL, internal.StorageBucketName, name), nil
}
