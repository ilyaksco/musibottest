package userbot

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

type UserbotClient struct {
	Client *telegram.Client
	Self   *tg.User
}

type customTerminalAuth struct {
	phone  string
	reader *bufio.Reader
}

func (a customTerminalAuth) Phone(_ context.Context) (string, error) {
	return a.phone, nil
}

func (a customTerminalAuth) Password(_ context.Context) (string, error) {
	fmt.Print("Enter 2FA Password: ")
	byteLine, _, err := a.reader.ReadLine()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(byteLine)), nil
}

func (a customTerminalAuth) Code(_ context.Context, _ *tg.AuthSentCode) (string, error) {
	fmt.Print("Enter code: ")
	byteLine, _, err := a.reader.ReadLine()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(byteLine)), nil
}

func (a customTerminalAuth) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("sign up not supported")
}

func (a customTerminalAuth) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	return &auth.SignUpRequired{TermsOfService: tos}
}

func Initialize(ctx context.Context, apiID int, apiHash, userbotPhone string) (*UserbotClient, error) {
	client := telegram.NewClient(apiID, apiHash, telegram.Options{
		SessionStorage: &session.FileStorage{
			Path: "userbot.session",
		},
	})

	status, err := client.Auth().Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("auth status failed: %w", err)
	}

	if !status.Authorized {
		log.Println("Log: Userbot session not found. Starting manual interactive login...")
		authenticator := customTerminalAuth{
			phone:  userbotPhone,
			reader: bufio.NewReader(os.Stdin),
		}
		flow := auth.NewFlow(authenticator, auth.SendCodeOptions{})
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return nil, fmt.Errorf("auth flow failed: %w", err)
		}
	}

	self, err := client.Self(ctx)
	if err != nil {
		return nil, fmt.Errorf("get self failed: %w", err)
	}

	log.Printf("Log: Userbot authorized as: %s (@%s)", self.FirstName, self.Username)
	return &UserbotClient{Client: client, Self: self}, nil
}

func (uc *UserbotClient) Start(ctx context.Context) error {
	return uc.Client.Run(ctx, func(ctx context.Context) error {
		log.Println("Log: Userbot client connection is active.")
		<-ctx.Done()
		return ctx.Err()
	})
}