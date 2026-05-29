package types

// SurfaceType represents the interface surface of a feature (web, api, cli, tui, mobile).
type SurfaceType string

// Interface surface type constants.
const (
	SurfaceWeb    SurfaceType = "web"
	SurfaceAPI    SurfaceType = "api"
	SurfaceCLI    SurfaceType = "cli"
	SurfaceTUI    SurfaceType = "tui"
	SurfaceMobile SurfaceType = "mobile"
)

// AllSurfaceTypes returns all defined SurfaceType constants.
func AllSurfaceTypes() []SurfaceType {
	return []SurfaceType{
		SurfaceWeb,
		SurfaceAPI,
		SurfaceCLI,
		SurfaceTUI,
		SurfaceMobile,
	}
}
