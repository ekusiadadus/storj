// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package metabase_test

import (
	"testing"

	"storj.io/common/storj"
	"storj.io/common/testcontext"
	"storj.io/common/testrand"
	"storj.io/storj/satellite/metabase"
	"storj.io/storj/satellite/metabase/metabasetest"
)

func TestGetStreamPieceCountByNodeID(t *testing.T) {
	metabasetest.Run(t, func(ctx *testcontext.Context, t *testing.T, db *metabase.DB) {
		obj := metabasetest.RandObjectStream()

		t.Run("StreamID missing", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			metabasetest.GetStreamPieceCountByNodeID{
				Opts:     metabase.GetStreamPieceCountByNodeID{},
				ErrClass: &metabase.ErrInvalidRequest,
				ErrText:  "StreamID missing",
			}.Check(ctx, t, db)

			metabasetest.Verify{}.Check(ctx, t, db)
		})

		t.Run("no segments", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			metabasetest.GetStreamPieceCountByNodeID{
				Opts: metabase.GetStreamPieceCountByNodeID{
					StreamID: obj.StreamID,
				},
				Result: map[storj.NodeID]int64{},
			}.Check(ctx, t, db)

			metabasetest.Verify{}.Check(ctx, t, db)
		})

		t.Run("inline segments", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			metabasetest.BeginObjectExactVersion{
				Opts: metabase.BeginObjectExactVersion{
					ObjectStream: obj,
					Encryption:   metabasetest.DefaultEncryption,
				},
				Version: 1,
			}.Check(ctx, t, db)

			encryptedKey := testrand.Bytes(32)
			encryptedKeyNonce := testrand.Bytes(32)

			metabasetest.CommitInlineSegment{
				Opts: metabase.CommitInlineSegment{
					ObjectStream: obj,
					Position:     metabase.SegmentPosition{Part: 0, Index: 0},
					InlineData:   []byte{1, 2, 3},

					EncryptedKey:      encryptedKey,
					EncryptedKeyNonce: encryptedKeyNonce,

					PlainSize:   512,
					PlainOffset: 0,
				},
			}.Check(ctx, t, db)

			metabasetest.GetStreamPieceCountByNodeID{
				Opts: metabase.GetStreamPieceCountByNodeID{
					StreamID: obj.StreamID,
				},
				Result: map[storj.NodeID]int64{},
			}.Check(ctx, t, db)
		})

		t.Run("remote segments", func(t *testing.T) {
			defer metabasetest.DeleteAll{}.Check(ctx, t, db)

			metabasetest.BeginObjectExactVersion{
				Opts: metabase.BeginObjectExactVersion{
					ObjectStream: obj,
					Encryption:   metabasetest.DefaultEncryption,
				},
				Version: 1,
			}.Check(ctx, t, db)

			encryptedKey := testrand.Bytes(32)
			encryptedKeyNonce := testrand.Bytes(32)

			n01 := testrand.NodeID()
			n02 := testrand.NodeID()
			n03 := testrand.NodeID()

			metabasetest.CommitSegment{
				Opts: metabase.CommitSegment{
					ObjectStream: obj,
					Position:     metabase.SegmentPosition{Part: 0, Index: 0},
					RootPieceID:  testrand.PieceID(),

					Pieces: metabase.Pieces{
						{Number: 1, StorageNode: n01},
						{Number: 2, StorageNode: n02},
					},

					EncryptedKey:      encryptedKey,
					EncryptedKeyNonce: encryptedKeyNonce,

					EncryptedSize: 1024,
					PlainSize:     512,
					PlainOffset:   0,
					Redundancy:    metabasetest.DefaultRedundancy,
				},
			}.Check(ctx, t, db)

			metabasetest.CommitSegment{
				Opts: metabase.CommitSegment{
					ObjectStream: obj,
					Position:     metabase.SegmentPosition{Part: 1, Index: 56},
					RootPieceID:  testrand.PieceID(),

					Pieces: metabase.Pieces{
						{Number: 1, StorageNode: n02},
						{Number: 2, StorageNode: n03},
						{Number: 3, StorageNode: n03},
					},

					EncryptedKey:      encryptedKey,
					EncryptedKeyNonce: encryptedKeyNonce,

					EncryptedSize: 1024,
					PlainSize:     512,
					PlainOffset:   0,
					Redundancy:    metabasetest.DefaultRedundancy,
				},
			}.Check(ctx, t, db)

			metabasetest.GetStreamPieceCountByNodeID{
				Opts: metabase.GetStreamPieceCountByNodeID{
					StreamID: obj.StreamID,
				},
				Result: map[storj.NodeID]int64{
					n01: 1,
					n02: 2,
					n03: 2,
				},
			}.Check(ctx, t, db)
		})
	})
}
