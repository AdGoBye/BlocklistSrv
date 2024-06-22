// Package Pushers.grafana: because playing ball with the rest of the ecosystem is for dweebs
package Pushers

import (
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/gofiber/fiber/v2/log"
	goapi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/client/service_accounts"
	"github.com/grafana/grafana-openapi-client-go/models"
	"net/url"
	"os"
	"slices"
	"strings"
	"time"
)

var (
	grafanaUrl, _ = url.Parse(os.Getenv("GRAFANA_LOCATION"))
	clientCfg     = &goapi.TransportConfig{
		Host:     grafanaUrl.Host,
		Schemes:  []string{grafanaUrl.Scheme},
		BasePath: "/api",
		OrgID:    1,
	}
	client goapi.GrafanaHTTPAPI
)

// pointer - Thanks to https://stackoverflow.com/a/77494966
func pointer[T any](d T) *T {
	return &d
}

func EnsureGrafanaUp() bool {
	return tryGetToken()
}

// There are multiple problems with config for Docker and Grafana:
//
// - Docker does not provide a mechanism to push dynamically secrets to containers
//
// - Grafana does not allow us to provision custom service credentials for our applications
//
// - Grafana does not take environment variables, not allowing us to share them
//
// This exists as crowbar approach, tryGetToken first tries to check if we have a service credential and prefer that.
// If it doesn't exist or is faulty, tryGetToken will issue itself a token with admin credentials and save it for reruns.
func tryGetToken() bool {
	token, err := readServiceCred()
	if err != nil || token == "" {
		token, err = issueServiceToken()
		if err != nil {
			log.Debugf("tryGetToken: Failed to issue service token ('%s'), refusing to enable webhook endpoint", err.Error())
			return false
		}
	}

	clientCfg.BasicAuth = nil
	clientCfg.APIKey = token
	client = *goapi.NewHTTPClientWithConfig(strfmt.Default, clientCfg)
	_, err = grafanaGetSelf()
	if err != nil {
		token, err = issueServiceToken() // The token we read is probably stale, get another one
		if err != nil {
			panic(err)
		}
		clientCfg.APIKey = token
		client = *goapi.NewHTTPClientWithConfig(strfmt.Default, clientCfg)
		_, err = grafanaGetSelf()
		if err != nil {
			log.Errorf("tryGetToken: We issued a token but it doesn't work?! ('%s'), refusing to enable webhook endpoint", err.Error())
			return false
		}
	}
	return true
}
func grafanaGetSelf() (*models.UserProfileDTO, error) {
	user, err := client.SignedInUser.GetSignedInUser()
	if err != nil {
		return nil, nil
	}
	return user.GetPayload(), err
}

func issueServiceToken() (string, error) {
	password, err := readFileFromEnvVarLocation("GF_SECURITY_ADMIN_PASSWORD__FILE")
	if err != nil {
		return "", err
	}
	clientCfg.APIKey = ""
	clientCfg.BasicAuth = url.UserPassword("admin", password)
	client = *goapi.NewHTTPClientWithConfig(strfmt.Default, clientCfg)

	// If some exec with power over the Grafana API team is reading this... this does not spark joy as developer.
	// Like really, this is barely better than just writing by hand. I shouldn't need to scrounge a utility function
	// for an API wrapper. You're supposed to abstract this for me as API consumer.
	params := service_accounts.NewSearchOrgServiceAccountsWithPagingParams()
	params.Disabled = pointer(false)

	paging, err := client.ServiceAccounts.SearchOrgServiceAccountsWithPaging(params)
	if err != nil {
		return "", err
	}
	accounts := paging.GetPayload()

	var serviceAccountId int64 = 0
	if len(accounts.ServiceAccounts) == 0 {
		serviceAccountId, err = createServiceAccount(&client)
		if err != nil {
			return "", err
		}
	} else {
		for _, account := range accounts.ServiceAccounts {
			if account.Name == "blocklistsrv" {
				serviceAccountId = account.ID
			}
		}
	}

	paramsServiceToken := service_accounts.NewCreateTokenParams()
	paramsServiceToken.ServiceAccountID = serviceAccountId
	paramsServiceToken.Body = pointer(models.AddServiceAccountTokenCommand{Name: "blocklistsrvannotations"})
	tokenContainer, err := client.ServiceAccounts.CreateToken(paramsServiceToken)
	if err != nil {
		return "", err
	}
	err = writeServiceCred([]byte(tokenContainer.GetPayload().Key))
	if err != nil {
		panic(err)
	}
	return tokenContainer.GetPayload().Key, nil
}

func createServiceAccount(client *goapi.GrafanaHTTPAPI) (int64, error) {
	paramsServiceAccount := service_accounts.NewCreateServiceAccountParams()
	paramsServiceAccount.Body = pointer(models.CreateServiceAccountForm{
		Name: "blocklistsrv",
		Role: "Editor",
	})

	serviceAccountCreationContainer, err := client.ServiceAccounts.CreateServiceAccount(paramsServiceAccount)
	if err != nil {
		return 0, err
	}
	return serviceAccountCreationContainer.GetPayload().ID, nil
}

func readFileFromEnvVarLocation(envvar string) (string, error) {
	passwordFileLocation, ok := os.LookupEnv(envvar)
	if passwordFileLocation == "" || ok == false {
		return "", errors.New(envvar + " unset or default")
	}

	if _, err := os.Stat(passwordFileLocation); err != nil {
		return "", err
	}

	content, err := os.ReadFile(passwordFileLocation)
	if err != nil {
		return "", err
	}
	if string(content) == "" {
		return "", errors.New(envvar + "was empty or something")
	}

	return string(content), nil
}

func constructAnnotationGrafana(WebhookObj GithubPushWebhookObj) {
	_, err := client.Annotations.PostAnnotation(pointer(models.PostAnnotationsCmd{
		Time:    time.Now().UnixMilli(),
		TimeEnd: time.Now().UnixMilli(),
		Tags:    []string{"gitpush"},
		Text:    generateGrafanaAnnotationText(WebhookObj),
	}))
	if err != nil {
		fmt.Println(err)
	}
}

func generateGrafanaAnnotationText(webhookObj GithubPushWebhookObj) *string {
	var filesChanged []string
	for _, commit := range webhookObj.Commits {
		for _, modifiedFile := range commit.Modified {
			if !slices.Contains(filesChanged, modifiedFile) {
				filesChanged = append(filesChanged, modifiedFile)
			}
		}
		for _, modifiedFile := range commit.Added {
			if !slices.Contains(filesChanged, modifiedFile) {
				filesChanged = append(filesChanged, modifiedFile)
			}
		}
	}
	var builder strings.Builder
	builder.WriteString("[" + strings.Join(filesChanged, ", ") + "]:\n")
	for _, commit := range webhookObj.Commits {
		builder.WriteString(commit.Id)
		builder.WriteString(": ")
		builder.WriteString(strings.Split(commit.Message, "\n")[0])
		builder.WriteString("\n")
	}
	return pointer(builder.String())
}
func readServiceCred() (string, error) {
	file, err := os.ReadFile("/.grafanaServiceCredential")
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func writeServiceCred(data []byte) error {
	err := os.WriteFile("/.grafanaServiceCredential", data, 0777)
	if err != nil {
		return err
	}
	return nil
}
