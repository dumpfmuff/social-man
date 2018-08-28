package mongo

import (
	"context"
	mgo "github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/osimono/social-man"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

var storage Storage

func init() {
	bundle := clientopt.ClientBundle{}
	timeout := time.Duration(2) * time.Second
	client, _ := mgo.NewClientWithOptions("mongodb://@localhost:27017",
		bundle.ServerSelectionTimeout(timeout),
		bundle.ConnectTimeout(timeout),
		bundle.SocketTimeout(timeout))
	connectErr := client.Connect(context.Background())
	if connectErr != nil {
		panic("no mongo")
	}
	storage = NewStorage(client)
}

func TestSomething(t *testing.T) {
	storage.StoreTenant(social_man.Tenant{Lastname: "social", Surname: "man"})
	tenants, _ := storage.FetchAllTenants()

	logrus.Infof("%v", tenants)
}
