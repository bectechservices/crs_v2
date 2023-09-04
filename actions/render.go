package actions

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/jinzhu/now"
	"github.com/leekchan/accounting"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/gobuffalo/tags"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	r = render.New(render.Options{
		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			"csrf":                                 csrfHelper,
			"is_admin":                             isAdmin,
			"is_manager":                           isManager,
			"is_sec_maker":                         isSecMaker,
			"is_sec_checker":                       isSecChecker,
			"is_npra_maker":                        isNPRAMaker,
			"is_npra_checker":                      isNPRAChecker,
			"is_trustee_maker":                     isTrusteeMaker,
			"is_trustee_checker":                   isTrusteeChecker,
			"is_nav_maker":                         isNAVMaker,
			"is_nav_checker":                       isNAVChecker,
			"can_do":                               canDo,
			"format_with_comma":                    formatWithComma,
			"format_with_comma_decimal":            formatWithComma2,
			"quarter_name_from_date":               quarterNameFromDate,
			"current_time":                         currentTime,
			"format_date_time":                     formatTimeDate,
			"money_null_float":                     formatNullFloatWithComma,
			"first_day_of_month":                   firstDayOfMonth,
			"last_day_of_month":                    lastDayOfMonth,
			"last_sec_license_renewal_date":        lastSecLicenseRenewalDate,
			"last_day_of_last_quarter":             lastDayOfLastQuarter,
			"last_day_of_next_month":               lastDayOfNextMonth,
			"handle_variance":                      handleVariance,
			"is_alarming_variance":                 isAlarmingVariance,
			"remove_spaces":                        removeSpacesFromString,
			"first_data_of_txn_vol_by_asset_class": firstDataOfAssetClassData,
			"format_date":                          formatDate,
			"first_basis_points_basis_point":       firstBasisPointsBasisPoint,
			"is_empty_time_pointer":                isEmptyTimePointer,
			"parse_holiday":                        parseHoliday,
			"activities_quater":                    parse_quarter_for_activities,
		},
	})
}

func csrfHelper(ctx plush.HelperContext) (template.HTML, error) {
	tok, ok := ctx.Value("authenticity_token").(string)
	if !ok {
		return "", fmt.Errorf("expected CSRF token got %T", ctx.Value("authenticity_token"))
	}
	t := tags.New("input", tags.Options{
		"value": tok,
		"type":  "hidden",
		"name":  "authenticity_token",
	})
	return t.HTML(), nil
}

func isAdmin(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsAdmin()
}

func isManager(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsManager()
}

func isSecMaker(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsSecMaker()
}

func isSecChecker(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsSecChecker()
}

func isNPRAMaker(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsNPRAMaker()
}

func isNPRAChecker(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsNPRAChecker()
}

func isTrusteeMaker(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsTrusteeMaker()
}

func isTrusteeChecker(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsTrusteeChecker()
}

func isNAVMaker(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsNAVMaker()
}

func isNAVChecker(help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.IsNAVChecker()
}

func canDo(permission string, help plush.HelperContext) bool {
	user := help.Value("auth_user").(User)
	return user.Can(permission)
}

func formatWithComma(amount interface{}, dp int, help plush.HelperContext) string {
	return FormatWithComma(amount, dp)
}

func quarterNameFromDate(date time.Time) string {
	if date.IsZero() {
		date = time.Now()
	}
	month := int(date.Month())
	if month >= 1 && month <= 3 {
		return "Fourth Quarter"
	} else if month >= 4 && month <= 6 {
		return "First Quarter"
	} else if month >= 7 && month <= 9 {
		return "Second Quarter"
	} else {
		return "Third Quarter"
	}
}

func currentTime(layout string) string {
	return time.Now().Format(layout)
}

func firstDayOfMonth(layout string) string {
	return now.New(time.Now()).BeginningOfMonth().Format(layout)
}

func lastDayOfMonth(layout string) string {
	return now.New(time.Now()).EndOfMonth().Format(layout)
}

func lastDayOfLastQuarter(layout string) string {
	currentQuarter, _ := time.Parse("2006-01-02", GetQuarterFormalDate())
	return currentQuarter.Format(layout)
}

func lastDayOfNextMonth(layout string) string {
	return now.New(time.Now().AddDate(0, 1, 0)).EndOfMonth().Format(layout)
}

func handleVariance(number interface{}) string {
	return replaceNegativeWithBraces(number)
}

func isAlarmingVariance(number float32) bool {
	return number <= -5 || number >= 5
}

func removeSpacesFromString(value string) string {
	return strings.ReplaceAll(value, " ", "")
}

func firstDataOfAssetClassData(data []TradeVolumeByAssetClassSummary) []TradeVolumeByAssetClass {
	return data[0].Data
}

func formatDate(date time.Time, layout string) string {
	return date.Format("2006-Jan-02")
}

func firstBasisPointsBasisPoint(data []ClientBasisPoints) float64 {
	return data[0].BasisPoints
}

func isEmptyTimePointer(date *time.Time) bool {
	return date == nil
}

func lastSecLicenseRenewalDate() string {
	year := time.Now().Year()
	expiryDate, _ := time.Parse("2006-01-02", fmt.Sprintf("%d-06-30", year))
	if time.Now().After(expiryDate) {
		return fmt.Sprintf("1st July %d", year+1)
	}
	return fmt.Sprintf("1st July %d", year)
}

func parseHoliday(date string) string {
	parsedDate, _ := time.Parse("02-01", date)
	return parsedDate.Format("02 January")
}

func parse_quarter_for_activities(date time.Time) string {
	month := int(date.Month())
	year := date.Year()
	if month == 3 {
		return fmt.Sprintf(`Q1 %d`, year)
	} else if month == 6 {
		return fmt.Sprintf(`Q2 %d`, year)
	} else if month == 9 {
		return fmt.Sprintf(`Q3 %d`, year)
	} else {
		return fmt.Sprintf(`Q4 %d`, year)
	}
}

func formatWithComma2(amount float64) string {
	p := message.NewPrinter(language.English)
	withCommaThousandSep := p.Sprintf("%.2f", amount)
	return withCommaThousandSep
}

func formatNullFloatWithComma(amount float64) string {
	ac := accounting.Accounting{Symbol: "¢", Precision: 2}
	withCommaThousandSep := ac.FormatMoney(amount)
	return withCommaThousandSep
}

func formatTimeDate(date, format string) string {
	parsedDate, _ := time.Parse("20060102", date)
	return parsedDate.Format(format)
}
