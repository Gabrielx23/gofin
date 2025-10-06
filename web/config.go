package web

const (
	RouteLogin    = "/login"
	RouteLogout   = "/logout"
	RouteDashboard = "/dashboard"
	RouteStatic   = "/static/*"
	
	TemplatesDir = "web/templates"
	BaseTemplate = "web/templates/base.html"
	
	SessionTokenCookie = "session_token"
	DatabaseFile = "database.db"
	
	CookiePath = "/"
	CookieMaxAge = 86400
	CookieMaxAgeClear = -1
	
	ContextAccessID = "accessID"
	ContextProjectID = "projectID"
	ContextProjectSlug = "projectSlug"
	
	EmptyString = ""
	EmptySessionTokenError = "empty session token"
	ProjectIDNotFoundError = "Project ID not found"
	ProjectSlugNotFoundError = "Project slug not found"
	AccessIDNotFoundError = "Access ID not found"
	ProjectNotFoundError = "Project not found"
	AccessNotFoundError = "Access not found"
	
	StaticDir = "web/static"
)

func GetTemplatePath(filename string) string {
	return TemplatesDir + "/" + filename
}
