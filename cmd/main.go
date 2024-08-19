package main

import (
	"net/http"
	"os"

	"grpc-memo-api/internal/db"
	"grpc-memo-api/internal/db/generated"
	"grpc-memo-api/internal/server"

	"github.com/bufbuild/connect-go"
	"github.com/joho/godotenv"
	"golang.org/x/exp/slog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		return
	}

	// Get database URL from environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		slog.Error("DATABASE_URL is not set in the environment")
		return
	}

	// Initialize database connection
	conn, err := db.NewDB(databaseURL)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		return
	}
	queries := generated.New(conn)

	// Initialize gRPC server
	// MemoServer is a server struct that processes memo-related gRPC requests using the generated database queries.
	memoServer := &server.MemoServer{DB: queries}

	// Create a multiplexer to handle HTTP requests
	// http.NewServeMux() is used to route HTTP requests to the appropriate handler based on the path.
	// Since the gRPC server operates over the HTTP protocol, a multiplexer is needed to map requests to the appropriate gRPC methods.
	mux := http.NewServeMux()

	// Create a handler for the ListMemos endpoint
	// connect.NewUnaryHandler() creates an HTTP handler by specifying the path of the gRPC endpoint and the function to handle the request.
	// This exposes the gRPC method as an HTTP request.
	handler := connect.NewUnaryHandler(
		"/memo.MemoService/ListMemos", // Endpoint path
		memoServer.ListMemos,          // Function to handle the request
	)
	mux.Handle("/memo.MemoService/ListMemos", handler) // Register the handler with the multiplexer

	// Create a handler for the GetMemo endpoint
	handler = connect.NewUnaryHandler(
		"/memo.MemoService/GetMemo", // Endpoint path
		memoServer.GetMemo,          // Function to handle the request
	)
	mux.Handle("/memo.MemoService/GetMemo", handler) // Register the handler with the multiplexer

	// Create a handler for the CreateMemo endpoint
	handler = connect.NewUnaryHandler(
		"/memo.MemoService/CreateMemo", // Endpoint path
		memoServer.CreateMemo,          // Function to handle the request
	)
	mux.Handle("/memo.MemoService/CreateMemo", handler) // Register the handler with the multiplexer

	// Initialize HTTP/2 server
	// &http2.Server{} is a struct that represents the configuration of the HTTP/2 server.
	// This struct provides options to customize the behavior of HTTP/2.
	// Here, we use the default settings.
	h2s := &http2.Server{}

	// Create HTTP server
	// &http.Server{} is a struct that represents the configuration of the HTTP server.
	// The Addr field specifies the address and port on which the server will listen.
	// The Handler field specifies the handler to process requests.
	// h2c.NewHandler(mux, h2s) creates a handler to process HTTP/2 cleartext (unencrypted) requests.
	// This allows the gRPC server to use the HTTP/2 protocol.
	server := &http.Server{
		Addr:    ":50051",                 // Server address and port
		Handler: h2c.NewHandler(mux, h2s), // Set HTTP/2 cleartext handler
	}

	// Log server startup message
	slog.Info("gRPC server is running on port 50051")

	// Start gRPC server
	// server.ListenAndServe() starts the HTTP server on the specified address and port.
	// This causes the gRPC server to start receiving requests.
	// If an error occurs, it is logged using slog.Error.
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Failed to serve", "error", err)
	}
}
