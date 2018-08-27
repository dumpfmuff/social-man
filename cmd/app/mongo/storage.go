package mongo

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/osimono/social-man/cmd/app"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

var client *mongo.Client

func init() {
	var err error
	client, err = mongo.NewClient("mongodb://@localhost:27017")
	if err != nil {
		panic("wooot " + err.Error())
	}

}

func AllClients(w http.ResponseWriter, r *http.Request) {
	tenants, err := FetchAllTenants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(tenants)
	w.Write(b)
}

func FetchAllTenants() ([]app.Tenant, error) {
	connectErr := client.Connect(context.Background())
	if connectErr != nil {
		return nil, errors.Wrap(connectErr, "cannot fetch clients")
	}
	collection := client.Database("social-man").Collection("clients")

	cursor, err := collection.Find(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var clients []app.Tenant
	for cursor.Next(context.Background()) {
		document := bson.NewDocument()
		err := cursor.Decode(document)
		if err != nil {
			log.Fatal(err)
		}

		extJSON := document.ToExtJSON(true)
		decoder := json.NewDecoder(bytes.NewReader([]byte(extJSON)))

		var c app.Tenant
		decodeErr := decoder.Decode(&c)
		if decodeErr != nil {
			panic(decodeErr.Error())
		}

		oid := document.Lookup("_id").ObjectID().Hex()
		c.Id = oid

		clients = append(clients, c)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return clients, nil
}

func NewTenant(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var c app.Tenant
	err := decoder.Decode(&c)
	if err != nil {
		panic(err)
	}

	clientWithId, err := StoreClient(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(clientWithId)
	if err != nil {
		panic(err.Error())
	}
}

func StoreClient(c app.Tenant) (app.Tenant, error) {
	connectErr := client.Connect(context.Background())
	if connectErr != nil {
		return app.Tenant{}, errors.Wrap(connectErr, "cannot fetch clients")
	}
	collection := client.Database("social-man").Collection("clients")

	insertOneResult, e := collection.InsertOne(context.Background(), c)

	if e != nil {
		panic(e.Error())
	}
	if oid, ok := insertOneResult.InsertedID.(objectid.ObjectID); ok {
		c.Id = oid.Hex()
	} else {
		// Not objectid.ObjectID, do what you want
	}
	return c, nil
}
