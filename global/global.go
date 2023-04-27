package global

import (
	"context"
	"log"
	"os"
	"sync"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/zicops/zicops-cass-pool/cassandra"
	cry "github.com/zicops/zicops-vilt-manager/lib/crypto"
	"github.com/zicops/zicops-vilt-manager/lib/identity"
	"github.com/zicops/zicops-vilt-manager/lib/sendgrid"
	"google.golang.org/api/option"
)

var (
	IDP             *identity.IDP
	CTX             context.Context
	Cancel          context.CancelFunc
	CryptSession    *cry.Cryptography
	CassPool        *cassandra.CassandraPool
	SGClient        *sendgrid.ClientSendGrid
	WaitGroupServer sync.WaitGroup
	App             *firebase.App
	Client          *firestore.Client
	Messanger       *messaging.Client
)

func init() {
	serviceAccountZicops := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if serviceAccountZicops == "" {
		log.Printf("failed to get right credentials for course creator")
	}
	targetScopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}
	currentCreds, _, err := identity.ReadCredentialsFile(CTX, serviceAccountZicops, targetScopes)
	if err != nil {
		log.Println(err)
	}

	opt := option.WithCredentials(currentCreds)
	App, err = firebase.NewApp(CTX, nil, opt)
	if err != nil {
		log.Printf("error initializing app: %v", err)
	}

	Client, err = App.Firestore(CTX)
	if err != nil {
		log.Printf("Error while initialising firestore %v", err)
	}

	messanger, err := App.Messaging(CTX)
	if err != nil {
		log.Printf("Error while initialising messaging: %v", err)
	}
	Messanger = messanger
}
