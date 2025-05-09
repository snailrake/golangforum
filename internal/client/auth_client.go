package client

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/snailrake/sstu-auth-proto/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client auth.AuthServiceClient
}

func NewClient(addr string) (*AuthClient, error) {
	dialer := func(ctx context.Context, address string) (net.Conn, error) {
		return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, "tcp", address)
	}
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(dialer),
		grpc.WithConnectParams(grpc.ConnectParams{MinConnectTimeout: 5 * time.Second}),
	)
	if err != nil {
		return nil, err
	}
	return &AuthClient{client: auth.NewAuthServiceClient(conn)}, nil
}

func (a *AuthClient) VerifyToken(token string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := a.client.VerifyToken(ctx, &auth.VerifyTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	return resp.Claims.AsMap(), nil
}

func (a *AuthClient) GetUsername(r *http.Request) (string, error) {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}
	claims, err := a.VerifyToken(parts[1])
	if err != nil {
		return "", err
	}
	name, ok := claims["username"].(string)
	if !ok {
		return "", errors.New("username not found")
	}
	return name, nil
}
