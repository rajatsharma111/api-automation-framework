package helpers

import (
	kms "cloud.google.com/go/kms/apiv1"
	"context"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"io/ioutil"
	"log"
	"os"
)

// Kms consist kmsClient and ctx
type Kms struct {
	kmsClient *kms.KeyManagementClient
	ctx       context.Context
}

var gcpKmsInstance *Kms

// KmsClient creates singlton instace for GcpKms struct
func KmsClient() *Kms {
	if gcpKmsInstance == nil {
		gcpKmsInstance = createClient()
	}
	return gcpKmsInstance
}

// Create the client
func createClient() *Kms {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", os.Getenv("AUTOMATION_PATH")+"/"+os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	ctx := context.Background()
	kmsClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	return &Kms{kmsClient, ctx}
}

// DecryptSymmetric file
func (gkms *Kms) DecryptSymmetric(writefile string, name string, ciphertext []byte) {
	req := &kmspb.DecryptRequest{
		Name:       name,
		Ciphertext: ciphertext,
	}

	result, err := gkms.kmsClient.Decrypt(gkms.ctx, req)
	if err != nil {
		LogError("failed to decrypt ciphertext")
	}
	mode := int(0777)
	err1 := ioutil.WriteFile(writefile, result.Plaintext, os.FileMode(mode))
	if err1 != nil {
		LogError("error writing data to the file")
	}
}
