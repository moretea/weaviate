/*                          _       _
 *__      _____  __ ___   ___  __ _| |_ ___
 *\ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
 * \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
 *  \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
 *
 * Copyright © 2016 Weaviate. All rights reserved.
 * LICENSE: https://github.com/weaviate/weaviate/blob/master/LICENSE
 * AUTHOR: Bob van Luijt (bob@weaviate.com)
 * See www.weaviate.com for details
 * See package.json for author and maintainer info
 * Contact: @weaviate_iot / yourfriends@weaviate.com
 */
 package adapters




import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// WeaveAdaptersDeactivateHandlerFunc turns a function with the right signature into a weave adapters deactivate handler
type WeaveAdaptersDeactivateHandlerFunc func(WeaveAdaptersDeactivateParams) middleware.Responder

// Handle executing the request and returning a response
func (fn WeaveAdaptersDeactivateHandlerFunc) Handle(params WeaveAdaptersDeactivateParams) middleware.Responder {
	return fn(params)
}

// WeaveAdaptersDeactivateHandler interface for that can handle valid weave adapters deactivate params
type WeaveAdaptersDeactivateHandler interface {
	Handle(WeaveAdaptersDeactivateParams) middleware.Responder
}

// NewWeaveAdaptersDeactivate creates a new http.Handler for the weave adapters deactivate operation
func NewWeaveAdaptersDeactivate(ctx *middleware.Context, handler WeaveAdaptersDeactivateHandler) *WeaveAdaptersDeactivate {
	return &WeaveAdaptersDeactivate{Context: ctx, Handler: handler}
}

/*WeaveAdaptersDeactivate swagger:route POST /adapters/{adapterId}/deactivate adapters weaveAdaptersDeactivate

Deactivates an adapter. This will also delete all devices provided by that adapter.

*/
type WeaveAdaptersDeactivate struct {
	Context *middleware.Context
	Handler WeaveAdaptersDeactivateHandler
}

func (o *WeaveAdaptersDeactivate) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewWeaveAdaptersDeactivateParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}