package actions

import "github.com/gobuffalo/buffalo"

// ShowSystemSetupPage display System Setup Page
func ShowSystemSetupPage(c buffalo.Context) error {
	return c.Render(200, r.HTML("system-setup.html"))
}

func HandleLoadSecurities(c buffalo.Context) error {
	c.Set("securities", LoadSecurities())
	return c.Render(200, r.HTML("securities.html"))
}

func HandleAddSecurity(c buffalo.Context) error {
	request := Securities{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	AddSecurity(request)
	c.Flash().Add("success", "Security created successfully")
	return c.Redirect(302, "/securities")
}

func HandleDeleteSecurity(c buffalo.Context) error {
	request := OnlyID{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteSecurity(request.ID)
	c.Flash().Add("success", "Security deleted successfully")
	return c.Redirect(302, "/securities")
}

func HandleEditSecurity(c buffalo.Context) error {
	request := Securities{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	EditSecurity(request)
	c.Flash().Add("success", "Security edited successfully")
	return c.Redirect(302, "/securities")
}

//------------------------------------------------------------------------------

func HandleLoadAssetClass(c buffalo.Context) error {
	c.Set("assetclasses", LoadAssetClassReports())
	return c.Render(200, r.HTML("assetclass.html"))
}

func HandleAddAssetClass(c buffalo.Context) error {
	request := AssetClass{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	AddAssetClass(request)
	c.Flash().Add("success", "Asset Class created successfully")
	return c.Redirect(302, "/assetclass")
}

func HandleDeleteAssetClass(c buffalo.Context) error {
	request := AssetClass{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	DeleteAssetClass(request.ID)
	c.Flash().Add("success", "Asset Class deleted successfully")
	return c.Redirect(302, "/assetclass")
}

func HandleEditAssetClass(c buffalo.Context) error {
	request := AssetClass{}
	if err := c.Bind(&request); err != nil {
		return err
	}
	EditAssetClass(request)
	c.Flash().Add("success", "Asset class edited successfully")
	return c.Redirect(302, "/assetclass")
}
