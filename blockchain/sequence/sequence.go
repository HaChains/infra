package sequence

//
//import (
//	"fmt"
//
//	"github.com/HaChains/infra/blockchain/types"
//	"github.com/emirpasic/gods/v2/queues/arrayqueue"
//)
//
//type fetcher interface {
//	GetBlockByHash(hash string) (types.IBlock, error)
//}
//
//type store interface {
//	SetBlock(block types.IBlock) error
//	DeleteBlock(height uint64) error
//}
//
//type Slot struct {
//	height   uint64
//	selected string // selected block hash
//	mBlocks  map[string]types.IBlock
//}
//
//func NewSlot(b types.IBlock) *Slot {
//	return &Slot{
//		height:   b.Height(),
//		selected: b.Hash(),
//		mBlocks:  map[string]types.IBlock{b.Hash(): b},
//	}
//}
//
//// Select selects a block by its hash. Returns false if the block does not exist in the slot.
//func (s *Slot) Select(hash string) bool {
//	if _, exists := s.mBlocks[hash]; !exists {
//		return false
//	}
//	s.selected = hash
//	return true
//}
//
//func (s *Slot) Selected() types.IBlock {
//	return s.mBlocks[s.selected]
//}
//
//func (s *Slot) Height() uint64 {
//	return s.height
//}
//
//func (s *Slot) Insert(b types.IBlock) error {
//	if b.Height() != s.height {
//		return fmt.Errorf("block height mismatch")
//	}
//	s.mBlocks[b.Hash()] = b
//	return nil
//}
//
//type Sequence struct {
//	size      int
//	slots     map[uint64]*Slot
//	queue     *arrayqueue.Queue[uint64]
//	queueTail uint64
//	low       uint64
//
//	blockFetter fetcher
//	store       store
//}
//
//// New
////
////	lastFinalized: the last finalized block to start the sequence. if nil, validation of first appended block is skipped.
//func New(confirmation int, lastFinalized types.IBlock, blockFetter fetcher, store store) (*Sequence, error) {
//	size := confirmation + 1 // first slot is kept for the last finalized block
//	slots := make(map[uint64]*Slot, size)
//	queue := arrayqueue.New[uint64]()
//	s := &Sequence{
//		size:  size,
//		slots: slots,
//		queue: queue,
//
//		blockFetter: blockFetter,
//		store:       store,
//	}
//	if lastFinalized == nil {
//		return s, nil
//	}
//	err := s.append(lastFinalized)
//	if err != nil {
//		return nil, err
//	}
//	s.low = lastFinalized.Height()
//	return s, nil
//}
//
//func (s *Sequence) Get() (seq []types.IBlock, err error) {
//	seq = make([]types.IBlock, 0, s.size)
//	for i := s.low; i <= s.queueTail; i++ {
//		slot, ok := s.slots[i]
//		// should never happen
//		if !ok {
//			return nil, fmt.Errorf("slot %d not found in sequence", i)
//		}
//		seq = append(seq, slot.Selected())
//	}
//	s.low = s.queueTail + 1
//	return seq, nil
//}
//
//func (s *Sequence) Append(b types.IBlock) error {
//	if s.queue.Size() == 0 {
//		return s.append(b)
//	}
//	err := s.check(b)
//	if err != nil {
//		return err
//	}
//	return s.append(b)
//}
//
//func (s *Sequence) check(b types.IBlock) error {
//	height := b.Height()
//	if height != s.queueTail+1 {
//		return fmt.Errorf("block height %d not consecutive to %d", height, s.queueTail)
//	}
//	return s.decide(b)
//}
//
//func (s *Sequence) append(b types.IBlock) error {
//	height := b.Height()
//	if s.slots[height] != nil {
//		return fmt.Errorf("slot %d already exists in sequence", height)
//	}
//	slot := NewSlot(b)
//	s.slots[height] = slot
//	s.queue.Enqueue(height)
//
//	err := s.store.SetBlock(b)
//	if err != nil {
//		return err
//	}
//
//	s.queueTail = height
//	for s.queue.Size() > s.size {
//		err = s.evict()
//		// should never happen
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func (s *Sequence) evict() error {
//	h, ok := s.queue.Dequeue()
//	if !ok {
//		return fmt.Errorf("failed to dequeue from sequence")
//	}
//	s.store.DeleteBlock(h)
//	delete(s.slots, h)
//	return nil
//}
