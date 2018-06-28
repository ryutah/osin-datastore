package datastore

import (
	"context"

	"go.mercari.io/datastore"
)

type mockDatastoreClient struct {
	get           func(ctx context.Context, key datastore.Key, dst interface{}) error
	getMulti      func(ctx context.Context, keys []datastore.Key, dst interface{}) error
	put           func(ctx context.Context, key datastore.Key, src interface{}) (datastore.Key, error)
	delete        func(ctx context.Context, key datastore.Key) error
	incompleteKey func(kind string, parent datastore.Key) datastore.Key
	idKey         func(kind string, id int64, parent datastore.Key) datastore.Key
	nameKey       func(kind string, name string, parent datastore.Key) datastore.Key
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
	return m.delete(ctx, key)
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
	return m.nameKey(kind, name, parent)
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
	return m.name
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

type mockClientHandler struct {
	_put    []func(ctx context.Context, c *Client) error
	_get    []func(ctx context.Context, id string) (*Client, error)
	_delete []func(ctx context.Context, id string) error
}

func (m *mockClientHandler) put(ctx context.Context, c *Client) error {
	mock := m._put[0]
	if len(m._put) > 1 {
		m._put = m._put[1:]
	}
	return mock(ctx, c)
}

func (m *mockClientHandler) get(ctx context.Context, id string) (*Client, error) {
	mock := m._get[0]
	if len(m._get) > 1 {
		m._get = m._get[1:]
	}
	return mock(ctx, id)
}

func (m *mockClientHandler) delete(ctx context.Context, id string) error {
	mock := m._delete[0]
	if len(m._delete) > 1 {
		m._delete = m._delete[1:]
	}
	return mock(ctx, id)
}

type mockAuthDataHandler struct {
	_put    []func(ctx context.Context, auth *authorizeData) error
	_get    []func(ctx context.Context, code string) (*authorizeData, error)
	_delete []func(ctx context.Context, code string) error
}

func (m *mockAuthDataHandler) put(ctx context.Context, auth *authorizeData) error {
	mock := m._put[0]
	if len(m._put) > 1 {
		m._put = m._put[1:]
	}
	return mock(ctx, auth)
}

func (m *mockAuthDataHandler) get(ctx context.Context, code string) (*authorizeData, error) {
	mock := m._get[0]
	if len(m._get) > 1 {
		m._get = m._get[1:]
	}
	return mock(ctx, code)
}

func (m *mockAuthDataHandler) delete(ctx context.Context, code string) error {
	mock := m._delete[0]
	if len(m._delete) > 1 {
		m._delete = m._delete[1:]
	}
	return mock(ctx, code)
}

type mockAccessDataHandler struct {
	_put    []func(ctx context.Context, ac *accessData) error
	_get    []func(ctx context.Context, token string) (*accessData, error)
	_delete []func(ctx context.Context, token string) error
}

func (m *mockAccessDataHandler) put(ctx context.Context, ac *accessData) error {
	mock := m._put[0]
	if len(m._put) > 1 {
		m._put = m._put[1:]
	}
	return mock(ctx, ac)
}

func (m *mockAccessDataHandler) get(ctx context.Context, token string) (*accessData, error) {
	mock := m._get[0]
	if len(m._get) > 1 {
		m._get = m._get[1:]
	}
	return mock(ctx, token)
}

func (m *mockAccessDataHandler) delete(ctx context.Context, token string) error {
	mock := m._delete[0]
	if len(m._delete) > 1 {
		m._delete = m._delete[1:]
	}
	return mock(ctx, token)
}

type mockRefreshHandler struct {
	_put    []func(ctx context.Context, ref *refresh) error
	_get    []func(ctx context.Context, token string) (*refresh, error)
	_delete []func(ctx context.Context, token string) error
}

func (m *mockRefreshHandler) put(ctx context.Context, ref *refresh) error {
	mock := m._put[0]
	if len(m._put) > 1 {
		m._put = m._put[1:]
	}
	return mock(ctx, ref)
}

func (m *mockRefreshHandler) get(ctx context.Context, token string) (*refresh, error) {
	mock := m._get[0]
	if len(m._get) > 1 {
		m._get = m._get[1:]
	}
	return mock(ctx, token)
}

func (m *mockRefreshHandler) delete(ctx context.Context, token string) error {
	mock := m._delete[0]
	if len(m._delete) > 1 {
		m._delete = m._delete[1:]
	}
	return mock(ctx, token)
}
