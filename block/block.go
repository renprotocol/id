package block

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/renproject/hyperdrive/sig"
	"github.com/renproject/hyperdrive/tx"
)

// The Round in which a `Block` was proposed.
type Round int64

// The Height at which a `Block` was proposed.
type Height int64

type Block struct {
	Time         time.Time
	Round        Round
	Height       Height
	Header       sig.Hash
	ParentHeader sig.Hash
	Txs          tx.Transactions
}

func New(round Round, height Height, parentHeader sig.Hash, txs tx.Transactions) Block {
	block := Block{
		Time:         time.Now(),
		Round:        round,
		Height:       height,
		ParentHeader: parentHeader,
		Txs:          txs,
	}
	// FIXME: At this point we should calculate the block header. It should be the SHA3 hash of the timestamp, round,
	// height, parentHeader, and transactions.
	block.Header = sig.Hash{}
	return block
}

func (block Block) Sign(signer sig.Signer) (SignedBlock, error) {
	signedBlock := SignedBlock{
		Block: block,
	}

	signature, err := signer.Sign(signedBlock.Header)
	if err != nil {
		return SignedBlock{}, err
	}
	signedBlock.Signature = signature
	signedBlock.Signatory = signer.Signatory()

	return signedBlock, nil
}

func (block Block) String() string {
	return fmt.Sprintf("Block(Header=%s,Round=%d,Height=%d)", base64.StdEncoding.EncodeToString(block.Header[:]), block.Round, block.Height)
}

type SignedBlock struct {
	Block

	Signature sig.Signature
	Signatory sig.Signatory
}

func Genesis() SignedBlock {
	return SignedBlock{
		Block: Block{
			Time:         time.Unix(0, 0),
			Round:        0,
			Height:       0,
			Header:       sig.Hash{},
			ParentHeader: sig.Hash{},
			Txs:          tx.Transactions{},
		},
		Signature: sig.Signature{},
		Signatory: sig.Signatory{},
	}
}

func (signedBlock SignedBlock) String() string {
	return fmt.Sprintf("Block(Header=%s,Round=%d,Height=%d)", base64.StdEncoding.EncodeToString(signedBlock.Header[:]), signedBlock.Round, signedBlock.Height)
}

type Blockchain struct {
	head   Commit
	blocks map[sig.Hash]Commit
}

func NewBlockchain() Blockchain {
	genesis := Genesis()
	genesisCommit := Commit{
		Polka: Polka{
			Block: &genesis,
		},
	}
	return Blockchain{
		head:   genesisCommit,
		blocks: map[sig.Hash]Commit{genesis.Header: genesisCommit},
	}
}

func (blockchain *Blockchain) Height() Height {
	if blockchain.head.Polka.Block == nil {
		return Genesis().Height
	}
	return blockchain.head.Polka.Block.Height
}

func (blockchain *Blockchain) Round() Round {
	if blockchain.head.Polka.Block == nil {
		return Genesis().Round
	}
	return blockchain.head.Polka.Block.Round
}

func (blockchain *Blockchain) Head() (SignedBlock, bool) {
	if blockchain.head.Polka.Block == nil {
		return Genesis(), false
	}
	return *blockchain.head.Polka.Block, true
}

func (blockchain *Blockchain) Block(header sig.Hash) (SignedBlock, bool) {
	commit, ok := blockchain.blocks[header]
	if !ok || commit.Polka.Block == nil {
		return Genesis(), false
	}
	return *commit.Polka.Block, true
}

func (blockchain *Blockchain) Extend(commitToNextBlock Commit) {
	if commitToNextBlock.Polka.Block == nil {
		return
	}
	blockchain.blocks[commitToNextBlock.Polka.Block.Header] = commitToNextBlock
	blockchain.head = commitToNextBlock
}
