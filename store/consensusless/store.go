package consensusless

import (
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
	s.DB.Set(key, value)
}

func (s Store) Delete(key []byte) {
	s.DB.Delete(key)
}

func (s Store) Write() {
	fmt.Println("consensusless storage Write() not implemented")
}

// Commit performs a no-op as entries are persistent between commitments.
func (s *Store) Commit() (id types.CommitID) {
	batch := s.Store.NewBatch()
	defer func() {
		_ = batch.Close()
	}()

	if err := batch.WriteSync(); err != nil {
		panic(fmt.Errorf("error on batch write %w", err))
	}
	return
}

func (s Store) Reset() error {
	var err error
	batch := s.NewBatch()
	defer func() {
		_ = batch.Close()
	}()

	itr := s.Iterator(nil, nil)
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		batch.Delete(itr.Key())
	}
	err = batch.WriteSync()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) SetPruning(pruning pruningtypes.PruningOptions) {}

// GetPruning is a no-op as pruning options cannot be directly set on this store.
// They must be set on the root commit multi-store.
func (s *Store) GetPruning() pruningtypes.PruningOptions {
	return pruningtypes.NewPruningOptions(pruningtypes.PruningUndefined)
}

func (s Store) LastCommitID() (id types.CommitID) {
	// fmt.Println("consensusless storage LastCommitID() not implemented")
	return
}

func (s Store) WorkingHash() (hash []byte) {
	// fmt.Println("consensusless storage WorkingHash() not implemented")
	return make([]byte, 0)
}