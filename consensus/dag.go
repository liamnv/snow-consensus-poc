package consensus

import "errors"

const decisionThreshold = 1

type Block struct {
	blockHash            string
	chit                 bool
	confidence           int
	consecutiveSuccesses int
	accepted             bool
	Ancestors            map[string]*Block
	Descendant           map[string]*Block
}

func NewBlock(blockHash string) Block {
	return Block{
		blockHash:            blockHash,
		chit:                 false,
		confidence:           0,
		consecutiveSuccesses: 0,
		accepted:             false,
		Ancestors:            make(map[string]*Block),
		Descendant:           make(map[string]*Block),
	}
}

func (b *Block) Hash() string {
	return b.blockHash
}

func (b *Block) Chit() bool {
	return b.chit
}

func (b *Block) Confidence() int {
	return b.confidence
}
func (b *Block) Accepted() bool {
	return b.accepted
}

func (b *Block) ConsecutiveSuccesses() int {
	return b.consecutiveSuccesses
}

type DAG struct {
	head           *Block
	processPointer *Block
	len            int
}

func InitialDAG() DAG {
	return DAG{
		head:           nil,
		processPointer: nil,
	}
}

func (d *DAG) AddBlock(block *Block, ancestorHashes []string) error {
	if d.head == nil && len(ancestorHashes) != 0 {
		return errors.New("invalid ancestor hashes, DAG with nil HEAD can only add new block to head")
	}
	if d.head == nil {
		d.head = block
		return nil
	}
	ancestors := d.search(ancestorHashes)
	if len(ancestors) != len(ancestorHashes) {
		return errors.New("can't find all ancestors with hashes")
	}
	for _, ancestor := range ancestors {
		ancestor.Descendant[block.Hash()] = block
		block.Ancestors[ancestor.Hash()] = ancestor
	}
	return nil
}

func (d *DAG) SuccessPool(result string, options []string) {
	blocks := d.search(options)
	for _, block := range blocks {
		ancestors := block.traverse(true)
		// if block is not prefered by almost of node
		if block.Hash() != result {
			block.chit = false
			block.confidence = 0
			block.consecutiveSuccesses = 0
			//update ancestors consecutive successes to 0
			for _, ancestor := range ancestors {
				if !ancestor.accepted {
					ancestor.consecutiveSuccesses = 0
				}
			}
			continue
		}
		// increase ancestors consecutive successes and confidence
		for _, ancestor := range ancestors {
			if ancestor.accepted {
				continue
			}
			ancestor.consecutiveSuccesses++
			ancestor.confidence++
			ancestor.chit = true
			ancestor.accepted = ancestor.consecutiveSuccesses > decisionThreshold
			if ancestor.accepted {
				d.head = ancestor
				d.head.Ancestors = make(map[string]*Block)
			}
		}

	}
}

func (d *DAG) search(blockHashes []string) []*Block {
	visited := make(map[string]bool, 0)
	queue := NewQueue()
	if d.head == nil {
		return make([]*Block, 0)
	}
	result := make([]*Block, 0)
	queue.Add(d.head)
	visited[d.head.Hash()] = true
	for queue.len != 0 {
		node := queue.Head()
		queue.Pop()

		// add all descendants to queue
		for i := range node.Value.Descendant {
			if !visited[node.Value.Descendant[i].Hash()] {
				queue.Add(node.Value.Descendant[i])
				visited[node.Value.Descendant[i].Hash()] = true
			}
		}

		//check if node have blockHash in list
		for _, hash := range blockHashes {
			if node.Value.Hash() == hash {
				result = append(result, node.Value)
			}
		}
	}
	return result
}

// traverse from block, reserve = false -> traverse from ancestor to descendant, reserve = true -> traverse from descendant
// to ancestors
func (b *Block) traverse(reverse bool) []*Block {
	visited := make(map[string]bool, 0)
	queue := NewQueue()
	if b == nil {
		return make([]*Block, 0)
	}
	result := make([]*Block, 0)
	pointer := b
	queue.Add(pointer)
	visited[pointer.Hash()] = true
	for queue.len != 0 {
		node := queue.Head()
		queue.Pop()

		// add all descendants to queue
		if reverse {
			for i := range node.Value.Ancestors {
				if !visited[node.Value.Ancestors[i].Hash()] {
					queue.Add(node.Value.Ancestors[i])
					visited[node.Value.Ancestors[i].Hash()] = true
				}
			}
		} else {
			for i := range node.Value.Descendant {
				if !visited[node.Value.Descendant[i].Hash()] {
					queue.Add(node.Value.Descendant[i])
					visited[node.Value.Descendant[i].Hash()] = true
				}
			}
		}

		result = append(result, node.Value)
	}
	return result
}
