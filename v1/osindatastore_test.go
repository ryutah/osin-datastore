package datastore

import (
	"go.mercari.io/datastore"
)

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
