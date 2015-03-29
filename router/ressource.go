package router

import (
	"net/http"
	"path"
)

// Ressource is a collection of http.HandlerFuncs for different http methods that will be mounted
// under a common router and prefix
type Ressource struct {
	// handlers for http methods that handle a single item
	GET, POST, PATCH, PUT, DELETE http.HandlerFunc
	// handler for http GET that returns a list of items
	INDEX http.HandlerFunc

	itemsURL func() string
	itemURL  func(itemID string) string
	mounted  bool
}

// IndexURL returns the URL that will invoke the SEARCH handler
func (rs *Ressource) IndexURL() string {
	return rs.itemsURL()
}

// PostURL returns the URL that will invoke the POST handler
func (rs *Ressource) PostURL() string {
	return rs.itemsURL()
}

// GetURL returns the URL that will invoke the GET handler
func (rs *Ressource) GetURL(itemID string) string {
	return rs.itemURL(itemID)
}

// DeleteURL returns the URL that will invoke the DELETE handler
func (rs *Ressource) DeleteURL(itemID string) string {
	return rs.itemURL(itemID)
}

// PatchURL returns the URL that will invoke the PATCH handler
func (rs *Ressource) PatchURL(itemID string) string {
	return rs.itemURL(itemID)
}

// PutURL returns the URL that will invoke the PUT handler
func (rs *Ressource) PutURL(itemID string) string {
	return rs.itemURL(itemID)
}

// Mount mounts the ressource http.Handlers under the given path segment inside the router
func (rs *Ressource) Mount(pathSegment string, rt *Router) {

	if rs.mounted {
		panic("ressource cant be mounted twice")
	}

	rs.mounted = true

	if rs.GET != nil {
		rt.GETFunc(pathSegment, rs.GET)
	}

	if rs.INDEX != nil {
		rt.INDEXFunc(pathSegment, rs.INDEX)
	}

	if rs.POST != nil {
		rt.POSTFunc(pathSegment, rs.POST)
	}

	if rs.PUT != nil {
		rt.PUTFunc(pathSegment, rs.PUT)
	}

	if rs.PATCH != nil {
		rt.PATCHFunc(pathSegment, rs.PATCH)
	}

	if rs.DELETE != nil {
		rt.DELETEFunc(pathSegment, rs.DELETE)
	}

	rs.itemsURL = func() string { return path.Join(rt.MountPath(), pathSegment) }
	rs.itemURL = func(param string) string { return path.Join(rt.MountPath(), pathSegment, param) }
}
