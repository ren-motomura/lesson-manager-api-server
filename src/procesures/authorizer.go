package procesures

import (
	"net/http"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

const SessionCookieName = "Lessonmanager-Session"

func Authorize(r *http.Request) (*models.User, error) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, err
	}

	sessionId := c.Value
	session, err := models.FindSession(sessionId)
	if err != nil {
		return nil, err
	}

	user, err := models.FindUser(session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
