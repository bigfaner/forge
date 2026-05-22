// Package cmd provides the CLI commands for the forge CLI tool.
//
// Error and output types are defined in the base sub-package and
// re-exported here for backward compatibility with existing callers.
package cmd

import "forge-cli/internal/cmd/base"

// ErrorCode is re-exported from base package for backward compatibility.
type ErrorCode = base.ErrorCode

// AIError is re-exported from base package for backward compatibility.
type AIError = base.AIError

// Re-export constants from base package.
const (
	ErrNoProject            = base.ErrNoProject
	ErrNoFeature            = base.ErrNoFeature
	ErrInvalidInput         = base.ErrInvalidInput
	ErrNotFound             = base.ErrNotFound
	ErrConflict             = base.ErrConflict
	ErrValidation           = base.ErrValidation
	ErrInvalidTransition    = base.ErrInvalidTransition
	ErrInvalidPath          = base.ErrInvalidPath
	ErrEvalParseFailure     = base.ErrEvalParseFailure
	ErrContractUnverifiable = base.ErrContractUnverifiable
)

// Re-export functions from base package for backward compatibility.
var (
	Exit                       = base.Exit
	NewAIError                 = base.NewAIError
	ErrProjectNotFound         = base.ErrProjectNotFound
	ErrFeatureNotSet           = base.ErrFeatureNotSet
	ErrTaskNotFound            = base.ErrTaskNotFound
	ErrNoInput                 = base.ErrNoInput
	ErrInvalidJSON             = base.ErrInvalidJSON
	ErrFileNotFound            = base.ErrFileNotFound
	ErrNoPendingTasks          = base.ErrNoPendingTasks
	ErrDependenciesNotMet      = base.ErrDependenciesNotMet
	ErrDataIntegrity           = base.ErrDataIntegrity
	ErrInvalidStatus           = base.ErrInvalidStatus
	ErrMissingFields           = base.ErrMissingFields
	WarnMissingFields          = base.WarnMissingFields
	ErrFeatureNotFound         = base.ErrFeatureNotFound
	ErrNoTestEvidence          = base.ErrNoTestEvidence
	ErrUnmetAcceptanceCriteria = base.ErrUnmetAcceptanceCriteria
	ErrTaskIDConflict          = base.ErrTaskIDConflict
	ErrInvalidDependency       = base.ErrInvalidDependency
	NewErrInvalidTransition    = base.NewErrInvalidTransition
	NewErrInvalidPath          = base.NewErrInvalidPath
	NewErrEvalParseFailure     = base.NewErrEvalParseFailure
	NewErrContractUnverifiable = base.NewErrContractUnverifiable
	ErrNotGitRepository        = base.ErrNotGitRepository
	ErrNotInsideWorktree       = base.ErrNotInsideWorktree
	ErrRefusingDefaultBranch   = base.ErrRefusingDefaultBranch
	ErrSlugRequired            = base.ErrSlugRequired
	ErrSourceBranchNotFound    = base.ErrSourceBranchNotFound
)
