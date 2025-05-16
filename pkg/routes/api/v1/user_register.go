// Vikunja is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public Licensee as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public Licensee for more details.
//
// You should have received a copy of the GNU Affero General Public Licensee
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package v1

import (
	"errors"
	"net/http"

	"code.vikunja.io/api/pkg/config"
	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/log"
	"code.vikunja.io/api/pkg/models"
	"code.vikunja.io/api/pkg/user"
	"code.vikunja.io/api/pkg/web/handler"

	"github.com/labstack/echo/v4"
)

type UserRegister struct {
	// The language of the new user. Must be a valid IETF BCP 47 language code and exist in Vikunja.
	Language string `json:"language" valid:"language"`
	user.APIUserPassword
}

// RegisterUser is the register handler
// @Summary Register
// @Description Creates a new user account.
// @tags auth
// @Accept json
// @Produce json
// @Param credentials body v1.UserRegister true "The user with credentials to create"
// @Success 200 {object} user.User
// @Failure 400 {object} web.HTTPError "No or invalid user register object provided / User already exists."
// @Failure 500 {object} models.Message "Internal error"
// @Router /register [post]
func RegisterUser(c echo.Context) error {
	if !config.ServiceEnableRegistration.GetBool() {
		return echo.ErrNotFound
	}
	// Check for Request Content
	var userIn *UserRegister
	if err := c.Bind(&userIn); err != nil {
		return c.JSON(http.StatusBadRequest, models.Message{Message: "No or invalid user model provided."})
	}
	if err := c.Validate(userIn); err != nil {
		e := models.ValidationHTTPError{}
		if is := errors.As(err, &e); is {
			return c.JSON(e.HTTPCode, e)
		}

		return handler.HandleHTTPError(err)
	}
	if userIn == nil {
		return c.JSON(http.StatusBadRequest, models.Message{Message: "No or invalid user model provided."})
	}

	s := db.NewSession()
	defer s.Close()

	// Insert the user
	newUser, err := user.CreateUser(s, &user.User{
		Username: userIn.Username,
		Password: userIn.Password,
		Email:    userIn.Email,
		Language: userIn.Language,
	})
	if err != nil {
		_ = s.Rollback()
		return handler.HandleHTTPError(err)
	}

	// Create their initial project
	err = models.CreateNewProjectForUser(s, newUser)
	if err != nil {
		_ = s.Rollback()
		return handler.HandleHTTPError(err)
	}

	// Get user's first project
	project := models.Project{}
	res, _, _, err := project.ReadAll(s, newUser, "", 1, 20)
	if err != nil {
		_ = s.Rollback()
		return handler.HandleHTTPError(err)
	}

	log.Infof("**user_register** User %s has %d projects", newUser.Username, len(res.([]*models.Project)))

	projects := res.([]*models.Project)

	var firstP *models.Project

	for _, p := range projects {
		if p.ID < 0 {
			continue
		} else {
			firstP = p
			break
		}
	}

	firstT, err := models.CreateDefaultTeamForNewUser(s, newUser)
	if err != nil {
		_ = s.Rollback()
		return handler.HandleHTTPError(err)
	}

	if firstP != nil {
		projectUser := &models.TeamProject{
			ProjectID: firstP.ID,
			TeamID:    firstT.ID,
			Right:     models.RightWrite,
		}

		err = projectUser.Create(s, newUser)
		if err != nil {
			_ = s.Rollback()
			return handler.HandleHTTPError(err)
		}

	}

	if err := s.Commit(); err != nil {
		_ = s.Rollback()
		return handler.HandleHTTPError(err)
	}

	return c.JSON(http.StatusOK, newUser)

}
