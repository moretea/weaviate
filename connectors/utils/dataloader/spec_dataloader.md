# Technical specification of the dataloader
This is Weaviate's implementation of the [graphql-gophers dataloader](https://github.com/graph-gophers/dataloader). This allows the use of batching and caching, reducing the amount of traffic to and from the database.

## What it does
The dataloader exists as middleware between Weaviate's GraphQL [endpoint](https://github.com/creativesoftwarefdn/weaviate/tree/develop/graphqlapi) and the database [connectors](https://github.com/creativesoftwarefdn/weaviate/tree/develop/connectors), as illustrated in this [image](https://github.com/creativesoftwarefdn/weaviate/blob/develop/docs/data-flow.png). It implements the [BaseConnector](https://github.com/creativesoftwarefdn/weaviate/blob/1d4b50f8f2a9fb47615c0c29da6540de1f386a0e/connectors/database_connector.go#L30) interface, causing Weaviate to communicate with it as if it is a database connector. The dataloader applies batching and caching where applicable, and then passes the requests on to the actual database connector.

## Function definitions: general
Note that only functions that do more than directly forward the request to the database connector are listed here

### Attach()
```go
Attach(ctx context.Context) (context.Context, error)
```
**Goal**: 
* instantiate the dataloaders for Things, Actions and Keys
* ensure that each request has its own dataloader

**Params**:
* `ctx context.Context`: a global variable for the request

**Returns**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `err error`: nil if everything goes correctly

***

## Function definitions: GraphQL

### GetLocalGraph()
```go
GetLocalGraph(request graphql.RequestParams) (interface{}, error)
```
**Goal**: 
* act as a central access point for all GraphQL endpoint resolvers
* interpret incoming request types by checking the value of `request.Args['weaviate_request_type']` 
* forward requests to the functions matching their types

**Params**:
* `request graphql.RequestParams`: a struct storing all information relevant to the request

**Returns**:
* `TODO`: to be determined

***
### resolveLocalThing()
```go
resolveLocalThing(request graphql.ResolveParams) (interface{}, error)
```
**Goal**: 
* handle both LocalGetThing and LocalGetMetaThing requests
* apply batching and caching
* forward the Thing request to the database connector

**Params**:
* `request graphql.ResolveParams`: a struct storing all information relevant to the request. `request.Args['weaviate_request_type']` contains the request type (e.g. `'local_get_thing'` or `'local_get_meta_action'`)

**Returns**:
* `TODO`: to be determined

***
### resolveLocalAction()
```go
resolveLocalAction(request graphql.ResolveParams) (interface{}, error)
```
**Goal**: 
* handle both LocalGetAction and LocalGetMetaAction requests
* apply batching and caching
* forward the Action request to the database connector

**Params**:
* `request graphql.ResolveParams`: a struct storing all information relevant to the request. `request.Args['weaviate_request_type']` contains the request type (e.g. `'local_get_thing'` or `'local_get_meta_action'`)

**Returns**:
* `TODO`: to be determined

***

## Function definitions: REST
Note that only functions that do more than directly forward the request to the database connector are listed here

### GetThing()
```go
GetThing(ctx context.Context, UUID strfmt.UUID, thingResponse *models.ThingGetResponse) error
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `UUID strfmt.UUID`: the uuid of the Thing to get
* `thingResponse *models.ThingGetResponse`: a class representing the thing object response to the request

**Returns**:
* `err error`: nil if everything goes correctly

***
### GetAction()
```go
GetAction(ctx context.Context, UUID strfmt.UUID, actionResponse *models.ActionGetResponse) error
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `UUID strfmt.UUID`: the uuid of the Action to get
* `actionResponse *models.ActionGetResponse`: a class representing the action object response to the request

**Returns**:
* `err error`: nil if everything goes correctly

***
### GetKey()
```go
GetKey(ctx context.Context, UUID strfmt.UUID, KeyResponse *models.KeyGetResponse) error
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `UUID strfmt.UUID`: the uuid of the Key to get
* `keyResponse *models.KeyGetResponse`: a class representing the key object response to the request

**Returns**:
* `err error`: nil if everything goes correctly

***
### UpdateThing()
```go
UpdateThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `thing *models.Thing`: a struct representing a Thing object
* `UUID strfmt.UUID`: the uuid of the Thing to update

**Returns**:
* `err error`: nil if everything goes correctly

***
### UpdateAction()
```go
UpdateAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `action *models.Action`: a struct representing an Action object
* `UUID strfmt.UUID`: the uuid of the Action to update

**Returns**:
* `err error`: nil if everything goes correctly

***
### UpdateKey()
```go
UpdateKey(ctx context.Context, key *models.Thing, UUID strfmt.UUID) error {
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `key *models.Key`: a struct representing a Key object
* `UUID strfmt.UUID`: the uuid of the Key to update

**Returns**:
* `err error`: nil if everything goes correctly

***
### DeleteThing()
```go
DeleteThing(ctx context.Context, thing *models.Thing, UUID strfmt.UUID) error {
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `thing *models.Thing`: a struct representing a Thing object
* `UUID strfmt.UUID`: the uuid of the Thing to delete

**Returns**:
* `err error`: nil if everything goes correctly

***
### DeleteAction()
```go
DeleteAction(ctx context.Context, action *models.Action, UUID strfmt.UUID) error {
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `action *models.Action`: a struct representing an Action object
* `UUID strfmt.UUID`: the uuid of the Action to delete

**Returns**:
* `err error`: nil if everything goes correctly

***
### DeleteKey()
```go
DeleteKey(ctx context.Context, key *models.Key, UUID strfmt.UUID) error {
```
**Goal**: 
* apply batching and caching
* forward REST update request to database connector

**Params**:
* `ctx context.Context`: the context with the added instantiated dataloaders, nil if an error occurs
* `key *models.Key`: a struct representing an Key object
* `UUID strfmt.UUID`: the uuid of the Key to delete

**Returns**:
* `err error`: nil if everything goes correctly

***
