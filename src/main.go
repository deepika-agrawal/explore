package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/deepika-agrawal/explore/domain"
	"github.com/deepika-agrawal/explore/pb"
	srvr "github.com/deepika-agrawal/explore/server"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	// PostgreSQL connection
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_SSL_MODE"),
	)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	repo := domain.NewDecisionRepositoryDatabase(db)
	exploreSrvr := srvr.ExploreServiceServer{Repo: repo}
	grpcServer := grpc.NewServer()
	pb.RegisterExploreServiceServer(grpcServer, &exploreSrvr)

	fmt.Println("Server is running on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
