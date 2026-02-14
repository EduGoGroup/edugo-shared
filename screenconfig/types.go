package screenconfig

// Pattern enumera los patterns de pantalla soportados
type Pattern string

const (
	PatternLogin        Pattern = "login"
	PatternForm         Pattern = "form"
	PatternList         Pattern = "list"
	PatternDashboard    Pattern = "dashboard"
	PatternSettings     Pattern = "settings"
	PatternDetail       Pattern = "detail"
	PatternSearch       Pattern = "search"
	PatternProfile      Pattern = "profile"
	PatternModal        Pattern = "modal"
	PatternNotification Pattern = "notification"
	PatternOnboarding   Pattern = "onboarding"
	PatternEmptyState   Pattern = "empty-state"
)

// ScreenType define como una pantalla se relaciona con un recurso
type ScreenType string

const (
	ScreenTypeList      ScreenType = "list"
	ScreenTypeDetail    ScreenType = "detail"
	ScreenTypeCreate    ScreenType = "create"
	ScreenTypeEdit      ScreenType = "edit"
	ScreenTypeDashboard ScreenType = "dashboard"
	ScreenTypeSettings  ScreenType = "settings"
)

// ActionType enumera las acciones estandar
type ActionType string

const (
	ActionNavigate     ActionType = "NAVIGATE"
	ActionNavigateBack ActionType = "NAVIGATE_BACK"
	ActionAPICall      ActionType = "API_CALL"
	ActionSubmitForm   ActionType = "SUBMIT_FORM"
	ActionRefresh      ActionType = "REFRESH"
	ActionConfirm      ActionType = "CONFIRM"
	ActionLogout       ActionType = "LOGOUT"
)
