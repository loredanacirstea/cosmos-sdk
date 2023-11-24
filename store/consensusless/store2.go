package consensusless

import (
	"encoding/hex"
	"fmt"
	"io"

	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/store/cachekv"
	"cosmossdk.io/store/dbadapter"
	pruningtypes "cosmossdk.io/store/pruning/types"
	"cosmossdk.io/store/tracekv"
	"cosmossdk.io/store/types"
)

var (
	_ types.KVStore   = (*Store)(nil)
	_ types.Committer = (*Store)(nil)
)

// Store implements an in-memory only KVStore. Entries are persisted between
// commits and thus between blocks. State in Memory store is not committed as part of app state but maintained privately by each node
type Store struct {
	dbadapter.Store
}

func NewStoreWithDB(db *dbm.DB) *Store { //nolint: interfacer // Concrete return type is fine here.
	return &Store{Store: dbadapter.Store{DB: *db}}
}

// GetStoreType returns the Store's type.
func (s Store) GetStoreType() types.StoreType {
	return types.StoreTypeConsensusless
}

// CacheWrap branches the underlying store.
func (s Store) CacheWrap() types.CacheWrap {
	return cachekv.NewStore(s)
}

// CacheWrapWithTrace implements KVStore.
func (s Store) CacheWrapWithTrace(w io.Writer, tc types.TraceContext) types.CacheWrap {
	return cachekv.NewStore(tracekv.NewStore(s, w, tc))
}

func (s Store) Set(key, value []byte) {
	fmt.Println("==STORE=Set=", hex.EncodeToString(key), hex.EncodeToString(value))
	s.DB.Set(key, value)
}

// Commit performs a no-op as entries are persistent between commitments.
func (s *Store) Commit() (id types.CommitID) {
	fmt.Println("==STORE=Commit=")
	batch := s.Store.NewBatch()
	defer func() {
		_ = batch.Close()
	}()

	if err := batch.WriteSync(); err != nil {
		panic(fmt.Errorf("error on batch write %w", err))
	}
	return
}

func (s *Store) SetPruning(pruning pruningtypes.PruningOptions) {}

// GetPruning is a no-op as pruning options cannot be directly set on this store.
// They must be set on the root commit multi-store.
func (s *Store) GetPruning() pruningtypes.PruningOptions {
	return pruningtypes.NewPruningOptions(pruningtypes.PruningUndefined)
}

func (s Store) LastCommitID() (id types.CommitID) { return }

func (s Store) WorkingHash() (hash []byte) { return }
