package main

import (
	authpb "auth_service/proto"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"net/mail"
	"os"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/argon2"
	"google.golang.org/grpc"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer
	db *pgxpool.Pool
	sessionSecret string
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)

	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		3,
		64*1024,
		4,
		32,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("%s$%s", encodedSalt, encodedHash), nil
}

func VerifyPassword(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 2 {
		return false, fmt.Errorf("invalid password hash format")
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	originalHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		3,
		64*1024,
		4,
		uint32(len(originalHash)),
	)

	return subtle.ConstantTimeCompare(hash, originalHash) == 1, nil
}
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}

	if addr.Address != email {
		return false
	}

	return true
}
func SessionGenerator(email string, user_id string, secret string) string {
    payload := "email="+email+"&user_id="+user_id
    mac := hmac.New(sha512.New,[]byte(secret));
    mac.Write([]byte(payload))
    return hex.EncodeToString(mac.Sum(nil))
}
func VerifyHmacSHA512(email, user_id, secretKey string, mac string) bool {
	providedMAC, err := hex.DecodeString(mac)
	if err != nil {
		return false
	}
	payload := "email="+email+"&user_id="+user_id

	h := hmac.New(sha512.New, []byte(secretKey))
	h.Write([]byte(payload))
	expectedMAC := h.Sum(nil)

	return hmac.Equal(expectedMAC, providedMAC)
}
func VerifySession(session_id string, secret string) bool {
    return false
}
func NewAuthServer(db *pgxpool.Pool, sessionSecret string) *AuthServer {
	return &AuthServer{
		db: db,
		sessionSecret: sessionSecret,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	if !IsValidEmail(email) {
	// Temporary will change in future.
	// I did not think about this part.
    	return &authpb.LoginResponse {
             Email: "None",
             Username: "None",
             SessionId: "None",
             UserId: "None",
    	}, nil
	}
	var email_val string
	var password_hash string
	var id string
	var username string
	error := s.db.QueryRow(ctx, "SELECT email, name, password, id FROM users WHERE email=$1", email).Scan(&email_val, &username, &password_hash, &id)
	if error ==  pgx.ErrNoRows {
	    return &authpb.LoginResponse{
            Email: "None",
            Username: "None",
            SessionId: "None",
            UserId: "None",
		}, nil
	}
	is_correct, err := VerifyPassword(password, password_hash)
	// when err exist, return all none
	// when is not correct, also return all none
	// so I guess we can use or in this case.
	if err != nil || !is_correct {
    	return &authpb.LoginResponse{
                Email: "None",
                Username: "None",
                SessionId: "None",
                UserId: "None",
        }, nil
	}
	session := SessionGenerator(email, id, s.sessionSecret)

	return &authpb.LoginResponse{
	    Email: email,
		Username: username,
		SessionId: session,
		UserId: id,
	}, nil
}
func (s *AuthServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()
	username := req.GetUsername()

	if !IsValidEmail(email) {
		message := "Email is not valid!"
		return &authpb.RegisterResponse{
			Status:  authpb.RegisterStatus_REGISTER_STATUS_INVALID_EMAIL,
			Message: &message,
		}, nil
	}
	password_hash, err := HashPassword(password)
	if err != nil {
		message := "Something went wrong during registration!"
		return &authpb.RegisterResponse{
			Status:  authpb.RegisterStatus_REGISTER_STATUS_UNSPECIFIED,
			Message: &message,
		}, nil
	}
	_, err = s.db.Exec(
		ctx,
		"INSERT INTO users(name, email, password) VALUES($1, $2, $3)",
		username,
		email,
		password_hash,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				message := "Email already exists!"

				return &authpb.RegisterResponse{
					Status:  authpb.RegisterStatus_REGISTER_STATUS_EMAIL_EXISTS,
					Message: &message,
				}, nil

			default:
				log.Printf("PostgreSQL error: %s", pgErr.Message)

				message := "Something went wrong during registration!"
				return &authpb.RegisterResponse{
					Status:  authpb.RegisterStatus_REGISTER_STATUS_UNSPECIFIED,
					Message: &message,
				}, nil
			}
		}

		log.Printf("Database error: %v", err)

		message := "Something went wrong during registration!"
		return &authpb.RegisterResponse{
			Status:  authpb.RegisterStatus_REGISTER_STATUS_UNSPECIFIED,
			Message: &message,
		}, nil
	}
	message := "Register success!"
	return &authpb.RegisterResponse{
		Status:  authpb.RegisterStatus_REGISTER_STATUS_SUCCESS,
		Message: &message,
	}, nil
}

func main() {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Fail to laod .env file")
	}
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	connectionString := os.Getenv("DATABASE_URL")
	// Recommended 512 hex char from TRNG
	sessionSecret := os.Getenv("SECRET_KEY")
	db, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer db.Close()
	if err := db.Ping(ctx); err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	authServer := NewAuthServer(db, sessionSecret)
	authpb.RegisterAuthServiceServer(grpcServer, authServer)
	log.Println("Auth server running on port 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
