package actions

import (
	"github.com/gobuffalo/buffalo"
	v "github.com/gobuffalo/validate"
	"golang.org/x/crypto/bcrypt"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	return c.Render(200, r.HTML("index.html"))
}

//LoginHandler handles the user login
func LoginHandler(c buffalo.Context) error {
	request := &LoginRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	validationErrors := v.Validate(request)
	if len(validationErrors.Errors) == 0 {
		user := GetUserByStaffID(request.StaffID)
		//store authenticated user
		if user.IsEmpty() {
			c.Flash().Add("error", "Invalid credentials.")
			return c.Redirect(302, "/")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err == nil {
			c.Session().Set("auth_id", user.ID)
			return c.Redirect(302, "/dashboard")
		}
		c.Flash().Add("error", "Invalid credentials.")
		return c.Redirect(302, "/")
	}
	c.Flash().Add("error", "Invalid credentials.")
	return c.Redirect(302, "/")
}

//Dashboard takes the user to the dashboard
func Dashboard(c buffalo.Context) error {
	c.Set("overall_summary", LoadOverallSecuritiesSummary())
	c.Set("overview_data", LoadOverviewData())
	c.Set("submitted_data", LoadUploadedPVData())
	return c.Render(200, r.HTML("dashboard.html"))
}

//DisplayQuarterOverviewPage loads the quarterly page reports ...
func DisplayQuarterOverviewPage(c buffalo.Context) error {
	c.Set("clients", LoadNClients(10))
	c.Set("summary", LoadQuarterOverview(c.Param("date")))
	return c.Render(200, r.HTML("quarter-overview.html"))
}

//Logout logs the user out
func Logout(c buffalo.Context) error {
	c.Session().Clear()
	return c.Redirect(302, "/")
}

func LoadPVHistoryPage(c buffalo.Context) error {
	return c.Render(200, r.HTML("pv_history.html"))
}
