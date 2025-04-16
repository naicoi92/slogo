package main

import (
	"errors"
	"log/slog"
	"net/netip"
	"os"

	"github.com/naicoi92/slogo"
)

func main() {
	// Create a new logger with default settings
	logger := slog.New(
		slogo.NewHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
			},
		),
	)

	// Set as default logger
	slog.SetDefault(logger)

	// Log messages
	slog.Info("Hello from SLogo")

	// Log with attributes
	slog.Info("User login", "user_id", 12345, "action", "login")

	// Log structs
	type User struct {
		ID        int        `json:"id"`
		Username  string     `json:"username"`
		Password  string     `json:"password"   slog:"restrict"` // This field will be redacted
		Email     string     `json:"email"`
		IpAddress netip.Addr `json:"ip_address"`
		Members   []*User    `json:"members"`
	}
	ip, err := netip.ParseAddr("10.0.0.1")
	if err != nil {
		slog.Error("Failed to parse IP address", "error", err)
	}
	user := User{
		ID:        1,
		Username:  "johndoe",
		Password:  "secret123",
		Email:     "john@example.com",
		IpAddress: ip,
		Members: []*User{
			{
				ID:       2,
				Username: "janedoe",
			},
			{
				ID:        3,
				Username:  "alice",
				IpAddress: ip,
			},
		},
	}

	slog.Info("User details", "user", user)

	// Log errors with stack traces
	if err := errors.New("an example error"); err != nil {
		slog.Error("Operation failed", "error", err)
	}
	CustomError := errors.New("custom error")

	if err := CustomError; err != nil {
		slog.Error("Operation failed", "error", err)
	}
}
