package main

import (
	"encoding/base64"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestTokenRenewZeroHours(t *testing.T) {
	tn := &token{}
	tn.renew(0)
	if tn.expiration.After(time.Now()) {
		t.Error("Expiration should be in the past")
	}
}

func TestTokenRenewFuture(t *testing.T) {
	tn := &token{}
	tn.renew(1)
	if tn.expiration.Before(time.Now()) {
		t.Error("Expiration should be in the future")
	}
}

func makeToken(t *testing.T, length int) *token {
	tn, err := generateToken(length)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	return tn
}

func validateSecretLength(t *testing.T, length int) {
	tn := makeToken(t, length)
	decoded, err := base64.URLEncoding.DecodeString(tn.secret)
	if err != nil {
		t.Errorf("Expected valid base64-encoded secret: %s", tn.secret)
	} else if len(decoded) != length {
		t.Errorf("Expected secret length %d but got %d: %s", length, len(decoded), tn.secret)
	}
}

func TestGenerateToken(t *testing.T) {
	validateSecretLength(t, 10)
	validateSecretLength(t, 100)
}

func TestSecretIsUnique(t *testing.T) {
	count := 100000
	oldTokens := make(map[string]bool)
	for x := 0; x < count; x++ {
		tn := makeToken(t, 50)
		if oldTokens[tn.secret] {
			t.Errorf("Expected unique token but found duplicate: %s", tn.secret)
			break
		}
		oldTokens[tn.secret] = true
	}
}

func TestReadPasswdContent(t *testing.T) {
	content := `
## 2018-03-19T14:43:23Z
user1:pass1
userMissingExp:passMissingExp
 ## 2018-03-18T14:43:23Z 
user2:pass2
## 2018-03-17T14:43:23Z 
#user3:pass3
`
	tokens := readTokens(content)
	if len(tokens) != 3 {
		t.Errorf("Expected 3 tokens, found: %d", len(tokens))
	}
	if tokens["user1"].secret != "pass1" {
		t.Errorf("Expected pass1 secret: %s", tokens["user1"].secret)
	}
	if tokens["user1"].expiration.Format(time.RFC3339) != "2018-03-19T14:43:23Z" {
		t.Errorf("Expected different expiration, but got: %s", tokens["user1"].expiration.Format(time.RFC3339))
	}
	if tokens["user2"].secret != "pass2" {
		t.Errorf("Expected pass1 secret: %s", tokens["user2"].secret)
	}
	if tokens["user2"].expiration.Format(time.RFC3339) != "2018-03-18T14:43:23Z" {
		t.Errorf("Expected different expiration, but got: %s", tokens["user2"].expiration.Format(time.RFC3339))
	}
	if tokens["user3"].secret != "pass3" {
		t.Errorf("Expected pass1 secret: %s", tokens["user3"].secret)
	}
	if tokens["user3"].expiration.Format(time.RFC3339) != "2018-03-17T14:43:23Z" {
		t.Errorf("Expected different expiration, but got: %s", tokens["user3"].expiration.Format(time.RFC3339))
	}
}

func sortString(w string) string {
	s := strings.Split(w, "\n")
	sort.Strings(s)
	return strings.Join(s, "\n")
}

func TestWritePasswdContent(t *testing.T) {
	content := `## 2918-03-19T14:43:23Z
user1:{SHA1}pass1
## 2018-02-18T14:43:23Z
#user2:{SHA1}pass2
`
	tokens := readTokens(content)
	newContent := writeTokens(tokens)
	if sortString(newContent) != sortString(content) {
		t.Errorf("Written token content does not match expected: %s", newContent)
	}
}

func TestHash(t *testing.T) {
	secret := "this is a test"
	expected := "{SHA}+ia+Gd5r/5P3C8IwhDTkpEC7rQI="
	actual := hash(secret)
	if expected != actual {
		t.Errorf("Hash mismatch; expected=%s; actual=%s; secret=%s", expected, actual, secret)
	}
}
