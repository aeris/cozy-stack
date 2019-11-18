package settings

import (
	"net/http"

	"github.com/cozy/cozy-stack/model/instance"
	"github.com/cozy/cozy-stack/model/permission"
	"github.com/cozy/cozy-stack/pkg/config/config"
	"github.com/cozy/cozy-stack/pkg/consts"
	"github.com/cozy/cozy-stack/pkg/couchdb"
	"github.com/cozy/cozy-stack/pkg/jsonapi"
	"github.com/cozy/cozy-stack/web/middlewares"
	"github.com/labstack/echo/v4"
)

type apiCapabilities struct {
	DocID          string `json:"_id,omitempty"`
	FileVersioning bool   `json:"file_versioning"`
}

func (c *apiCapabilities) ID() string                             { return c.DocID }
func (c *apiCapabilities) Rev() string                            { return "" }
func (c *apiCapabilities) DocType() string                        { return consts.Settings }
func (c *apiCapabilities) Clone() couchdb.Doc                     { cloned := *c; return &cloned }
func (c *apiCapabilities) SetID(id string)                        { c.DocID = id }
func (c *apiCapabilities) SetRev(rev string)                      {}
func (c *apiCapabilities) Relationships() jsonapi.RelationshipMap { return nil }
func (c *apiCapabilities) Included() []jsonapi.Object             { return nil }
func (c *apiCapabilities) Links() *jsonapi.LinksList {
	return &jsonapi.LinksList{Self: "/settings/capabilities"}
}
func (c *apiCapabilities) Fetch(field string) []string { return nil }

func newCapabilities(inst *instance.Instance) *apiCapabilities {
	// File versioning is enabled for all instances, except for the Swift
	// layout v1 and v2
	versioning := true
	switch config.FsURL().Scheme {
	case config.SchemeSwift, config.SchemeSwiftSecure:
		versioning = inst.SwiftLayout >= 2
	}
	return &apiCapabilities{
		DocID:          consts.CapabilitiesSettingsID,
		FileVersioning: versioning,
	}
}

func getCapabilities(c echo.Context) error {
	inst := middlewares.GetInstance(c)
	doc := newCapabilities(inst)
	if err := middlewares.Allow(c, permission.GET, doc); err != nil {
		return err
	}
	return jsonapi.Data(c, http.StatusOK, doc, nil)
}
