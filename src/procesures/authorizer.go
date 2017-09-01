package procesures

import (
	"net/http"

	"github.com/ren-motomura/lesson-manager-api-server/src/models"
)

const SessionCookieName = "Lessonmanager-Session"

func Authorize(r *http.Request) (*models.Company, error) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, err
	}

	sessionId := c.Value
	session, err := models.FindSession(sessionId)
	if err != nil {
		return nil, err
	}

	company, err := models.FindCompany(session.CompanyID)
	if err != nil {
		return nil, err
	}

	return company, nil
}
