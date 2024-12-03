package domain

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/deepika-agrawal/explore/pb"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestListLikedYou(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDecisionRepositoryDatabase(sqlxDB)

	ctx := context.Background()
	recipientID := "123"
	rows := sqlmock.NewRows([]string{"actor_user_id", "decision_timestamp"}).
		AddRow("456", time.Now()).
		AddRow("789", time.Now())

	mock.ExpectQuery(`SELECT actor_user_id, decision_timestamp FROM user_decisions WHERE recipient_user_id = \$1 AND liked_recipient = true ORDER BY decision_timestamp ASC LIMIT \$2 OFFSET \$3;`).
		WithArgs(recipientID, defaultPageSize, 0).
		WillReturnRows(rows)

	req := &pb.ListLikedYouRequest{RecipientUserId: recipientID}
	resp, err := repo.ListLikedYou(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, resp.Likers, 2)
	assert.Equal(t, "456", resp.Likers[0].ActorId)
}

func TestListNewLikedYou(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDecisionRepositoryDatabase(sqlxDB)

	ctx := context.Background()
	recipientID := "123"
	rows := sqlmock.NewRows([]string{"actor_user_id", "decision_timestamp"}).
		AddRow("456", time.Now())

	mock.ExpectQuery(`SELECT actor_user_id, decision_timestamp FROM user_decisions WHERE recipient_user_id = \$1 AND liked_recipient = TRUE AND actor_user_id NOT IN \(
			SELECT recipient_user_id
			FROM user_decisions
			WHERE actor_user_id = \$1
				AND liked_recipient = TRUE
		\) ORDER BY decision_timestamp ASC LIMIT \$2 OFFSET \$3;`).
		WithArgs(recipientID, defaultPageSize, 0).
		WillReturnRows(rows)

	req := &pb.ListLikedYouRequest{RecipientUserId: recipientID}
	resp, err := repo.ListNewLikedYou(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, resp.Likers, 1)
	assert.Equal(t, "456", resp.Likers[0].ActorId)
}

func TestPutDecision_MutualLike(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDecisionRepositoryDatabase(sqlxDB)

	ctx := context.Background()

	actorID := "123"
	recipientID := "456"
	liked := true

	// Mock insert or update
	mock.ExpectExec(`INSERT INTO user_decisions \(recipient_user_id, actor_user_id, liked_recipient, decision_timestamp\) VALUES \(\$1, \$2, \$3, \$4\) ON CONFLICT \(recipient_user_id, actor_user_id\) DO UPDATE SET liked_recipient = EXCLUDED.liked_recipient, decision_timestamp = EXCLUDED.decision_timestamp;`).
		WithArgs(recipientID, actorID, liked, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Mock mutual like check
	mock.ExpectQuery(`SELECT liked_recipient FROM user_decisions WHERE recipient_user_id = \$1 AND actor_user_id = \$2;`).
		WithArgs(actorID, recipientID).
		WillReturnRows(sqlmock.NewRows([]string{"liked_recipient"}).AddRow(true))

	req := &pb.PutDecisionRequest{ActorUserId: actorID, RecipientUserId: recipientID, LikedRecipient: liked}
	resp, err := repo.PutDecision(ctx, req)

	assert.NoError(t, err)
	assert.True(t, resp.MutualLikes)
}

func TestPutDecision_MutualLikeFalse(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDecisionRepositoryDatabase(sqlxDB)

	ctx := context.Background()

	actorID := "123"
	recipientID := "456"
	liked := false

	// Mock insert or update
	mock.ExpectExec(`INSERT INTO user_decisions \(recipient_user_id, actor_user_id, liked_recipient, decision_timestamp\) VALUES \(\$1, \$2, \$3, \$4\) ON CONFLICT \(recipient_user_id, actor_user_id\) DO UPDATE SET liked_recipient = EXCLUDED.liked_recipient, decision_timestamp = EXCLUDED.decision_timestamp;`).
		WithArgs(recipientID, actorID, liked, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	req := &pb.PutDecisionRequest{ActorUserId: actorID, RecipientUserId: recipientID, LikedRecipient: liked}
	resp, err := repo.PutDecision(ctx, req)

	assert.NoError(t, err)
	assert.False(t, resp.MutualLikes)
}

func TestCountLikedYou(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewDecisionRepositoryDatabase(sqlxDB)

	ctx := context.Background()
	recipientID := "123"
	count := 42

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM user_decisions WHERE recipient_user_id = \$1 AND liked_recipient = true;`).
		WithArgs(recipientID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(count))

	req := &pb.CountLikedYouRequest{RecipientUserId: recipientID}
	resp, err := repo.CountLikedYou(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, uint64(count), resp.Count)
}
