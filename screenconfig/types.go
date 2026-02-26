package screenconfig

// Pattern enumera los patterns de pantalla soportados
type Pattern string

// Pattern constants define the supported screen layout patterns
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

// ScreenType constants define how a screen relates to a resource
const (
	ScreenTypeList      ScreenType = "list"
	ScreenTypeDetail    ScreenType = "detail"
	ScreenTypeCreate    ScreenType = "create"
	ScreenTypeEdit      ScreenType = "edit"
	ScreenTypeDashboard ScreenType = "dashboard"
	ScreenTypeSettings  ScreenType = "settings"
)

// Platform identifica la plataforma del cliente para aplicar overrides de UI
type Platform string

// Platform constants identify supported client platforms for UI overrides
const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
	PlatformMobile  Platform = "mobile" // fallback generico para mobile
	PlatformDesktop Platform = "desktop"
	PlatformWeb     Platform = "web"
)
