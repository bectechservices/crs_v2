package actions

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// SetCurrentUser attempts to find a user based on the auth_id in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("auth_id"); uid != nil {
			switch id := uid.(type) {
			case int:
				user := GetUserByID(id)
				if user.IsEmpty() {
					return errors.WithStack(errors.New("User not found"))
				}
				c.Set("auth_user", user)
			}
		}
		return next(c)
	}
}

//SetConstantsInContext sets all the app's constants in the context
func SetConstantsInContext(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		c.Set("UploadPV", UploadPV)
		c.Set("InputClientInformation", InputClientInformation)
		c.Set("InputSECInformation", InputSECInformation)
		c.Set("InputClientBankDetails", InputClientBankDetails)
		c.Set("ViewPVDetails", ViewPVDetails)
		c.Set("ViewClientInformation", ViewClientInformation)
		c.Set("ApproveSECMakerRequest", ApproveSECMakerRequest)
		c.Set("ApproveNPRAMakerRequest", ApproveNPRAMakerRequest)
		c.Set("ApproveTrusteeMakerRequest", ApproveTrusteeMakerRequest)
		c.Set("ApproveNAVMakerRequest", ApproveNAVMakerRequest)
		c.Set("AccessClientBasedAnalytics", AccessClientBasedAnalytics)
		c.Set("AccessOverallBasedAnalytics", AccessOverallBasedAnalytics)
		c.Set("PVReportOutstanding", PVReportOutstanding)
		c.Set("SECRequests", SECRequests)
		c.Set("NPRARequests", NPRARequests)
		c.Set("ClientLetters", ClientLetters)
		c.Set("MissedLiabilities", MissedLiabilities)
		c.Set("NPRAReportDelays", NPRAReportDelays)
		c.Set("SECReportDelays", SECReportDelays)
		c.Set("ClientReportDelays", ClientReportDelays)
		c.Set("InvoiceDelay", InvoiceDelay)
		c.Set("InputNPRADetails", InputNPRADetails)
		c.Set("InputTrusteeDetails", InputTrusteeDetails)
		c.Set("AccessMailPage", AccessMailPage)
		c.Set("ViewSentMails", ViewSentMails)
		c.Set("ModifyMailTemplate", ModifyMailTemplate)

		c.Set("AdminRole", AdminRole)
		c.Set("ManagerRole", ManagerRole)
		c.Set("SECMakerRole", SECMakerRole)
		c.Set("SECCheckerRole", SECCheckerRole)
		c.Set("NPRAMakerRole", NPRAMakerRole)
		c.Set("NPRACheckerRole", NPRACheckerRole)
		c.Set("TrusteeMakerRole", TrusteeMakerRole)
		c.Set("TrusteeCheckerRole", TrusteeCheckerRole)
		c.Set("NAVMakerRole", NAVMakerRole)
		c.Set("NAVCheckerRole", NAVCheckerRole)

		c.Set("UploadScheme", UploadScheme)
		c.Set("UploadSecPV", UploadSecPV)
		c.Set("UploadOutstandingFDReport", UploadOutstandingFDReport)
		c.Set("UploadGOGMaturities", UploadGOGMaturities)
		c.Set("UploadTransactionVolumes", UploadTransactionVolumes)
		c.Set("UploadedMaturedSecurities", UploadedMaturedSecurities)
		c.Set("current_year_constant", time.Now().Year())
		c.Set("current_quarter_constant", GetQuarterNumber(time.Now()))
		c.Set("current_quarter_fancy_date", MakeQuarterFancyDate(time.Now()))
		c.Set("current_quarter_ordinal_date", GetShortQuarterDate())
		c.Set("previous_quarter_ordinal_date", GetShortPreviousQuarterDate())
		return next(c)
	}
}

// RequiresAuthorization require a user be logged in before accessing a route
func RequiresAuthorization(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if id := c.Session().Get("auth_id"); id == nil {
			c.Flash().Add("error", "You must be logged in to proceed.")
			return c.Redirect(302, "/")
		}
		return next(c)
	}
}

//RedirectIfAuthenticated redirects to the dashboard if the user is authenticated
func RedirectIfAuthenticated(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if id := c.Session().Get("auth_id"); id != nil {
			if c.Request().URL.Path == "/" {
				return c.Redirect(302, "/dashboard")
			}
		}
		return next(c)
	}
}

// MustBeAdmin is a middleware for pages only the admin can view
func MustBeAdmin(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		routes := []string{"users", "trigger-password-reset"}
		routeInRoutes := func(route string) bool {
			for _, val := range routes {
				if route == val {
					return true
				}
			}
			return false
		}
		if routeInRoutes(strings.Trim(c.Request().URL.Path, "/")) {
			user := c.Value("auth_user").(User)
			if user.IsAdmin() {
				return next(c)
			}
			return c.Redirect(302, "/dashboard")
		}
		return next(c)
	}
}

// MustBeSec is a middleware for pages only the sec users, managers and the admin can view
func MustBeSec(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		routes := []string{"custodian", "sec", "upload-governance", "upload-scheme-details", "upload-other-information", "upload-official-remarks", "sec-quarterly-report", "fetch-scheme-details", "secreport", "sec-report-quarterly-performance", "sec-report-contributions", "local-variance", "foreign-variance", "client-letters", "send-client-letters"}
		routeInRoutes := func(route string) bool {
			for _, val := range routes {
				if route == val {
					return true
				}
			}
			return false
		}
		if routeInRoutes(strings.Trim(c.Request().URL.Path, "/")) {
			user := c.Value("auth_user").(User)
			if user.IsSecChecker() || user.IsSecMaker() || user.IsAdmin() || user.IsManager() {
				return next(c)
			}
			return c.Redirect(302, "/dashboard")
		}
		return next(c)
	}
}

// MustBeNPRA is a middleware for pages only the npra users, managers and the admin can view
func MustBeNPRA(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		routes := []string{"npra", "npra-dashboard", "npra-report"}
		routeInRoutes := func(route string) bool {
			for _, val := range routes {
				if route == val {
					return true
				}
			}
			return false
		}
		if routeInRoutes(strings.Trim(c.Request().URL.Path, "/")) {
			user := c.Value("auth_user").(User)
			if user.IsNPRAMaker() || user.IsNPRAChecker() || user.IsAdmin() || user.IsManager() {
				return next(c)
			}
			return c.Redirect(302, "/dashboard")
		}
		return next(c)
	}
}

// MustBeTrustee is a middleware for pages only the trustee users, managers and the admin can view
func MustBeTrustee(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		routes := []string{"trustee-dashboard", "trustee-report"}
		routeInRoutes := func(route string) bool {
			for _, val := range routes {
				if route == val {
					return true
				}
			}
			return false
		}
		if routeInRoutes(strings.Trim(c.Request().URL.Path, "/")) {
			user := c.Value("auth_user").(User)
			if user.IsTrusteeChecker() || user.IsTrusteeMaker() || user.IsAdmin() || user.IsManager() {
				return next(c)
			}
			return c.Redirect(302, "/dashboard")
		}
		return next(c)
	}
}

// MustBeNAV is a middleware for pages only the nav users, managers and the admin can view
func MustBeNAV(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		routes := []string{}
		routeInRoutes := func(route string) bool {
			for _, val := range routes {
				if route == val {
					return true
				}
			}
			return false
		}
		if routeInRoutes(strings.Trim(c.Request().URL.Path, "/")) {
			user := c.Value("auth_user").(User)
			if user.IsNAVChecker() || user.IsNAVMaker() || user.IsAdmin() || user.IsManager() {
				return next(c)
			}
			return c.Redirect(302, "/dashboard")
		}
		return next(c)
	}
}

//LoadSystemAlerts loads all system alerts
func LoadSystemAlerts(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		routes := []string{"dashboard", "sec", "npra-dashboard"}
		routeInRoutes := func(route string) bool {
			for _, val := range routes {
				if route == val {
					return true
				}
			}
			return false
		}
		if routeInRoutes(strings.Trim(c.Request().URL.Path, "/")) {
			settings := LoadAlertSettings()
			alerts := make([]Alert, 0)
			end, _ := time.Parse("2006-01-02", GetQuarterFormalDate())
			daysToEndOfQuarter := int(end.Sub(time.Now()).Hours() / 24)
			if daysToEndOfQuarter <= settings.DaysToReportPV {
				due := "Overdue"
				if daysToEndOfQuarter > 0 {
					due = fmt.Sprintf("In %d days", daysToEndOfQuarter)
				}
				data := LoadUploadedPVData()
				alerts = append(alerts, Alert{
					Type:   "PV",
					Number: data.Pending,
					Due:    due,
				})
			}
			c.Set("system_alerts", alerts)
		}
		return next(c)
	}
}

//MustResetPassword is a middleware for pages only the admin can view
func MustResetPassword(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		user := c.Value("auth_user").(User)
		if user.MustResetPassword {
			return c.Redirect(302, "/profile")
		}
		return next(c)
	}
}
