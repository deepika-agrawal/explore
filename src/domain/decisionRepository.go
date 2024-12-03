package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/deepika-agrawal/explore/pb"
	"github.com/jmoiron/sqlx"
)

var defaultPageSize = 20

type DecisionRepository interface {
	ListLikedYou(ctx context.Context, req *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error)    // List all users who liked the recipient
	ListNewLikedYou(ctx context.Context, req *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error) // List all users who liked the recipient excluding those who have been liked in return
	CountLikedYou(ctx context.Context, req *pb.CountLikedYouRequest) (*pb.CountLikedYouResponse, error) // Count the number of users who liked the recipient
	PutDecision(ctx context.Context, req *pb.PutDecisionRequest) (*pb.PutDecisionResponse, error)       // Record the decision of the actor to like or pass the recipient
}

type DecisionRepositoryDatabase struct {
	db *sqlx.DB
}

func NewDecisionRepositoryDatabase(database *sqlx.DB) DecisionRepositoryDatabase {
	return DecisionRepositoryDatabase{
		db: database,
	}
}

func (r DecisionRepositoryDatabase) ListLikedYou(ctx context.Context, req *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error) {
	if r.db == nil {
		log.Println("Db connection is not available")
		return nil, errors.New("internal error")
	}

	offset := 0
	page := 0
	var err error
	if req.PaginationToken != nil {
		page, err = strconv.Atoi(*req.PaginationToken)
		if err != nil {
			log.Printf("Error in parsing pagination token: %v", err)
			return nil, err
		}
		if page > 1 {
			offset = page * defaultPageSize
		}
	}

	var likers []*pb.ListLikedYouResponse_Liker
	query := `SELECT actor_user_id, decision_timestamp
		FROM user_decisions
		WHERE recipient_user_id = $1 AND liked_recipient = true
		ORDER BY decision_timestamp ASC
		LIMIT $2
		OFFSET $3;`

	rows, err := r.db.QueryContext(ctx, query, req.RecipientUserId, defaultPageSize, offset)
	if err != nil {
		log.Printf("Error in getting liked users: %v", err)
		return nil, fmt.Errorf("Error in getting liked users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var actorID string
		var timestamp time.Time
		if err := rows.Scan(&actorID, &timestamp); err != nil {
			log.Printf("Failed to scan row in getting liked users: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		likers = append(likers, &pb.ListLikedYouResponse_Liker{
			ActorId:       actorID,
			UnixTimestamp: uint64(timestamp.Unix()),
		})
	}

	nextPage := ""
	if len(likers) == defaultPageSize {
		strconv.Itoa(page + 1)
	}
	response := pb.ListLikedYouResponse{
		Likers:              likers,
		NextPaginationToken: &nextPage,
	}
	return &response, nil
}

func (r DecisionRepositoryDatabase) ListNewLikedYou(ctx context.Context, req *pb.ListLikedYouRequest) (*pb.ListLikedYouResponse, error) {
	if r.db == nil {
		log.Println("Db connection is not available")
		return nil, errors.New("internal error")
	}

	offset := 0
	page := 0
	var err error
	if req.PaginationToken != nil {
		page, err = strconv.Atoi(*req.PaginationToken)
		if err != nil {
			log.Printf("Error in parsing pagination token: %v", err)
			return nil, err
		}
		if page > 1 {
			offset = page * defaultPageSize
		}
	}

	var likers []*pb.ListLikedYouResponse_Liker
	query := `SELECT actor_user_id, decision_timestamp
		FROM user_decisions
		WHERE recipient_user_id = $1
		AND liked_recipient = TRUE
		AND actor_user_id NOT IN (
			SELECT recipient_user_id
			FROM user_decisions
			WHERE actor_user_id = $1
				AND liked_recipient = TRUE
		)
		ORDER BY decision_timestamp ASC
		LIMIT $2
		OFFSET $3;`

	rows, err := r.db.QueryContext(ctx, query, req.RecipientUserId, defaultPageSize, offset)
	if err != nil {
		log.Printf("Error in getting new liked users: %v", err)
		return nil, fmt.Errorf("Error in getting new liked users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var actorID string
		var timestamp time.Time
		if err := rows.Scan(&actorID, &timestamp); err != nil {
			log.Printf("Failed to scan rows in getting new liked users: %v", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		likers = append(likers, &pb.ListLikedYouResponse_Liker{
			ActorId:       actorID,
			UnixTimestamp: uint64(timestamp.Unix()),
		})
	}

	nextPage := ""
	if len(likers) == defaultPageSize {
		strconv.Itoa(page + 1)
	}
	response := pb.ListLikedYouResponse{
		Likers:              likers,
		NextPaginationToken: &nextPage,
	}
	return &response, nil
}

func (r DecisionRepositoryDatabase) CountLikedYou(ctx context.Context, req *pb.CountLikedYouRequest) (*pb.CountLikedYouResponse, error) {

	if r.db == nil {
		log.Println("Db connection is not available")
		return nil, errors.New("internal error")
	}

	query := `
	SELECT COUNT(*)
	FROM user_decisions
	WHERE recipient_user_id = $1 AND liked_recipient = true;`

	var count int
	err := r.db.QueryRowContext(ctx, query, req.RecipientUserId).Scan(&count)
	if err != nil {
		log.Printf("Error in getting counting users who liked you: %v", err)
		return nil, fmt.Errorf("Error in couting users who liked you: %w", err)
	}

	return &pb.CountLikedYouResponse{Count: uint64(count)}, nil
}

func (r DecisionRepositoryDatabase) PutDecision(ctx context.Context, req *pb.PutDecisionRequest) (*pb.PutDecisionResponse, error) {

	if r.db == nil {
		log.Println("Db connection is not available")
		return nil, errors.New("internal error")
	}

	// Insert or update the decision
	query := `
	INSERT INTO user_decisions (recipient_user_id, actor_user_id, liked_recipient, decision_timestamp)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (recipient_user_id, actor_user_id)
	DO UPDATE SET liked_recipient = EXCLUDED.liked_recipient, decision_timestamp = EXCLUDED.decision_timestamp;`

	_, err := r.db.ExecContext(ctx, query, req.RecipientUserId, req.ActorUserId, req.LikedRecipient, time.Now())
	if err != nil {
		log.Printf("Error in putting decision: %v", err)
		return nil, fmt.Errorf("Error in putting the decision: %w", err)
	}

	// Check for mutual likes
	mutualLikes := false
	if req.LikedRecipient {
		query = `
			SELECT liked_recipient
			FROM user_decisions
			WHERE recipient_user_id = $1 AND actor_user_id = $2;
		`
		err := r.db.QueryRowContext(ctx, query, req.ActorUserId, req.RecipientUserId).Scan(&mutualLikes)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error in checking mutual likes: %v", err)
			return nil, fmt.Errorf("failed to check mutual likes: %w", err)
		}
	}

	return &pb.PutDecisionResponse{MutualLikes: mutualLikes}, nil
}
