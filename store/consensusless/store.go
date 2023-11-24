package consensusless

import (
	"cosmossdk.io/store/iavl"
	"cosmossdk.io/store/types"
)

// Store implements an in-memory only KVStore. Entries are persisted between
// commits and thus between blocks. State in Memory store is not committed as part of app state but maintained privately by each node
type Store struct {
	*iavl.Store
}

// GetStoreType returns the Store's type.
func (s Store) GetStoreType() types.StoreType {
	return types.StoreTypeConsensusless
}
