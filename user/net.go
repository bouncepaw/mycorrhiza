package user

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"github.com/bouncepaw/mycorrhiza/cfg"
	"github.com/bouncepaw/mycorrhiza/util"
	"golang.org/x/crypto/bcrypt"
)

// CanProceed returns `true` if the user in `rq` has enough rights to access `route`.
func CanProceed(rq *http.Request, route string) bool {
	return FromRequest(rq).CanProceed(route)
}

// FromRequest returns user from `rq`. If there is no user, an anon user is returned instead.
func FromRequest(rq *http.Request) *User {
	cookie, err := rq.Cookie("mycorrhiza_token")
	if err != nil {
		return EmptyUser()
	}
	return UserByToken(cookie.Value)
}

// LogoutFromRequest logs the user in `rq` out and rewrites the cookie in `w`.
func LogoutFromRequest(w http.ResponseWriter, rq *http.Request) {
	cookieFromUser, err := rq.Cookie("mycorrhiza_token")
	if err == nil {
		http.SetCookie(w, cookie("token", "", time.Unix(0, 0)))
		terminateSession(cookieFromUser.Value)
	}
}

// Register registers the given user. If it fails, a non-nil error is returned.
func Register(username, password, group, source string, force bool) error {
	username = util.CanonicalName(username)

	switch {
	case !util.IsPossibleUsername(username):
		return fmt.Errorf("illegal username \"%s\"", username)
	case !ValidGroup(group):
		return fmt.Errorf("invalid group \"%s\"", group)
	case !ValidSource(source):
		return fmt.Errorf("invalid source \"%s\"", source)
	case HasUsername(username):
		return fmt.Errorf("username \"%s\" is already taken", username)
	case !force && cfg.RegistrationLimit > 0 && Count() >= cfg.RegistrationLimit:
		return fmt.Errorf("reached the limit of registered users (%d)", cfg.RegistrationLimit)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u := User{
		Name:         username,
		Group:        group,
		Source: source,
		Password:     string(hash),
		RegisteredAt: time.Now(),
	}
	users.Store(username, &u)
	return SaveUserDatabase()
}

// LoginDataHTTP logs such user in and returns string representation of an error if there is any.
//
// The HTTP parameters are used for setting header status (bad request, if it is bad) and saving a cookie.
func LoginDataHTTP(w http.ResponseWriter, rq *http.Request, username, password string) string {
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	if !HasUsername(username) {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Unknown username", username, "was entered")
		return "unknown username"
	}
	if !CredentialsOK(username, password) {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("A wrong password was entered for username", username)
		return "wrong password"
	}
	token, err := AddSession(username)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return err.Error()
	}
	http.SetCookie(w, cookie("token", token, time.Now().Add(365*24*time.Hour)))
	return ""
}

// AddSession saves a session for `username` and returns a token to use.
func AddSession(username string) (string, error) {
	token, err := util.RandomString(16)
	if err == nil {
		commenceSession(username, token)
		log.Println("New token for", username, "is", token)
	}
	return token, err
}

// A handy cookie constructor
func cookie(name_suffix, val string, t time.Time) *http.Cookie {
	return &http.Cookie{
		Name:    "mycorrhiza_" + name_suffix,
		Value:   val,
		Expires: t,
		Path:    "/",
	}
}

// TelegramAuthParamsAreValid is true if the given params are ok.
func TelegramAuthParamsAreValid(params map[string][]string) bool {
	// According to the Telegram documentation,
	// > You can verify the authentication and the integrity of the data received by comparing the received hash parameter with the hexadecimal representation of the HMAC-SHA-256 signature of the data-check-string with the SHA256 hash of the bot's token used as a secret key.
	tokenHash := sha256.New()
	tokenHash.Write([]byte(cfg.TelegramBotToken))
	secretKey := tokenHash.Sum(nil)

	hash := hmac.New(sha256.New, secretKey)
	hash.Write([]byte(telegramDataCheckString(params)))
	hexHash := hex.EncodeToString(hash.Sum(nil))

	passedHash := params["hash"][0]
	return passedHash == hexHash
}

// According to the Telegram documentation,
// > Data-check-string is a concatenation of all received fields, sorted in alphabetical order, in the format key=<value> with a line feed character ('\n', 0x0A) used as separator – e.g., 'auth_date=<auth_date>\nfirst_name=<first_name>\nid=<id>\nusername=<username>'.
//
// Note that hash is not used here.
func telegramDataCheckString(params map[string][]string) string {
	var lines []string
	for key, value := range params {
		if key == "hash" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s=%s", key, value[0]))
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}
