package consensus

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

/*
                     A--------
                   /   \     |
                  B     C    |
                   \   /     |
	                 D       |
	                         E
*/
func TestSearch(t *testing.T) {
	blockA := NewBlock("blockHashA")
	blockB := NewBlock("blockHashB")
	blockC := NewBlock("blockHashC")
	blockD := NewBlock("blockHashD")
	blockE := NewBlock("blockHashE")
	blockA.Descendant = map[string]*Block{
		"blockHashB": &blockB,
		"blockHashC": &blockC,
		"blockHashD": &blockD,
		"blockHashE": &blockE,
	}
	blockB.Descendant = map[string]*Block{
		"blockHashD": &blockD,
	}

	blockC.Descendant = map[string]*Block{
		"blockHashD": &blockD,
	}
	blockB.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
	}
	blockC.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
	}
	blockD.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
		"blockHashB": &blockB,
		"blockHashC": &blockC,
	}
	blockE.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
	}
	dag := InitialDAG()
	dag.head = &blockA

	blocks := dag.search([]string{"blockHashA"})
	assert.Equal(t, []*Block{&blockA}, blocks)
}

/*
                     A------
                   /   \    |
                  B     C   |
                   \   /    |
	                 D      |
	                   \----E


*/
func TestSearch2(t *testing.T) {
	blockA := NewBlock("blockHashA")
	blockB := NewBlock("blockHashB")
	blockC := NewBlock("blockHashC")
	blockD := NewBlock("blockHashD")
	blockE := NewBlock("blockHashE")
	blockA.Descendant = map[string]*Block{
		"blockHashB": &blockB,
		"blockHashC": &blockC,
		"blockHashD": &blockD,
		"blockHashE": &blockE,
	}
	blockB.Descendant = map[string]*Block{
		"blockHashD": &blockD,
	}

	blockC.Descendant = map[string]*Block{
		"blockHashD": &blockD,
	}

	blockD.Descendant = map[string]*Block{
		"blockHashE": &blockE,
	}
	blockB.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
	}
	blockC.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
	}
	blockD.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
		"blockHashB": &blockB,
		"blockHashC": &blockC,
	}
	blockE.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
		"blockHashD": &blockD,
	}
	dag := InitialDAG()
	dag.head = &blockA

	blocks := dag.search([]string{"blockHashA", "blockHashD"})
	assert.Equal(t, []*Block{&blockA, &blockD}, blocks)
}

func TestBlockTraverse(t *testing.T) {
	blockA := NewBlock("blockHashA")
	blockB := NewBlock("blockHashB")
	blockC := NewBlock("blockHashC")
	blockD := NewBlock("blockHashD")
	blockE := NewBlock("blockHashE")
	blockA.Descendant = map[string]*Block{
		"blockHashB": &blockB,
		"blockHashC": &blockC,
		"blockHashD": &blockD,
		"blockHashE": &blockE,
	}
	blockB.Descendant = map[string]*Block{
		"blockHashD": &blockD,
	}

	blockC.Descendant = map[string]*Block{
		"blockHashD": &blockD,
	}

	blockD.Descendant = map[string]*Block{
		"blockHashE": &blockE,
	}
	blockB.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
	}
	blockC.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
	}
	blockD.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
		"blockHashB": &blockB,
		"blockHashC": &blockC,
	}
	blockE.Ancestors = map[string]*Block{
		"blockHashA": &blockA,
		"blockHashD": &blockD,
	}

	blocks := blockA.traverse(false)
	for _, block := range blocks {
		fmt.Println(block.Hash())
	}

	reverseBlocks := blockE.traverse(true)
	for _, block := range reverseBlocks {
		fmt.Println(block.Hash())
	}
}

func TestDAG_AddBlock_InvalidAncestorsToNilDAG(t *testing.T) {
	dag := InitialDAG()
	blockA := NewBlock("blockHashA")
	err := dag.AddBlock(&blockA, []string{"blockB"})
	assert.NotNil(t, err)
	assert.Equal(t, errors.New("invalid ancestor hashes, DAG with nil HEAD can only add new block to head"), err)
}

func TestDAG_AddBlock_ToHead(t *testing.T) {
	dag := InitialDAG()
	blockA := NewBlock("blockHashA")
	err := dag.AddBlock(&blockA, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "blockHashA", dag.head.Hash())
}

/*
                     A
                   / | \
                  B  |  C
                   \ | /
	                 D

*/
func TestDAG_AddBlock(t *testing.T) {
	dag := InitialDAG()
	blockA := NewBlock("blockHashA")
	err := dag.AddBlock(&blockA, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "blockHashA", dag.head.Hash())

	blockB := NewBlock("blockHashB")
	err = dag.AddBlock(&blockB, []string{"blockHashA"})
	assert.Nil(t, err)
	assert.Equal(t, map[string]*Block{
		"blockHashB": &blockB,
	}, dag.head.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashA": &blockA,
	}, blockB.Ancestors)

	blockC := NewBlock("blockHashC")

	err = dag.AddBlock(&blockC, []string{"blockHashA"})
	assert.Nil(t, err)
	assert.Equal(t, map[string]*Block{
		"blockHashB": &blockB,
		"blockHashC": &blockC,
	}, dag.head.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashA": &blockA,
	}, blockC.Ancestors)

	blockD := NewBlock("blockHashD")
	err = dag.AddBlock(&blockD, []string{"blockHashA", "blockHashB", "blockHashC"})

	assert.Nil(t, err)
	assert.Equal(t, map[string]*Block{
		"blockHashB": &blockB,
		"blockHashC": &blockC,
		"blockHashD": &blockD,
	}, dag.head.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashD": &blockD,
	}, blockB.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashD": &blockD,
	}, blockC.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashA": &blockA,
		"blockHashB": &blockB,
		"blockHashC": &blockC,
	}, blockD.Ancestors)
}

/*
		 A
	   /   \
	   B    C
	   \   / \
		 Y    Z

*/
func TestDAG_SuccessPool(t *testing.T) {

	dag := InitialDAG()
	blockA := NewBlock("blockHashA")
	blockA.accepted = true
	err := dag.AddBlock(&blockA, []string{})
	assert.Nil(t, err)
	assert.Equal(t, "blockHashA", dag.head.Hash())

	blockB := NewBlock("blockHashB")
	err = dag.AddBlock(&blockB, []string{"blockHashA"})
	assert.Nil(t, err)
	assert.Equal(t, map[string]*Block{
		"blockHashB": &blockB,
	}, dag.head.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashA": &blockA,
	}, blockB.Ancestors)

	blockC := NewBlock("blockHashC")

	err = dag.AddBlock(&blockC, []string{"blockHashA"})
	assert.Nil(t, err)
	assert.Equal(t, map[string]*Block{
		"blockHashB": &blockB,
		"blockHashC": &blockC,
	}, dag.head.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashA": &blockA,
	}, blockC.Ancestors)

	dag.SuccessPool("blockHashB", []string{"blockHashB", "blockHashC"})
	assert.Equal(t, true, blockB.chit)
	assert.Equal(t, 1, blockB.Confidence())
	assert.Equal(t, 1, blockB.ConsecutiveSuccesses())

	assert.Equal(t, false, blockC.chit)
	assert.Equal(t, 0, blockC.Confidence())
	assert.Equal(t, 0, blockC.ConsecutiveSuccesses())

	blockY := NewBlock("blockHashY")
	err = dag.AddBlock(&blockY, []string{"blockHashB", "blockHashC"})
	assert.Nil(t, err)
	assert.Equal(t, map[string]*Block{
		"blockHashY": &blockY,
	}, blockB.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashY": &blockY,
	}, blockC.Descendant)

	assert.Equal(t, map[string]*Block{
		"blockHashB": &blockB,
		"blockHashC": &blockC,
	}, blockY.Ancestors)

	blockZ := NewBlock("blockHashZ")
	err = dag.AddBlock(&blockZ, []string{"blockHashC"})
	assert.Equal(t, map[string]*Block{
		"blockHashY": &blockY,
		"blockHashZ": &blockZ,
	}, blockC.Descendant)
	assert.Equal(t, map[string]*Block{
		"blockHashC": &blockC,
	}, blockZ.Ancestors)

	dag.SuccessPool("blockHashY", []string{"blockHashY", "blockHashZ"})
	assert.Equal(t, true, blockY.Chit())
	assert.Equal(t, 1, blockY.Confidence())
	assert.Equal(t, 1, blockY.ConsecutiveSuccesses())

	assert.Equal(t, false, blockZ.Chit())
	assert.Equal(t, 0, blockZ.Confidence())
	assert.Equal(t, 0, blockZ.ConsecutiveSuccesses())

	assert.Equal(t, 2, blockB.Confidence())
	assert.Equal(t, 2, blockB.ConsecutiveSuccesses())
	assert.Equal(t, true, blockB.accepted)

	assert.Equal(t, 1, blockC.Confidence())
	assert.Equal(t, 0, blockC.ConsecutiveSuccesses())
}
