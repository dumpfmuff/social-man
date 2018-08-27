package mongo

import (
	"github.com/labstack/gommon/log"
	"testing"
)

func TestSomething(t *testing.T) {



	clients := FetchAllTenants()

	log.Printf("clients: %v", clients)

}
