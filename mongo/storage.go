package mongo

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/osimono/social-man"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	Client  *mongo.Client
	tenants *mongo.Collection
}

func NewStorage(c *mongo.Client) Storage {
	collection := c.Database("social-man").Collection("tenants")
	return Storage{Client: c, tenants: collection}
}

func (s *Storage) FetchAllTenants() ([]social_man.Tenant, error) {
	cursor, err := s.tenants.Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var clients []social_man.Tenant
	for cursor.Next(context.Background()) {
		document := bson.NewDocument()
		err := cursor.Decode(document)
		if err != nil {
			return nil, err
		}
		var c social_man.Tenant
		parseFromDocument(document, &c)

		clients = append(clients, c)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return clients, nil
}

func (s *Storage) FindTenant(id string) (social_man.Tenant, error) {
	logrus.Infof("looking for tenant with id: %v", id)
	var t social_man.Tenant
	result := s.tenants.FindOne(context.Background(), filterDocument(id))
	document := bson.NewDocument()

	if decodeErr := result.Decode(document); decodeErr != nil {
		return t, decodeErr
	}

	if parseErr := parseFromDocument(document, &t); parseErr != nil {
		return t, parseErr
	}

	logrus.Infof("found tenant with id: %v", t.Id)
	return t, nil
}

func (s *Storage) StoreTenant(t social_man.Tenant) (social_man.Tenant, error) {
	logrus.Infof("storing new tenant %v, %v", t.Lastname, t.Surname)

	insertOneResult, err := s.tenants.InsertOne(context.Background(), t)

	if err != nil {
		return t, err
	}
	if oid, ok := insertOneResult.InsertedID.(objectid.ObjectID); ok {
		t.Id = oid.Hex()
	} else {
		// Not objectid.ObjectID, do what you want
	}

	logrus.Infof("stored new tenant by id: %v", t.Id)
	return t, nil
}

func (s *Storage) UpdateTenant(t social_man.Tenant) (social_man.Tenant, error) {
	logrus.Infof("starting update of tenant with id: %v", t.Id)

	result := s.tenants.FindOneAndUpdate(context.Background(), filterDocument(t.Id), t)
	document := bson.NewDocument()

	if decodeErr := result.Decode(document); decodeErr != nil {
		return t, decodeErr
	}

	if parseErr := parseFromDocument(document, &t); parseErr != nil {
		return t, parseErr
	}

	logrus.Infof("updated tenant with id: %v", t.Id)
	return t, nil
}

func filterDocument(id string) *bson.Document {
	// do not check for error - relying only internal IDs are used ;-)
	oid, _ := objectid.FromHex(id)
	return bson.NewDocument(bson.EC.ObjectID("_id", oid))
}

func parseFromDocument(document *bson.Document, t *social_man.Tenant) error {
	extJSON := document.ToExtJSON(true)
	decoder := json.NewDecoder(bytes.NewReader([]byte(extJSON)))

	decodeErr := decoder.Decode(t)
	if decodeErr != nil {
		return decodeErr
	}

	t.Id = document.Lookup("_id").ObjectID().Hex()
	return nil
}
