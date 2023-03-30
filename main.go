package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	iam "google.golang.org/api/iam/v1"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	f, err := os.Create(path + "/gcp/gcp-credentials.json")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	w := bufio.NewWriter(f)

	defer w.Flush()

	_, err = createKey(w, "infra-builder@infra-builder-380512.iam.gserviceaccount.com")
	if err != nil {
		panic(err)
	}
}

// createKey creates a service account key.
func createKey(w io.Writer, serviceAccountEmail string) (*iam.ServiceAccountKey, error) {
	ctx := context.Background()
	service, err := iam.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("iam.NewService: %v", err)
	}

	resource := "projects/-/serviceAccounts/" + serviceAccountEmail
	request := &iam.CreateServiceAccountKeyRequest{}
	key, err := service.Projects.ServiceAccounts.Keys.Create(resource, request).Do()
	if err != nil {
		return nil, fmt.Errorf("Projects.ServiceAccounts.Keys.Create: %v", err)
	}
	// The PrivateKeyData field contains the base64-encoded service account key
	// in JSON format.
	// TODO(Developer): Save the below key (jsonKeyFile) to a secure location.
	// You cannot download it later.
	jsonKeyFile, _ := base64.StdEncoding.DecodeString(key.PrivateKeyData)
	_, err = w.Write(jsonKeyFile)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(w, "Key created successfully")
	return key, nil
}
