package types

import (
	"bytes"
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"time"
)

var (
	KeyDelimiter = []byte(":")
	// Key for getting a the next available proposalID from the store
	KeyNextProposalID = []byte("newProposalID")
	// Prefix key for voting period proposal
	PrefixActiveProposalQueue = []byte("activeProposalQueue")
	// Prefix key for deposit period proposal
	PrefixInactiveProposalQueue = []byte("inactiveProposalQueue")
	// Key for upgrade flag
	KeySoftUpgradeProposal = []byte("upgradeProposal")
)

func KeyProposalSubspace() []byte {
	return []byte("proposals:")
}

// Key for getting a specific proposal from the store
func KeyProposal(proposalID int64) []byte {
	return []byte(fmt.Sprintf("proposals:%d", proposalID))
}

// Key for getting a specific deposit from the store
func KeyDeposit(proposalID int64, depositorAddr btypes.AccAddress) []byte {
	return []byte(fmt.Sprintf("deposits:%d:%s", proposalID, depositorAddr.String()))
}

// Key for getting a specific vote from the store
func KeyVote(proposalID int64, voterAddr btypes.AccAddress) []byte {
	return []byte(fmt.Sprintf("votes:%d:%s", proposalID, voterAddr.String()))
}

// Key for validators set at entering voting period.
func KeyVotingPeriodValidators(proposalID int64) []byte {
	return []byte(fmt.Sprintf("validators:%d", proposalID))
}

// Key for getting all deposits on a proposal from the store
func KeyDepositsSubspace(proposalID int64) []byte {
	return []byte(fmt.Sprintf("deposits:%d:", proposalID))
}

// Key for getting all votes on a proposal from the store
func KeyVotesSubspace(proposalID int64) []byte {
	return []byte(fmt.Sprintf("votes:%d:", proposalID))
}

// Returns the key for a proposalID in the activeProposalQueue
func PrefixActiveProposalQueueTime(endTime time.Time) []byte {
	return bytes.Join([][]byte{
		PrefixActiveProposalQueue,
		btypes.FormatTimeBytes(endTime),
	}, KeyDelimiter)
}

// Returns the key for a proposalID in the activeProposalQueue
func KeyActiveProposalQueueProposal(endTime time.Time, proposalID int64) []byte {
	return bytes.Join([][]byte{
		PrefixActiveProposalQueue,
		btypes.FormatTimeBytes(endTime),
		types.Uint64ToBigEndian(uint64(proposalID)),
	}, KeyDelimiter)
}

// Returns the key for a proposalID in the activeProposalQueue
func PrefixInactiveProposalQueueTime(endTime time.Time) []byte {
	return bytes.Join([][]byte{
		PrefixInactiveProposalQueue,
		btypes.FormatTimeBytes(endTime),
	}, KeyDelimiter)
}

// Returns the key for a proposalID in the activeProposalQueue
func KeyInactiveProposalQueueProposal(endTime time.Time, proposalID int64) []byte {
	return bytes.Join([][]byte{
		PrefixInactiveProposalQueue,
		btypes.FormatTimeBytes(endTime),
		types.Uint64ToBigEndian(uint64(proposalID)),
	}, KeyDelimiter)
}
