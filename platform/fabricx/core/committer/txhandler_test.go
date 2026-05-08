/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package committer

import (
	"testing"

	cb "github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-x-common/api/committerpb"
	"github.com/stretchr/testify/require"

	"github.com/hyperledger-labs/fabric-smart-client/platform/fabric/core/generic/committer"
	"github.com/hyperledger-labs/fabric-smart-client/platform/fabric/driver"
)

func TestConvertValidationCode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		status   committerpb.Status
		expected driver.ValidationCode
	}{
		{
			name:     "COMMITTED maps to Valid",
			status:   committerpb.Status_COMMITTED,
			expected: driver.Valid,
		},
		{
			name:     "STATUS_UNSPECIFIED maps to Invalid",
			status:   committerpb.Status_STATUS_UNSPECIFIED,
			expected: driver.Invalid,
		},
		{
			name:     "ABORTED_SIGNATURE_INVALID maps to Invalid",
			status:   committerpb.Status_ABORTED_SIGNATURE_INVALID,
			expected: driver.Invalid,
		},
		{
			name:     "unknown status maps to Invalid",
			status:   committerpb.Status(99),
			expected: driver.Invalid,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := convertValidationCode(tc.status)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestNewHandler(t *testing.T) {
	t.Parallel()
	com := &committer.Committer{}
	h := NewHandler(com)
	require.NotNil(t, h)
	require.Same(t, com, h.committer)
}

func TestRegisterTransactionHandler(t *testing.T) {
	t.Parallel()
	com := &committer.Committer{
		Handlers: make(map[cb.HeaderType]committer.TransactionHandler),
	}
	RegisterTransactionHandler(com)
	_, ok := com.Handlers[cb.HeaderType_MESSAGE]
	require.True(t, ok, "Handler for MESSAGE type should be registered")
}

func TestHandleFabricxTransaction_MetadataLacksFilter(t *testing.T) {
	t.Parallel()
	com := &committer.Committer{
		Handlers: make(map[cb.HeaderType]committer.TransactionHandler),
	}
	h := NewHandler(com)

	// Block metadata with insufficient entries (less than statusIdx)
	blkMeta := &cb.BlockMetadata{
		Metadata: [][]byte{},
	}
	tx := committer.CommitTx{
		TxID: "tx_meta_short",
	}

	event, err := h.HandleFabricxTransaction(t.Context(), blkMeta, tx)
	require.Error(t, err)
	require.Nil(t, event)
	require.Contains(t, err.Error(), "block metadata lacks transaction filter")
}

func TestStatusIdxConstant(t *testing.T) {
	t.Parallel()
	// Verify statusIdx matches the fabric constant for TRANSACTIONS_FILTER
	require.Equal(t, int(cb.BlockMetadataIndex_TRANSACTIONS_FILTER), statusIdx)
}
