package server

import (
	"context"
	"grpc-memo-api/api/proto"
	"grpc-memo-api/internal/db/generated"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

// MemoServer is the implementation of MemoService.
type MemoServer struct {
	DB *generated.Queries // Field for executing database queries
}

// ListMemos lists memos for the specified user.
func (s *MemoServer) ListMemos(ctx context.Context, req *connect.Request[proto.ListMemosRequest]) (*connect.Response[proto.ListMemosResponse], error) {
	// Retrieve memos from the database
	memos, err := s.DB.ListMemos(ctx, req.Msg.UserId)
	if err != nil {
		slog.Error("Failed to list memos", "error", err)
		return nil, err
	}

	// Convert to protocol buffer memo format
	var protoMemos []*proto.Memo
	for _, memo := range memos {
		protoMemos = append(protoMemos, &proto.Memo{
			Id:        memo.ID.String(),
			UserId:    memo.UserID,
			Content:   memo.Content,
			CreatedAt: memo.CreatedAt.Format(time.RFC3339),
		})
	}

	// Return the response
	return connect.NewResponse(&proto.ListMemosResponse{Memos: protoMemos}), nil
}

// GetMemo retrieves a memo corresponding to the specified user and memo ID.
func (s *MemoServer) GetMemo(ctx context.Context, req *connect.Request[proto.GetMemoRequest]) (*connect.Response[proto.GetMemoResponse], error) {
	// Retrieve the memo from the database
	memo, err := s.DB.GetMemo(ctx, generated.GetMemoParams{
		UserID: req.Msg.UserId,
		ID:     uuid.MustParse(req.Msg.MemoId),
	})
	if err != nil {
		slog.Error("Failed to get memo", "error", err)
		return nil, err
	}

	// Return the response
	return connect.NewResponse(&proto.GetMemoResponse{Memo: &proto.Memo{
		Id:        memo.ID.String(),
		UserId:    memo.UserID,
		Content:   memo.Content,
		CreatedAt: memo.CreatedAt.Format(time.RFC3339),
	}}), nil
}

// CreateMemo creates a new memo.
func (s *MemoServer) CreateMemo(ctx context.Context, req *connect.Request[proto.CreateMemoRequest]) (*connect.Response[proto.CreateMemoResponse], error) {
	// Generate a new memo ID and creation timestamp
	memoID := uuid.New()
	createdAt := time.Now()

	// Insert the memo into the database
	params := generated.CreateMemoParams{
		ID:        memoID,
		UserID:    req.Msg.UserId,
		Content:   req.Msg.Content,
		CreatedAt: createdAt,
	}

	err := s.DB.CreateMemo(ctx, params)
	if err != nil {
		slog.Error("Failed to create memo", "error", err)
		return nil, err
	}

	// Return the response
	return connect.NewResponse(&proto.CreateMemoResponse{Memo: &proto.Memo{
		Id:        memoID.String(),
		UserId:    req.Msg.UserId,
		Content:   req.Msg.Content,
		CreatedAt: createdAt.Format(time.RFC3339),
	}}), nil
}

// UpdateMemo update memo.
func (s *MemoServer) UpdateMemo(ctx context.Context, req *connect.Request[proto.UpdateMemoRequest]) (*connect.Response[proto.UpdateMemoResponse], error) {
	// Parse UUID
	memoID, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		slog.Error("id must be uuid", "error", err)
		return nil, err
	}

	// Execute the update operation
	err = s.DB.UpdateMemo(ctx, generated.UpdateMemoParams{
		ID:      memoID,
		Content: req.Msg.Content,
	})
	if err != nil {
		slog.Error("Failed to update memo", "error", err)
		return nil, err
	}

	// After the update, retrieve the updated memo from the database
	updatedMemo, err := s.DB.GetMemo(ctx, generated.GetMemoParams{
		ID:     memoID,
		UserID: req.Msg.UserId,
	})
	if err != nil {
		slog.Error("Failed to retrieve updated memo", "error", err)
		return nil, err
	}

	// Return the response with the updated memo
	return connect.NewResponse(&proto.UpdateMemoResponse{
		Memo: &proto.Memo{
			Id:        updatedMemo.ID.String(),
			UserId:    updatedMemo.UserID,
			Content:   updatedMemo.Content,
			CreatedAt: updatedMemo.CreatedAt.Format(time.RFC3339),
		},
	}), nil
}
