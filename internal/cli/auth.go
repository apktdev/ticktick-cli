package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/apktdev/ticktick-cli/internal/config"
	"github.com/apktdev/ticktick-cli/internal/ticktick"
	"github.com/spf13/cobra"
)

func (a *app) newAuthCmd() *cobra.Command {
	auth := &cobra.Command{Use: "auth", Short: "OAuth setup and token management"}

	var clientID, clientSecret, redirectURI string
	setClient := &cobra.Command{
		Use:   "set-client",
		Short: "Set OAuth client credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := required(clientID, "--client-id"); err != nil {
				return err
			}
			if err := required(clientSecret, "--client-secret"); err != nil {
				return err
			}
			if err := required(redirectURI, "--redirect-uri"); err != nil {
				return err
			}
			a.cfg.ClientID = clientID
			a.cfg.ClientSecret = clientSecret
			a.cfg.RedirectURI = redirectURI
			return a.print("client credentials saved")
		},
	}
	setClient.Flags().StringVar(&clientID, "client-id", "", "TickTick OAuth client id")
	setClient.Flags().StringVar(&clientSecret, "client-secret", "", "TickTick OAuth client secret")
	setClient.Flags().StringVar(&redirectURI, "redirect-uri", "", "OAuth redirect URI")

	var scope, state string
	loginURL := &cobra.Command{
		Use:   "login-url",
		Short: "Print OAuth authorization URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			if scope == "" {
				scope = "tasks:read tasks:write"
			}
			if state == "" {
				state = fmt.Sprintf("tickcli-%d", time.Now().Unix())
			}
			u, err := ticktick.AuthURL(a.cfg, scope, state)
			if err != nil {
				return err
			}
			return a.print(u)
		},
	}
	loginURL.Flags().StringVar(&scope, "scope", "", "OAuth scope")
	loginURL.Flags().StringVar(&state, "state", "", "OAuth state")

	var code string
	exchange := &cobra.Command{
		Use:   "exchange",
		Short: "Exchange auth code for access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := required(code, "--code"); err != nil {
				return err
			}
			if err := a.client.ExchangeCode(context.Background(), code); err != nil {
				return err
			}
			return a.print("token saved")
		},
	}
	exchange.Flags().StringVar(&code, "code", "", "OAuth authorization code")

	status := &cobra.Command{
		Use:   "status",
		Short: "Show auth/token status",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := config.Path()
			masked := map[string]any{
				"client_id":       a.cfg.ClientID,
				"redirect_uri":    a.cfg.RedirectURI,
				"has_access":      a.cfg.AccessToken != "",
				"has_refresh":     a.cfg.RefreshToken != "",
				"expiry":          a.cfg.Expiry,
				"scope":           a.cfg.Scope,
				"client_secret":   a.cfg.ClientSecret != "",
				"token_type":      a.cfg.TokenType,
				"config_location": path,
			}
			return a.print(masked)
		},
	}

	auth.AddCommand(setClient, loginURL, exchange, status)
	return auth
}
