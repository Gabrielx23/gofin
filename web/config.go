package web

const (
	RouteLogin             = "/login"
	RouteLogout            = "/logout"
	RouteDashboard         = "/dashboard"
	RouteCreateTransaction = "/transactions/create"
	RouteCreateAccount     = "/accounts/create"
	RouteStatic            = "/static/*"

	TemplatesDir = "web/templates"
	BaseTemplate = "web/templates/base.html"

	SessionTokenCookie = "session_token"
	DatabaseFile       = "database.db"

	CookiePath        = "/"
	CookieMaxAge      = 86400
	CookieMaxAgeClear = -1

	ContextAccessID    = "accessID"
	ContextProjectID   = "projectID"
	ContextProjectSlug = "projectSlug"

	EmptyString              = ""
	EmptySessionTokenError   = "empty session token"
	ProjectIDNotFoundError   = "Project ID not found"
	ProjectSlugNotFoundError = "Project slug not found"
	AccessIDNotFoundError    = "Access ID not found"
	ProjectNotFoundError     = "Project not found"
	AccessNotFoundError      = "Access not found"

	SuccessTransactionsCreated = "Transactions created successfully!"
	SuccessLoginSuccessful     = "Login successful!"

	SuccessKeyTransactionsCreated = "transactions_created"
	SuccessKeyLoginSuccessful     = "login_successful"

	SuccessQueryParam = "success"

	StaticDir = "web/static"
)
