package datastore

import (
	"context"

	"go.mercari.io/datastore"
)

type mockDatastoreClient struct {
	get           func(ctx context.Context, key datastore.Key, dst interface{}) error
	getMulti      func(ctx context.Context, keys []datastore.Key, dst interface{}) error
	put           func(ctx context.Context, key datastore.Key, src interface{}) (datastore.Key, error)
	incompleteKey func(kind string, parent datastore.Key) datastore.Key
	idKey         func(kind string, id int64, parent datastore.Key) datastore.Key
}

func (m *mockDatastoreClient) Get(ctx context.Context, key datastore.Key, dst interface{}) error {
	return m.get(ctx, key, dst)
}

func (m *mockDatastoreClient) GetMulti(ctx context.Context, keys []datastore.Key, dst interface{}) error {
	return m.getMulti(ctx, keys, dst)
}

func (m *mockDatastoreClient) Put(ctx context.Context, key datastore.Key, src interface{}) (datastore.Key, error) {
	return m.put(ctx, key, src)
}

func (m *mockDatastoreClient) PutMulti(ctx context.Context, keys []datastore.Key, src interface{}) ([]datastore.Key, error) {
	panic("not implemented")
}

func (m *mockDatastoreClient) Delete(ctx context.Context, key datastore.Key) error {
	panic("not implemented")
}

func (m *mockDatastoreClient) DeleteMulti(ctx context.Context, keys []datastore.Key) error {
	panic("not implemented")
}

func (m *mockDatastoreClient) NewTransaction(ctx context.Context) (datastore.Transaction, error) {
	panic("not implemented")
}

func (m *mockDatastoreClient) RunInTransaction(ctx context.Context, f func(tx datastore.Transaction) error) (datastore.Commit, error) {
	panic("not implemented")
}

func (m *mockDatastoreClient) Run(ctx context.Context, q datastore.Query) datastore.Iterator {
	panic("not implemented")
}

func (m *mockDatastoreClient) AllocateIDs(ctx context.Context, keys []datastore.Key) ([]datastore.Key, error) {
	panic("not implemented")
}

func (m *mockDatastoreClient) Count(ctx context.Context, q datastore.Query) (int, error) {
	panic("not implemented")
}

func (m *mockDatastoreClient) GetAll(ctx context.Context, q datastore.Query, dst interface{}) ([]datastore.Key, error) {
	panic("not implemented")
}

func (m *mockDatastoreClient) IncompleteKey(kind string, parent datastore.Key) datastore.Key {
	return m.incompleteKey(kind, parent)
}

func (m *mockDatastoreClient) NameKey(kind string, name string, parent datastore.Key) datastore.Key {
	panic("not implemented")
}

func (m *mockDatastoreClient) IDKey(kind string, id int64, parent datastore.Key) datastore.Key {
	return m.idKey(kind, id, parent)
}

func (m *mockDatastoreClient) NewQuery(kind string) datastore.Query {
	panic("not implemented")
}

func (m *mockDatastoreClient) Close() error {
	panic("not implemented")
}

func (m *mockDatastoreClient) DecodeKey(encoded string) (datastore.Key, error) {
	panic("not implemented")
}

func (m *mockDatastoreClient) DecodeCursor(s string) (datastore.Cursor, error) {
	panic("not implemented")
}

func (m *mockDatastoreClient) Batch() *datastore.Batch {
	panic("not implemented")
}

func (m *mockDatastoreClient) AppendMiddleware(middleware datastore.Middleware) {
	panic("not implemented")
}

func (m *mockDatastoreClient) RemoveMiddleware(middleware datastore.Middleware) bool {
	panic("not implemented")
}

func (m *mockDatastoreClient) Context() context.Context {
	panic("not implemented")
}

func (m *mockDatastoreClient) SetContext(ctx context.Context) {
	panic("not implemented")
}

type mockKey struct {
	kind      string
	id        int64
	name      string
	parent    datastore.Key
	namespace string
}

func (m *mockKey) Kind() string {
	panic("not implemented")
}

func (m *mockKey) ID() int64 {
	return m.id
}

func (m *mockKey) Name() string {
	panic("not implemented")
}

func (m *mockKey) ParentKey() datastore.Key {
	panic("not implemented")
}

func (m *mockKey) Namespace() string {
	panic("not implemented")
}

func (m *mockKey) SetNamespace(namespace string) {
	panic("not implemented")
}

func (m *mockKey) String() string {
	panic("not implemented")
}

func (m *mockKey) GobEncode() ([]byte, error) {
	panic("not implemented")
}

func (m *mockKey) GobDecode(buf []byte) error {
	panic("not implemented")
}

func (m *mockKey) MarshalJSON() ([]byte, error) {
	panic("not implemented")
}

func (m *mockKey) UnmarshalJSON(buf []byte) error {
	panic("not implemented")
}

func (m *mockKey) Encode() string {
	panic("not implemented")
}

func (m *mockKey) Equal(o datastore.Key) bool {
	panic("not implemented")
}

func (m *mockKey) Incomplete() bool {
	panic("not implemented")
}
