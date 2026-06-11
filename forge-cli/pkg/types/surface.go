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

// allSurfaceTypesSet is the cached set for AllSurfaceTypesSet.
var allSurfaceTypesSet map[SurfaceType]bool

// AllSurfaceTypesSet returns a set of all defined SurfaceType constants
// for O(1) membership checks. The returned map must not be mutated.
func AllSurfaceTypesSet() map[SurfaceType]bool {
	if allSurfaceTypesSet == nil {
		allSurfaceTypesSet = make(map[SurfaceType]bool, len(AllSurfaceTypes()))
		for _, st := range AllSurfaceTypes() {
			allSurfaceTypesSet[st] = true
		}
	}
	return allSurfaceTypesSet
}
