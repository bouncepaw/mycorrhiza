package util

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"io"
	"log"
	"sync"
)

var (
	performed sync.Once
	provider  *oidc.Provider
	config    *oauth2.Config
)

func generateOidcProviderAndConfig(clientId string, clientSecret string, providerUri string, redirectUri string, scopes []string) (*oidc.Provider, *oauth2.Config, error) {
	ctx := context.TODO()

	provider, err := oidc.NewProvider(ctx, providerUri)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	config := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectUri,
		Scopes:       scopes,
	}

	return provider, config, nil
}

func GetOidcProviderAndConfig(clientId string, clientSecret string, providerUri string, redirectUri string, scopes []string) (*oidc.Provider, *oauth2.Config, error) {

	// TODO: produce different providers for different parametes?
	// TODO: refresh?
	var err error
	performed.Do(func() {
		provider, config, err = generateOidcProviderAndConfig(clientId, clientSecret, providerUri, redirectUri, scopes)
	})

	if provider == nil || config == nil {
		err = fmt.Errorf("provider and config are already failed")
	}

	if err != nil {
		return nil, nil, err
	}

	return provider, config, nil

}

func RandString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
