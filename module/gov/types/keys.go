package types

import (
	"bytes"
	"fmt"
	btypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qos/types"
	"time"
)

// Key for getting a the next available proposalID from the store
var (
	KeyDelimiter = []byte(":")

	KeyNextProposalID           = []byte("newProposalID")
	PrefixActiveProposalQueue   = []byte("activeProposalQueue")
	PrefixInactiveProposalQueue = []byte("inactiveProposalQueue")
	KeySoftUpgradeProposal      = []byte("upgradeProposal")
)

func KeyProposalSubspace() []byte {
	return []byte("proposals:")
}

// Key for getting a specific proposal from the store
func KeyProposal(proposalID uint64) []byte {
	return []byte(fmt.Sprintf("proposals:%d", proposalID))
}

// Key for getting a specific deposit from the store
func KeyDeposit(proposalID uint64, depositorAddr btypes.Address) []byte {
	return []byte(fmt.Sprintf("deposits:%d:%s", proposalID, depositorAddr.String()))
}

// Key for getting a specific vote from the store
func KeyVote(proposalID uint64, voterAddr btypes.Address) []byte {
	return []byte(fmt.Sprintf("votes:%d:%s", proposalID, voterAddr.String()))
}

// Key for validators set at entering voting period.
func KeyVotingPeriodValidators(proposalID uint64) []byte {
	return []byte(fmt.Sprintf("validators:%d", proposalID))
}

// Key for getting all deposits on a proposal from the store
func KeyDepositsSubspace(proposalID uint64) []byte {
	return []byte(fmt.Sprintf("deposits:%d:", proposalID))
}

// Key for getting all votes on a proposal from the store
func KeyVotesSubspace(proposalID uint64) []byte {
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
func KeyActiveProposalQueueProposal(endTime time.Time, proposalID uint64) []byte {
	return bytes.Join([][]byte{
		PrefixActiveProposalQueue,
		btypes.FormatTimeBytes(endTime),
		types.Uint64ToBigEndian(proposalID),
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
func KeyInactiveProposalQueueProposal(endTime time.Time, proposalID uint64) []byte {
	return bytes.Join([][]byte{
		PrefixInactiveProposalQueue,
		btypes.FormatTimeBytes(endTime),
		types.Uint64ToBigEndian(proposalID),
	}, KeyDelimiter)
}
