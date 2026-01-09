package sequence

//
//import (
//	"fmt"
//
//	"github.com/HaChains/infra/blockchain/types"
//)
//
//func (s *Sequence) decide(post types.IBlock) error {
//	iter := s.queue.Iterator()
//	iter.End()
//	for iter.Prev() {
//		height := iter.Value()
//		slot, ok := s.slots[height]
//		if !ok {
//			return fmt.Errorf("slot %d not found", height)
//		}
//		if slot.Selected().Hash() == post.ParentHash() {
//			return nil
//		}
//		ok = slot.Select(post.ParentHash())
//		// if parent block is not included, fetch it
//		if !ok {
//			ib, err := s.blockFetter.GetBlockByHash(post.ParentHash())
//			if err != nil {
//				return fmt.Errorf("parent block %s not found for block #%d", post.ParentHash(), post.Height())
//			}
//			err = slot.Insert(ib)
//			if err != nil {
//				return fmt.Errorf("failed to insert block %s into slot #%d: %v", ib.Hash(), height, err)
//			}
//			ok = slot.Select(ib.Hash())
//			// should never happen
//			if !ok {
//				return fmt.Errorf("failed to select parent block %s in slot #%d", ib.Hash(), height)
//			}
//		}
//		if height < s.low {
//			s.low = height
//		}
//		post = slot.Selected()
//		s.store.SetBlock(post)
//	}
//	return nil
//}
