package main

import (
	"bufio"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/apex/log"
)

type token struct {
	expiration time.Time
	secret     string
}

func (t *token) renew(hours int) {
	t.expiration = time.Now()
	t.expiration = t.expiration.Add(time.Hour * time.Duration(hours))
}

func (t *token) isExpired() bool {
	return t.expiration.Before(time.Now())
}

func generateToken(length int) (*token, error) {
	var t *token
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err == nil {
		t = &token{}
		t.secret = base64.URLEncoding.EncodeToString(bytes)
	}
	return t, err
}

func maintainTokens(config *Config) {
	log.Info("Running token maintenance")
	tokens, err := readTokensFromFile(config.PasswdFilename)
	if err != nil {
		log.WithField("err", err).Error("Unable to read htpasswd file for maintenance, creating new file")
	}
	err = writeTokensToFile(tokens, config.PasswdFilename)
	if err != nil {
		log.WithField("err", err).Error("Unable to write htpasswd file for maintenance")
	}
	log.Info("Token maintenance complete")
}

func hash(password string) string {
	d := sha1.New()
	d.Write([]byte(password))
	return "{SHA}" + base64.StdEncoding.EncodeToString(d.Sum(nil))
}

func createOrUpdateToken(config *Config, user string, uri string) (string, error) {
	response := ""
	tokens, err := readTokensFromFile(config.PasswdFilename)
	msg1 := ""
	msg2 := ""
	msg3 := ""
	if err == nil {
		t := tokens[user]
		if t == nil || strings.HasSuffix(uri, "?new") {
			t, err = generateToken(config.TokenByteLength)
			if err == nil {
				msg1 = fmt.Sprintf("Your new Basic HTTP Auth credentials have been generated.")
				msg2 = fmt.Sprintf("User: %s", user)
				msg3 = fmt.Sprintf("Password: %s", t.secret)
				t.secret = hash(t.secret)
				log.WithField("user", user).Info("Generated token")
			}
		} else {
			msg1 = fmt.Sprintf("Refreshed existing token for user '%s'.", user)
			msg2 = "Navigate to /?new to generate a new token."
		}
		response = fmt.Sprintf("<html><head><title>SSO Token</title> <meta http-equiv=\"refresh\" content=\"1800;URL='/'\"></head><body>%s<p>%s<p>%s<p>Note: It can take a few minutes for the new token to synchronize with the SSO proxy.</body></html>", msg1, msg2, msg3)
		if err == nil {
			t.renew(config.TokenValidityHours)
			tokens[user] = t
			err = writeTokensToFile(tokens, config.PasswdFilename)
			if err == nil {
				log.WithField("user", user).Info("Refreshed token")
			}
		}
	}
	return response, err
}

func readTokensFromFile(filename string) (map[string]*token, error) {
	content, err := ioutil.ReadFile(filename)
	if err == nil {
		return readTokens(string(content)), nil
	}
	return nil, err
}

func readTokens(content string) map[string]*token {
	tokens := make(map[string]*token)
	scanner := bufio.NewScanner(strings.NewReader(content))
	var exp time.Time
	var err error
	var haveExp = false
	for scanner.Scan() {
		value := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(value, "## ") {
			value = strings.TrimSpace(value[2:])
			exp, err = time.Parse(time.RFC3339, value)
			if err == nil {
				haveExp = true
			}
		} else if haveExp {
			values := strings.SplitN(value, ":", 2)
			if len(values) > 1 {
				t := &token{
					expiration: exp,
					secret:     values[1],
				}
				user := values[0]
				if strings.HasPrefix(user, "#") {
					user = user[1:]
				}
				tokens[user] = t
			}
			haveExp = false
		}
	}
	return tokens
}

func writeTokensToFile(tokens map[string]*token, filename string) error {
	content := writeTokens(tokens)
	err := ioutil.WriteFile(filename, []byte(content), 0600)
	return err
}

func writeTokens(tokens map[string]*token) string {
	content := ""
	for user, t := range tokens {
		content += fmt.Sprintf("## %s\n", t.expiration.Format(time.RFC3339))
		expired := ""
		if t.isExpired() {
			expired = "#"
		}
		content += fmt.Sprintf("%s%s:%s\n", expired, user, t.secret)
	}
	return content
}
