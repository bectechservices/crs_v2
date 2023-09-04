package actions

import "github.com/gobuffalo/buffalo"

//LoadUsersPage displays the users page
func LoadUsersPage(c buffalo.Context) error {
	c.Set("roles", LoadRolesFromDB())
	c.Set("users", LoadUsersFromDB())
	return c.Render(200, r.HTML("users.html"))
}

//HandleAddUser handles adding users to the system
func HandleAddUser(c buffalo.Context) error {
	request := UserAddRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	CreateNewUser(request)
	c.Flash().Add("success", "User created successfully")
	return c.Redirect(302, "/users")
}

//HandleTriggerPasswordReset handles the request to reset the users password
func HandleTriggerPasswordReset(c buffalo.Context) error {
	request := OnlyUserIDRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	ResetUserPassword(request.UserID, request.Email)
	c.Flash().Add("success", "Password reset successful")
	return c.Redirect(302, "/users")
}

//HandleUserDelete handles the request which deletes a user
func HandleUserDelete(c buffalo.Context) error {
	request := OnlyUserIDRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteUser(request)
	c.Flash().Add("success", "User deleted successfully")
	return c.Redirect(302, "/users")
}

//HandleEditUser handles the request to edit the users account
func HandleEditUser(c buffalo.Context) error {
	request := EditUserRequest{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	EditUser(request)
	c.Flash().Add("success", "User edited successfully")
	return c.Redirect(302, "/users")
}

func HandlePasswordReset(c buffalo.Context) error {
	request := OnlyStaffID{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	user := GetUserByStaffID(request.StaffID)
	if user.ID > 0 {
		ResetUserPassword(user.ID, user.Email)
		c.Flash().Add("success", "Password reset successful")
	} else {
		c.Flash().Add("error", "Invalid staff id")
	}
	return c.Redirect(302, "/")
}
