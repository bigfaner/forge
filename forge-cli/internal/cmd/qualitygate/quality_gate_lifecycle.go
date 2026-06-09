package qualitygate

import (
	"time"

	"forge-cli/pkg/feature"
	"forge-cli/pkg/forgeconfig"
	"forge-cli/pkg/forgelog"
	"forge-cli/pkg/just"
	"forge-cli/pkg/serverprobe"
	"forge-cli/pkg/testrunner"
	"forge-cli/pkg/types"
)

// runTestRegression runs the full test regression suite when a justfile with
// a test recipe is present. When surfaces are configured in .forge/config.yaml,
// it orchestrates per-surface lifecycle (dev->probe->test->teardown for web/api/mobile;
// test->teardown for cli/tui). Falls back to legacy behavior when no surfaces configured.
// Returns an error when a gate failure is detected, nil otherwise.
func runTestRegression(projectRoot, featureSlug string) error {
	if !just.HasJustfile(projectRoot) || !just.HasRecipe(projectRoot, "test") {
		return nil
	}

	// Detect surface types from config.
	surfaces, _ := forgeconfig.ReadSurfaces(projectRoot)
	surfaceTypes := forgeconfig.SurfaceTypes(surfaces)

	if len(surfaceTypes) == 0 {
		// No surfaces configured — fall back to legacy behavior.
		return runTestRegressionLegacy(projectRoot, featureSlug)
	}

	// Surface-aware orchestration: run lifecycle per surface type.
	return runTestRegressionSurface(projectRoot, featureSlug, surfaceTypes)
}

// runTestRegressionLegacy is the pre-surface-aware test regression logic.
// Runs test-setup (optional), serverprobe health check, then just test.
func runTestRegressionLegacy(projectRoot, featureSlug string) error {
	// Optional setup step — skip regression on failure.
	if just.HasRecipe(projectRoot, "test-setup") {
		forgelog.Info("--- Ensuring test dependencies (just test-setup) ---\n")
		setupOutput, setupSuccess := just.RunCapture(projectRoot, "just", "test-setup")
		if !setupSuccess {
			forgelog.Warn("WARNING: test-setup failed; skipping test regression\n")
			forgelog.Info("  To retry manually: just test-setup && just test\n")
			if setupOutput != "" {
				if err := testrunner.WriteRegressionRawOutput(projectRoot, "=== test-setup failure ===\n"+setupOutput); err != nil {
					forgelog.Warn("WARNING: failed to write setup output: %v\n", err)
				} else {
					forgelog.Info("  Setup output saved to " + feature.TestResultsDir + "/" + feature.TestOutputFileName + "\n")
				}
			}
			return nil
		}
	}

	// Health check — skip regression if servers aren't ready.
	if !serverprobe.ProbeServers(projectRoot, "") {
		forgelog.Warn("WARNING: server health check failed; skipping test regression\n")
		forgelog.Info("  Start dev server and retry: just dev && just test\n")
		return nil
	}

	// Run the regression suite.
	forgelog.Info("--- Running full test regression (just test) ---\n")
	regressionOutput, regSuccess := just.RunCapture(projectRoot, "just", "test")
	if !regSuccess {
		forgelog.Error("ERROR: test regression failed\n")
		errorDocPath := feature.TestResultsDir + "/" + feature.TestOutputFileName
		if regressionOutput != "" {
			if err := testrunner.WriteRegressionRawOutput(projectRoot, regressionOutput); err != nil {
				forgelog.Warn("WARNING: failed to write raw-output.txt: %v\n", err)
			}
		}
		fixID, fixErr := addRegressionFixTasks(projectRoot, featureSlug, regressionOutput, errorDocPath)
		if fixErr != nil {
			forgelog.Warn("WARNING: %v\n", fixErr)
		}
		return HandleGateFailure("test", errorDocPath, fixID, just.ExtractConciseError(regressionOutput, conciseErrorMaxLines), true)
	}
	return nil
}

// runTestRegressionSurface orchestrates per-surface-type lifecycle sequences.
// For each unique surface type, runs the appropriate sequence:
//   - web/api: dev -> probe -> test -> teardown (full lifecycle)
//   - mobile: dev -> probe -> test-setup -> test -> teardown (full lifecycle with mobile setup)
//   - cli/tui: test -> teardown (simplified)
//
// Surfaces of the same type share a single lifecycle (dev/probe run once per type).
// Teardown is mandatory regardless of prior step success/failure.
func runTestRegressionSurface(projectRoot, featureSlug string, surfaceTypes []string) error {
	var lastErr error
	for _, surfaceType := range surfaceTypes {
		forgelog.Info("--- Running surface orchestration for %s ---\n", surfaceType)
		result := runSurfaceLifecycle(projectRoot, surfaceType)
		if !result.success {
			errorDocPath := feature.TestResultsDir + "/" + feature.TestOutputFileName
			if result.output != "" {
				if err := testrunner.WriteRegressionRawOutput(projectRoot, result.output); err != nil {
					forgelog.Warn("WARNING: failed to write raw-output.txt: %v\n", err)
				}
			}
			fixID, fixErr := addRegressionFixTasks(projectRoot, featureSlug, result.output, errorDocPath)
			if fixErr != nil {
				forgelog.Warn("WARNING: %v\n", fixErr)
			}
			lastErr = HandleGateFailure("test", errorDocPath, fixID, just.ExtractConciseError(result.output, conciseErrorMaxLines), true)
		}
	}
	return lastErr
}

// lifecycleResult holds the result of a surface lifecycle execution.
type lifecycleResult struct {
	success bool
	output  string
}

// needsFullLifecycle returns true for surface types that require dev->probe->test->teardown.
// cli and tui surfaces use the simplified test->teardown sequence.
func needsFullLifecycle(surfaceType types.SurfaceType) bool {
	return surfaceType == types.SurfaceWeb || surfaceType == types.SurfaceAPI || surfaceType == types.SurfaceMobile
}

// resolveRecipe resolves a recipe name using surface-type prefix.
// When surfaceType is non-empty, it returns "<surfaceType>-<recipe>" if that
// prefixed recipe exists, or empty string if it does not (no fallback to generic).
// When surfaceType is empty, it returns the generic recipe name if it exists.
func resolveRecipe(projectRoot, surfaceType, genericRecipe string) string {
	if surfaceType != "" {
		specificRecipe := surfaceType + "-" + genericRecipe
		if just.HasRecipe(projectRoot, specificRecipe) {
			return specificRecipe
		}
		// Prefixed recipe not found — no fallback to generic.
		return ""
	}
	// No surface type — use generic recipe name.
	if just.HasRecipe(projectRoot, genericRecipe) {
		return genericRecipe
	}
	return ""
}

// runSurfaceLifecycle executes the per-surface lifecycle sequence.
// For web/api: dev -> probe -> test -> teardown
// For mobile: dev -> probe -> mobile-test-setup -> test -> teardown
// For cli/tui: test -> teardown
// Teardown always executes (via defer-like pattern).
func runSurfaceLifecycle(projectRoot, surfaceType string) lifecycleResult {
	full := needsFullLifecycle(types.SurfaceType(surfaceType))

	// Phase 1: Dev (full lifecycle only)
	if full {
		recipe := resolveRecipe(projectRoot, surfaceType, "dev")
		if recipe != "" {
			forgelog.Info("  Starting dev server (just %s)...\n", recipe)
			output, success := just.RunCapture(projectRoot, "just", recipe)
			if !success {
				forgelog.Error("  ERROR: dev failed (just %s)\n", recipe)
				runTeardown(projectRoot, surfaceType)
				return lifecycleResult{success: false, output: output}
			}
		}
	}

	// Phase 2: Probe (full lifecycle only)
	if full {
		probeRecipe := resolveRecipe(projectRoot, surfaceType, "probe")
		if !probeWithRetry(projectRoot, probeRecipe, maxProbeRetries, probeRetryInterval) {
			forgelog.Error("  ERROR: probe failed after retries\n")
			runTeardown(projectRoot, surfaceType)
			return lifecycleResult{success: false, output: "probe failed: server not responding after 3 retries"}
		}
	}

	// Phase 2b: Mobile test setup (mobile surfaces only)
	if types.SurfaceType(surfaceType) == types.SurfaceMobile {
		setupRecipe := resolveRecipe(projectRoot, surfaceType, "test-setup")
		if setupRecipe != "" {
			forgelog.Info("  Running mobile test setup (just %s)...\n", setupRecipe)
			output, success := just.RunCapture(projectRoot, "just", setupRecipe)
			if !success {
				forgelog.Error("  ERROR: mobile-test-setup failed (just %s)\n", setupRecipe)
				runTeardown(projectRoot, surfaceType)
				return lifecycleResult{success: false, output: output}
			}
		}
	}

	// Phase 3: Test
	var result lifecycleResult
	testRecipe := resolveRecipe(projectRoot, surfaceType, "test")
	if testRecipe != "" {
		forgelog.Info("  Running tests (just %s)...\n", testRecipe)
		output, success := just.RunCapture(projectRoot, "just", testRecipe)
		result = lifecycleResult{success: success, output: output}
		if !success {
			forgelog.Error("  ERROR: test failed\n")
		}
	} else {
		result = lifecycleResult{success: true}
	}

	// Phase 4: Teardown (always)
	runTeardown(projectRoot, surfaceType)

	return result
}

// runTeardown executes the teardown recipe for a surface type.
// Errors are logged but never fail the lifecycle — teardown is best-effort cleanup.
func runTeardown(projectRoot, surfaceType string) {
	recipe := resolveRecipe(projectRoot, surfaceType, "teardown")
	if recipe != "" {
		forgelog.Info("  Running teardown (just %s)...\n", recipe)
		output, success := just.RunCapture(projectRoot, "just", recipe)
		if !success {
			forgelog.Warn("  WARNING: teardown failed (just %s)\n", recipe)
			if output != "" {
				forgelog.Info("  %s\n", just.ExtractConciseError(output, 3))
			}
		}
	}
}

// probeWithRetry runs the probe recipe with the specified number of retries.
// Returns true if the probe succeeds within the retry count.
// Returns true (skip) if the probe recipe doesn't exist.
// interval is the delay between retries (0 for no delay, useful in tests).
func probeWithRetry(projectRoot, probeRecipe string, maxRetries int, interval time.Duration) bool {
	if probeRecipe == "" {
		return true // no probe recipe — skip
	}

	// Verify the recipe actually exists before retrying.
	if !just.HasRecipe(projectRoot, probeRecipe) {
		return true // recipe not found — skip
	}

	for attempt := range maxRetries {
		if attempt > 0 && interval > 0 {
			forgelog.Info("  Probe retry %d/%d (waiting %v)...\n", attempt+1, maxRetries, interval)
			time.Sleep(interval)
		}
		forgelog.Info("  Probing (just %s) attempt %d/%d...\n", probeRecipe, attempt+1, maxRetries)
		_, success := just.RunCapture(projectRoot, "just", probeRecipe)
		if success {
			forgelog.Info("  Probe succeeded\n")
			return true
		}
	}
	return false
}
