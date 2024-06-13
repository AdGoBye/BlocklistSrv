package Pushers

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"os"
	"strings"
	"time"
)

type GrafanaGithubWebhookAnnotation struct{}

var HMACKey = []byte(os.Getenv("GITHUB_WEBHOOK_SECRET"))

func init() {
	if len(HMACKey) == 0 {
		log.Warn("GITHUB_WEBHOOK_SECRET is unset, this bypasses the signature check for webhooks!\n" +
			"You most definitely don't want this in production, this enables anybody to send arbitrary webhook data.")
	}
}

type GithubPushWebhookObj struct {
	Ref        string      `json:"-"`
	Before     string      `json:"-"`
	After      string      `json:"-"`
	Repository struct{}    `json:"-"`
	Pusher     struct{}    `json:"-"`
	Sender     struct{}    `json:"-"`
	Created    bool        `json:"-"`
	Deleted    bool        `json:"-"`
	Forced     bool        `json:"-"`
	BaseRef    interface{} `json:"-"`
	Compare    string      `json:"-"`
	Commits    []struct {
		Id        string    `json:"id"`
		TreeId    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Url       string    `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []string `json:"added"`
		Removed  []string `json:"removed"`
		Modified []string `json:"modified"`
	} `json:"commits"`
	HeadCommit struct{} `json:"-"`
}

func verifyGithubSignature(c *fiber.Ctx) error {
	sigHeader, exists := c.GetReqHeaders()["X-Hub-Signature-256"]
	if !exists || len(sigHeader) == 0 {
		return fiber.NewError(fiber.StatusUnauthorized, "missing signature")
	}

	hmacObj := hmac.New(sha256.New, HMACKey)
	hmacObj.Write(c.Body())
	hmacObjSignature := hmacObj.Sum(nil)

	signatureSplit := strings.Split(sigHeader[0], "=")[1]
	signature, err := hex.DecodeString(signatureSplit)
	if err != nil {
		panic(err)
	}

	// I believe (armchair cryptography creature, correct me) this comparison is overkill, as knowing the length is not
	// too useful to an attacker. I'll still do it as sanity check for ConstantTimeCompare, which fails opaquely
	// if the length is incorrect.
	if subtle.ConstantTimeEq(int32(len(signature)), int32(len(hmacObjSignature))) != 1 {
		return fiber.NewError(fiber.StatusUnauthorized, "signature length mismatch")
	}

	if subtle.ConstantTimeCompare(signature, hmacObjSignature) != 1 {
		return fiber.NewError(fiber.StatusUnauthorized, "hmac mismatch")
	}
	return nil
}
func (grafghanno GrafanaGithubWebhookAnnotation) HandlePushRequest(c *fiber.Ctx) error {
	c.Accepts("application/json")
	var Callback GithubPushWebhookObj
	if len(HMACKey) > 0 {
		err := verifyGithubSignature(c)
		if err != nil {
			return err
		}
	}

	switch c.GetReqHeaders()["X-Github-Event"][0] {
	case "push":
		if err := c.BodyParser(&Callback); err != nil {
			panic(err)
		}
		go constructAnnotationGrafana(Callback)
	case "ping": // Needs no processing
	default:
		return c.SendStatus(fiber.StatusNotImplemented)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (grafghanno GrafanaGithubWebhookAnnotation) CanPusherOperate() bool {
	// ? Is there some way we can check that GitHub is good?
	return EnsureGrafanaUp()
}
