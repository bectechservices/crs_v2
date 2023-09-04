package actions

import (
	"github.com/gobuffalo/buffalo"
)

//showProfilePage displays the password reset page
func showProfilePage(c buffalo.Context) error {
	c.Set("user", AuthUser(c))
	return c.Render(200, r.HTML("profile-page.html"))
}

func HandleUserProfileUpdate(c buffalo.Context) error {
	request := UserProfileUpdateRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	UpdateUserProfile(AuthID(c), request)
	c.Flash().Add("success", "Profile updated successfully")
	return c.Redirect(302, "/dashboard")
}
