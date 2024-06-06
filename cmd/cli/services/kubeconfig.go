package services

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
)

func readKubeconfig() (*api.Config, error) {
	home := homedir.HomeDir()
	kubeconfigPath := filepath.Join(home, ".kube", "config")

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.ExplicitPath = kubeconfigPath

	config, err := loadingRules.Load()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// tries to find the jwtoken for the specifies cluster by clustername
func getJWTokenFromConfig(clusterName string, config api.Config) (string, error) {
	context := config.Contexts[clusterName]
	if context == nil {
		return "", errors.New("could not find context")
	}
	username := config.Contexts[clusterName].AuthInfo
	if username == "" {
		return "", errors.New("could not find Authinfo")
	}

	token := config.AuthInfos[username].Token
	if token == "" {
		return "", errors.New("could not find JWToken")
	}

	return token, nil
}

func getJWTokenClaims(rawToken string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(rawToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}

func IsTokenExpired(name string) error {
	config, err := readKubeconfig()
	if err != nil {
		return err
	}

	token, err := getJWTokenFromConfig(name, *config)
	if err != nil {
		return err
	}

	claims, err := getJWTokenClaims(token)
	if err != nil {
		return err
	}

	expiry := int64(claims["exp"].(float64))
	now := time.Now().Unix()

	if now > expiry {
		return errors.New("token has expired")
	}
	return nil
}
