package session

import (
	"encoding/base32"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

type InMemoryStore struct {
	sessions map[string]*sessions.Session
	Codecs   []securecookie.Codec
	Options  *sessions.Options // default configuration
	mu       sync.RWMutex
}

const maxAge = 600 // 10 minutes
// 86400 * 30 // 30 days
// NewInMemoryStore creates a new in-memory session store.
func NewInMemoryStore(keyPairs ...[]byte) *InMemoryStore {
	store := &InMemoryStore{
		sessions: make(map[string]*sessions.Session),
		Codecs:   securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:   "/",
			MaxAge: maxAge,
		},
	}

	store.MaxAge(maxAge)
	return store
}

// MaxAge sets the maximum age for the store and the underlying cookie
// implementation. Individual sessions can be deleted by setting Options.MaxAge
// = -1 for that session.
func (store *InMemoryStore) MaxAge(age int) {
	store.Options.MaxAge = age

	// Set the maxAge for each securecookie instance.
	for _, codec := range store.Codecs {
		if sc, ok := codec.(*securecookie.SecureCookie); ok {
			sc.MaxAge(age)
		}
	}
}

// Get retrieves a session by name and request.
func (store *InMemoryStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	if !isCookieNameValid(name) {
		return nil, fmt.Errorf("sessions: invalid character in cookie name: %s", name)
	}
	store.mu.RLock()
	defer store.mu.RUnlock()
	cookie, err := r.Cookie(name)
	var session *sessions.Session
	if err != nil {
		return nil, fmt.Errorf("sessions: missing cookie: %s", err.Error())
	}
	exists := false
	for j := range store.sessions {
		// encoded, err := securecookie.EncodeMulti(store.sessions[j].Name(), store.sessions[j].ID, store.Codecs...)
		// if err != nil {
		// 	return nil, fmt.Errorf("sessions: invalid session content: %s", err.Error())
		// }
		// if cookie.Value == encoded {
		if isOldSession(store.sessions[j]) {
			// Session is old, delete it
			if err := store.delete(store.sessions[j]); err != nil {
				return nil, err
			}
			continue
		}
		if cookie.Value == store.sessions[j].ID {
			exists = true
			session = store.sessions[j]
		}
	}
	if !exists {
		// Create a new session if it doesn't exist
		session = sessions.NewSession(store, name)
		session.Options = store.Options
	}
	return session, nil
}

// New creates a new session
func (store *InMemoryStore) New(r *http.Request, name string) (*sessions.Session, error) {
	if !isCookieNameValid(name) {
		return nil, fmt.Errorf("sessions: invalid character in cookie name: %s", name)
	}
	session := sessions.NewSession(store, name)
	session.Options = store.Options
	return session, nil
}

var base32RawStdEncoding = base32.StdEncoding.WithPadding(base32.NoPadding)

const lastSaveKey = "_last_save"

// isOldSession decides that the session have expired based on MaxAge.
func isOldSession(session *sessions.Session) bool {
	var lastSave time.Time
	if flashes := session.Flashes(lastSaveKey); len(flashes) > 0 {
		if t, ok := flashes[0].(time.Time); ok {
			lastSave = t
		}
	} else {
		return false
	}
	return int64(session.Options.MaxAge) <= int64(time.Since(lastSave).Seconds())
}

// Save saves the session to the in-memory store.
func (store *InMemoryStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if isOldSession(session) {
		if err := store.delete(session); err != nil {
			return err
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		session.ID = base32RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	}
	session.AddFlash(time.Now(), lastSaveKey)
	store.sessions[session.ID] = session

	//  encoded, err := securecookie.EncodeMulti(session.Name(), session.ID,
	// 	store.Codecs...)
	// if err != nil {
	// 	return err
	// }
	// http.SetCookie(w, sessions.NewCookie(session.Name(), encoded, session.Options))
	http.SetCookie(w, sessions.NewCookie(session.Name(), session.ID, session.Options))
	return nil
}

// delete session
func (store *InMemoryStore) delete(session *sessions.Session) error {
	delete(store.sessions, session.ID)
	return nil
}

var isTokenTable = [127]bool{
	'!':  true,
	'#':  true,
	'$':  true,
	'%':  true,
	'&':  true,
	'\'': true,
	'*':  true,
	'+':  true,
	'-':  true,
	'.':  true,
	'0':  true,
	'1':  true,
	'2':  true,
	'3':  true,
	'4':  true,
	'5':  true,
	'6':  true,
	'7':  true,
	'8':  true,
	'9':  true,
	'A':  true,
	'B':  true,
	'C':  true,
	'D':  true,
	'E':  true,
	'F':  true,
	'G':  true,
	'H':  true,
	'I':  true,
	'J':  true,
	'K':  true,
	'L':  true,
	'M':  true,
	'N':  true,
	'O':  true,
	'P':  true,
	'Q':  true,
	'R':  true,
	'S':  true,
	'T':  true,
	'U':  true,
	'W':  true,
	'V':  true,
	'X':  true,
	'Y':  true,
	'Z':  true,
	'^':  true,
	'_':  true,
	'`':  true,
	'a':  true,
	'b':  true,
	'c':  true,
	'd':  true,
	'e':  true,
	'f':  true,
	'g':  true,
	'h':  true,
	'i':  true,
	'j':  true,
	'k':  true,
	'l':  true,
	'm':  true,
	'n':  true,
	'o':  true,
	'p':  true,
	'q':  true,
	'r':  true,
	's':  true,
	't':  true,
	'u':  true,
	'v':  true,
	'w':  true,
	'x':  true,
	'y':  true,
	'z':  true,
	'|':  true,
	'~':  true,
}

func isToken(r rune) bool {
	i := int(r)
	return i < len(isTokenTable) && isTokenTable[i]
}

func isNotToken(r rune) bool {
	return !isToken(r)
}

func isCookieNameValid(raw string) bool {
	if raw == "" {
		return false
	}
	return strings.IndexFunc(raw, isNotToken) < 0
}
