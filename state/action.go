package state

import (
	"io"

	"github.com/renproject/hyperdrive/block"
)

// An Action is emitted by the state Machine to signal to other packages that
// some Action needs to be broadcast to other state Machines that are
// participating in the consensus algorithm.
type Action interface {

	// IsAction is a marker function. It is implemented by types to ensure that
	// we cannot accidentally use the wrong types is some functions. We use type
	// switching is used to enumerate the possible concrete types.
	IsAction()
}

// Propose a Block for consensus in the current round. A previously found Commit
// can be included to help locked state Machines to unlock.
type Propose struct {
	block.SignedPropose

	LastCommit Commit
}

// IsAction is a marker function that implements the Action interface for the Propose type.
func (Propose) IsAction() {
}

func (propose Propose) Write(w io.Writer) error {
	if err := propose.SignedPropose.Write(w); err != nil {
		return err
	}
	if err := propose.LastCommit.Write(w); err != nil {
		return err
	}
	return nil
}

func (propose *Propose) Read(r io.Reader) error {
	if err := propose.SignedPropose.Read(r); err != nil {
		return err
	}
	if err := propose.LastCommit.Read(r); err != nil {
		return err
	}
	return nil
}

type PreVote struct {
	block.PreVote
}

// IsAction is a marker function that implements the Action interface for the PreVote type.
func (PreVote) IsAction() {
}

func (preVote PreVote) Write(w io.Writer) error {
	return preVote.PreVote.Write(w)
}

func (preVote *PreVote) Read(r io.Reader) error {
	return preVote.PreVote.Read(r)
}

type SignedPreVote struct {
	block.SignedPreVote
}

// IsAction is a marker function that implements the Action interface for the SignedPreVote type.
func (SignedPreVote) IsAction() {
}

func (preVote SignedPreVote) Write(w io.Writer) error {
	return preVote.SignedPreVote.Write(w)
}

func (preVote *SignedPreVote) Read(r io.Reader) error {
	return preVote.SignedPreVote.Read(r)
}

type PreCommit struct {
	block.PreCommit
}

// IsAction is a marker function that implements the Action interface for the PreCommit type.
func (PreCommit) IsAction() {
}

func (preCommit PreCommit) Write(w io.Writer) error {
	return preCommit.PreCommit.Write(w)
}

func (preCommit *PreCommit) Read(r io.Reader) error {
	return preCommit.PreCommit.Read(r)
}

type SignedPreCommit struct {
	block.SignedPreCommit
}

// IsAction is a marker function that implements the Action interface for the SignedPreCommit type.
func (SignedPreCommit) IsAction() {
}

func (signedPreCommit SignedPreCommit) Write(w io.Writer) error {
	return signedPreCommit.SignedPreCommit.Write(w)
}

func (signedPreCommit *SignedPreCommit) Read(r io.Reader) error {
	return signedPreCommit.SignedPreCommit.Read(r)
}

type Commit struct {
	block.Commit
}

// IsAction is a marker function that implements the Action interface for the Commit type.
func (Commit) IsAction() {
}

func (commit Commit) Write(w io.Writer) error {
	return commit.Commit.Write(w)
}

func (commit *Commit) Read(r io.Reader) error {
	return commit.Commit.Read(r)
}
