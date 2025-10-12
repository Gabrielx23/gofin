package web

import (
	"net/http"
	"net/url"

	"gofin/web"
)

const TemplatesDir = "web/templates"

func RedirectWithSuccess(w http.ResponseWriter, r *http.Request, redirectPath, successMessage string) {
	u, _ := url.Parse(redirectPath)
	q := u.Query()
	q.Set(web.SuccessQueryParam, successMessage)
	u.RawQuery = q.Encode()
	http.Redirect(w, r, u.String(), http.StatusSeeOther)
}

func GetTemplatePath(filename string) string {
	return TemplatesDir + "/" + filename
}

func RedirectToProjectLogin(w http.ResponseWriter, r *http.Request, projectSlug string) {
	http.Redirect(w, r, "/"+projectSlug+web.RouteLogin, http.StatusSeeOther)
}

func RedirectToProjectDashboard(w http.ResponseWriter, r *http.Request, projectSlug string) {
	http.Redirect(w, r, "/"+projectSlug+web.RouteDashboard, http.StatusSeeOther)
}
