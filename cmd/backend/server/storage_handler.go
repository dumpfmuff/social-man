package server

import (
	"encoding/json"
	"github.com/osimono/social-man"
	"github.com/osimono/social-man/cmd/backend/webutils"
	"net/http"
)

type TenantStorage interface {
	FetchAllTenants() ([]social_man.Tenant, error)
	StoreTenant(c social_man.Tenant) (social_man.Tenant, error)
}

type AllTenantHandler struct {
	storage TenantStorage
}

func (h *AllTenantHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tenants, err := h.storage.FetchAllTenants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, _ := json.Marshal(tenants)
	w.Write(b)
}

type NewTenantHandler struct {
	storage TenantStorage
}

func (h *NewTenantHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var c social_man.Tenant
	if err := webutils.ParseJson(w, r.Body, &c); err != nil {
		return
	}

	clientWithId, err := h.storage.StoreTenant(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(clientWithId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
