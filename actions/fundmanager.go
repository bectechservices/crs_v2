package actions

import (
	"fmt"
	"github.com/gobuffalo/buffalo"
	"strconv"
)

//ShowFundManager displays the fund managers page
func ShowFundManager(c buffalo.Context) error {
	c.Set("fund_managers", LoadFundManagers())
	return c.Render(200, r.HTML("fundmanagers.html"))
}

//ShowFundManagerEdit displays the fund managers Edit Page
func ShowFundManagerEdit(c buffalo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}
	c.Set("manager", GetFundManagerByID(id))
	return c.Render(200, r.HTML("fundmanager-edit.html"))
}

//ShowAddFundManager displays the fund managers add Page
func ShowAddFundManager(c buffalo.Context) error {
	return c.Render(200, r.HTML("fundmanager-add.html"))
}

//HandleAddFundManager handles the request to add a fundmanager
func HandleAddFundManager(c buffalo.Context) error {
	request := &FundManagerAddRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	//TODO: show alerts
	AddFundManagerToDB(*request)
	return c.Redirect(302, "/fundmanager-add")
}

//HandleEditFundManager handles the request to edit a fund manager
func HandleEditFundManager(c buffalo.Context) error {
	request := &FundManagerEditRequest{}
	if err := c.Bind(request); err != nil {
		return err
	}
	//TODO: show alerts
	EditFundManager(*request)
	return c.Redirect(302, fmt.Sprintf("/fundmanager-edit?id=%d", request.ID))

}
