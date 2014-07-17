package api

const (
	DefaultUser = "anonymous"
	Namespace = "core"
)

type GroupEntry struct {
	Id int `datastore:"-"`
	Name string `datastore:"name,noindex"`
	TimeoutSeconds int64 `datastore:"TimeoutSeconds,noindex"`
}

type UserEntry struct {
	Id string `datastore:"-"`

}