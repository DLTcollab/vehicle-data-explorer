package kvstore

type KVStore interface {
	Set(key string, value string) error
	Get(key string) (value string, err error)
	Has(key string) (bool, error)
	Delete(key string) error
}
