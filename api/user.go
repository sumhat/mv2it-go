package api

type GroupEntry struct {
	Id int `datastore:"-"`
	Name string `datastore:"name,noindex"`
	TimeoutSeconds int64 `datastore:"TimeoutSeconds,noindex"`
}

type UserEntry struct {
	Id string `datastore:"-"`

}