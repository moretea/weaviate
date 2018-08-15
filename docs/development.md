# Development of Weaviate

## FAQ
- Based on `go-swagger` tool.
- The following files are completely generated.
  - `models`/ 
  - `restapi/`
  - `restapi/server.go`
  - `cmd/weaviate-server/main.go`
- The file `restapi/configure_weaviate.go` is partially automatically generated, partially hand-edited.

## Data Model
- Weaviate stores Things, Actions and Keys.
- Keys are used both for authentication and authorization in Weaviate.
- Owners of a key can create more keys; the new key points to the parent key that is used to create the key.
- Permissions (read, write, delete, execute) are linked to a key.
- Each piece of data (e.g. Things & Actions) is associated with a Key.
