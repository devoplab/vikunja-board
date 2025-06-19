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

package cmd

import (
	"os"
	"strconv"
	"time"

	"code.vikunja.io/api/pkg/db"
	"code.vikunja.io/api/pkg/initialize"
	"code.vikunja.io/api/pkg/log"
	"code.vikunja.io/api/pkg/models"
	"code.vikunja.io/api/pkg/user"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"xorm.io/xorm"
)

func init() {

	teamCmd.AddCommand(teamListCmd, addUserToAllTeamsCmd, removeUserFromAllTeamsCmd)
	rootCmd.AddCommand(teamCmd)
}

var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage team locally through the cli.",
}

// listAllTeams returns all teams
func listAllTeams(s *xorm.Session) (teams []*models.Team, err error) {
	err = s.Find(&teams)
	return
}

func getCreatorSafe(c *user.User) (creator *user.User) {
	if c == nil {
		return &user.User{
			ID:   0,
			Name: "",
		}
	} else {
		return c
	}
}

var teamListCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows a list of all teams.",
	PreRun: func(_ *cobra.Command, _ []string) {
		initialize.FullInit()
	},

	Run: func(_ *cobra.Command, _ []string) {
		s := db.NewSession()
		defer s.Close()

		// Get all teams
		teams, err := listAllTeams(s)
		if err != nil {
			log.Fatalf("Error getting teams: %s", err)
		}

		// fmt.Println("Teams loaded:", teams)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{
			"ID",
			"Name",
			"Description",
			"CreatedBy",
			"CreatedByID",
			"Members",
			"Created",
			"Updated",
		})

		for _, t := range teams {
			// fmt.Println("Team:", t)
			table.Append([]string{
				strconv.FormatInt(t.ID, 10),
				t.Name,
				t.Description,
				getCreatorSafe(t.CreatedBy).Name,
				strconv.FormatInt(getCreatorSafe(t.CreatedBy).ID, 10),
				strconv.FormatInt(int64(len(t.Members)), 10),
				t.Created.Format(time.RFC3339),
				t.Updated.Format(time.RFC3339),
			})
		}

		table.Render()

	},
}

var addUserToAllTeamsCmd = &cobra.Command{
	Use:   "adduser",
	Short: "Add a user to all teams.",
	PreRun: func(_ *cobra.Command, _ []string) {
		initialize.FullInit()
	},
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Please provide a username.")
		}

		s := db.NewSession()
		defer s.Close()

		u, err := user.GetUserByUsername(s, args[0])
		if err != nil {
			log.Fatalf("Error getting user: %s", err)
		}

		if u == nil {
			log.Fatalf("User %s does not exist.", args[0])
		}

		teams, err := listAllTeams(s)
		if err != nil {
			log.Fatalf("Error getting teams: %s", err)
		}

		for _, t := range teams {
			tm := &models.TeamMember{
				UserID:   u.ID,
				TeamID:   t.ID,
				Username: u.Username,
				Admin:    false,
			}
			err = tm.Create(s, u)
			if err != nil {
				// log.Fatalf("Error adding user to team %s: %s", t.Name, err)
				log.Infof("User not added to %s because: %s", t.Name, err)
			} else {
				log.Infof("Added user %s to team %s\n", u.Username, t.Name)
			}
		}
	},
}

var removeUserFromAllTeamsCmd = &cobra.Command{
	Use:   "removeuser",
	Short: "Remove a user to all teams.",
	PreRun: func(_ *cobra.Command, _ []string) {
		initialize.FullInit()
	},
	Args: cobra.ExactArgs(1),
	Run: func(_ *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("Please provide a username.")
		}

		s := db.NewSession()
		defer s.Close()

		u, err := user.GetUserByUsername(s, args[0])
		if err != nil {
			log.Fatalf("Error getting user: %s", err)
		}

		if u == nil {
			log.Fatalf("User %s does not exist.", args[0])
		}

		teams, err := listAllTeams(s)
		if err != nil {
			log.Fatalf("Error getting teams: %s", err)
		}

		for _, t := range teams {
			tm := &models.TeamMember{
				UserID:   u.ID,
				TeamID:   t.ID,
				Username: u.Username,
				Admin:    false,
			}
			err = tm.Delete(s, u)
			// The delete method does not return an error if the user is not in the team
			// but it does return an error if the team does not exist
			// or if the team only has one member.
			if err != nil {
				// log.Fatalf("Error adding user to team %s: %s", t.Name, err)
				log.Infof("User might not have been removed from %s because: %s", t.Name, err)
			} else {
				log.Infof("Removed user %s from team %s. The function does NOT error if the user is not in the team.\n", u.Username, t.Name)
			}
		}
	},
}
