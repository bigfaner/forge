package task

import (
	"fmt"
	"strings"

	"forge-cli/pkg/types"
)

// Dependency Resolver functions

// ResolveLastRunTest returns the ID of the last task in the run-test expansion chain.
var ResolveLastRunTest DepResolveFunc = func(ctx *GenContext) []string {
	if len(ctx.RunTestChain) == 0 {
		return nil
	}
	return []string{ctx.RunTestChain[len(ctx.RunTestChain)-1]}
}

// ResolveUpstream returns the IDs of the immediately preceding generated node(s).
var ResolveUpstream DepResolveFunc = func(ctx *GenContext) []string {
	if len(ctx.UpstreamIDs) == 0 {
		return nil
	}
	return ctx.UpstreamIDs
}

// ResolveDocTasks returns the IDs of all doc-category business tasks.
var ResolveDocTasks DepResolveFunc = func(ctx *GenContext) []string {
	var ids []string
	for _, t := range ctx.BusinessTasks {
		if CategoryForType(t.Type) == CategoryDoc {
			ids = append(ids, t.ID)
		}
	}
	return ids
}

// ResolveLastBusinessTask returns the ID of the highest-numbered business task.
var ResolveLastBusinessTask DepResolveFunc = func(ctx *GenContext) []string {
	if len(ctx.BusinessTasks) == 0 {
		return nil
	}
	var maxID string
	var maxNum int
	for _, t := range ctx.BusinessTasks {
		num := numericID(t.ID)
		if num > maxNum {
			maxNum = num
			maxID = t.ID
		}
	}
	if maxID == "" {
		return nil
	}
	return []string{maxID}
}

// ResolveHighestGateOrLastBiz returns the ID of the highest-phase gate/summary,
// or the last business task if its phase is higher. Gate priority over summary.
var ResolveHighestGateOrLastBiz DepResolveFunc = func(ctx *GenContext) []string {
	var dep string
	var depPhase int
	for _, t := range ctx.ExistingTasks {
		if strings.HasSuffix(t.ID, IDSuffixGate) {
			p := phaseFromID(t.ID)
			if p > depPhase {
				depPhase = p
				dep = t.ID
			}
		}
	}
	if dep == "" {
		for _, t := range ctx.ExistingTasks {
			if strings.HasSuffix(t.ID, IDSuffixSummary) {
				p := phaseFromID(t.ID)
				if p > depPhase {
					depPhase = p
					dep = t.ID
				}
			}
		}
	}
	var maxBizID string
	var maxBizNum int
	for _, t := range ctx.BusinessTasks {
		num := numericID(t.ID)
		if num > maxBizNum {
			maxBizNum = num
			maxBizID = t.ID
		}
	}
	if maxBizID != "" && maxBizNum > depPhase {
		dep = maxBizID
	}
	if dep == "" {
		return nil
	}
	return []string{dep}
}

// ResolveLastRunTestOrBusiness returns the last run-test task ID when test pipeline
// is active, otherwise falls back to the last business task ID.
var ResolveLastRunTestOrBusiness DepResolveFunc = func(ctx *GenContext) []string {
	if len(ctx.RunTestChain) > 0 {
		return []string{ctx.RunTestChain[len(ctx.RunTestChain)-1]}
	}
	if len(ctx.BusinessTasks) > 0 {
		return []string{ctx.BusinessTasks[len(ctx.BusinessTasks)-1].ID}
	}
	return nil
}

// ResolveIfGenerated returns the task ID if it was already generated, nil otherwise.
func ResolveIfGenerated(id string) DepResolveFunc {
	return func(ctx *GenContext) []string {
		for _, generated := range ctx.AllGenerated {
			if generated == id {
				return []string{id}
			}
		}
		return nil
	}
}

// Registry-derived ID lookup

// matchRegistryID attempts to match a task ID against registry node ID patterns.
// Returns the node's Type if matched, or "" if no match.
func matchRegistryID(id string, surfaces map[string]string) string {
	for _, node := range PipelineRegistry {
		if node.Expansion == "" {
			if id == node.ID {
				return node.Type
			}
			continue
		}
		switch node.Expansion {
		case "per-surface-type":
			if matchTypeSuffixedID(id, node.ID) {
				return node.Type
			}
		case "per-surface-key":
			if matched := matchSurfaceKeyID(id, node.ID, surfaces); matched {
				return node.Type
			}
		}
	}
	return ""
}

// matchTypeSuffixedID checks if id matches a template with {surface-type} placeholder.
// Also accepts degenerate form (no suffix) for backward compatibility.
func matchTypeSuffixedID(id, idTemplate string) bool {
	placeholder := "{surface-type}"
	idx := strings.Index(idTemplate, placeholder)
	if idx < 0 {
		return false
	}
	prefix := idTemplate[:idx]
	stripPrefix := strings.TrimSuffix(prefix, "-")
	if id == stripPrefix {
		return true
	}
	if !strings.HasPrefix(id, prefix) {
		return false
	}
	rem := id[len(prefix):]
	if len(rem) == 0 {
		return false
	}
	for _, c := range rem {
		if (c < 'a' || c > 'z') && c != '-' {
			return false
		}
	}
	return true
}

// matchSurfaceKeyID checks if id matches a template with {surface-key} placeholder.
// Also handles single-surface degenerate case.
func matchSurfaceKeyID(id, idTemplate string, surfaces map[string]string) bool {
	placeholder := "{surface-key}"
	idx := strings.Index(idTemplate, placeholder)
	if idx < 0 {
		return false
	}
	prefix := idTemplate[:idx]
	stripTemplate := strings.ReplaceAll(idTemplate, "-{surface-key}", "")
	if id == stripTemplate {
		return true
	}
	if !strings.HasPrefix(id, prefix) {
		return false
	}
	suffix := id[len(prefix):]
	if suffix == "" {
		return false
	}
	_, ok := surfaces[suffix]
	return ok
}

// PipelineRegistry — single source of truth for auto-generated tasks.
// Order determines generation sequence; execution order is determined by DependsOn.
var PipelineRegistry = []PipelineNode{
	{Type: TypeDocReview, Key: "review-doc", ID: "T-review-doc",
		Title: "Review Documentation Quality", Priority: string(types.PriorityP1), EstimatedTime: "30min",
		ConfigGate: nil, IntentGate: GateAllowAll, GenerateCondition: CondHasDocTasks,
		DependsOn: []DepRef{{Resolve: ResolveDocTasks}}},
	{Type: TypeCleanCode, Key: "clean-code", ID: "T-clean-code",
		Title: "Simplify and Clean Code", Priority: string(types.PriorityP2), EstimatedTime: "20min",
		ConfigGate: GateCleanCode, IntentGate: GateAllowAll, GenerateCondition: CondAlways,
		DependsOn: []DepRef{{Resolve: ResolveHighestGateOrLastBiz}}},
	{Type: TypeTestGenJourneys, Key: "gen-journeys", ID: "T-test-gen-journeys",
		Title: "Generate Test Journeys", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, StrategyKind: "interface",
		DependsOn: []DepRef{
			{Resolve: ResolveIfGenerated("T-review-doc")},
			{Resolve: ResolveIfGenerated("T-clean-code")},
			{Resolve: ResolveLastBusinessTask},
		}},
	{Type: TypeEvalJourney, Key: "eval-journey", ID: "T-eval-journey",
		Title: "Evaluate Journey Quality", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", MainSession: true,
		DependsOn: []DepRef{{Ref: "T-test-gen-journeys"}}},
	{Type: TypeTestGenContracts, Key: "gen-contracts", ID: "T-test-gen-contracts",
		Title: "Generate Test Contracts", Priority: string(types.PriorityP1), EstimatedTime: "30-45min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown",
		DependsOn: []DepRef{{Ref: "T-eval-journey"}}},
	{Type: TypeEvalContract, Key: "eval-contract", ID: "T-eval-contract",
		Title: "Evaluate Contract Quality", Priority: string(types.PriorityP1), EstimatedTime: "20-30min",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", MainSession: true,
		DependsOn: []DepRef{{Ref: "T-test-gen-contracts"}}},
	{Type: TypeTestGenScripts, Key: "gen-test-scripts-{surface-key}", ID: "T-test-gen-scripts-{surface-key}",
		Title: "Generate {test-type-title} Scripts", Priority: string(types.PriorityP1), EstimatedTime: estTimeMedium,
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks, Mode: "breakdown", Expansion: "per-surface-key",
		DependsOn: []DepRef{{Ref: "T-eval-contract"}}, StrategyKind: "generate"},
	{Type: TypeTestRun, Key: "run-test-{surface-key}", ID: "T-test-run-{surface-key}",
		Title: "Run {test-type-title}", Priority: string(types.PriorityP1), EstimatedTime: "30min-1h",
		ConfigGate: GateTest, GenerateCondition: CondHasTestableTasks,
		DependsOn: []DepRef{{Resolve: ResolveUpstream}}, Expansion: "per-surface-key", StrategyKind: "run"},
	{Type: TypeValidationCode, Key: "validate-code", ID: "T-validate-code",
		Title: "Validate Code Quality", Priority: string(types.PriorityP2), EstimatedTime: estTimeQuickTask,
		ConfigGate: GateValidation, IntentGate: GateAllowAll, GenerateCondition: CondAlways,
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}}},
	{Type: TypeValidationUx, Key: "validate-ux", ID: "T-validate-ux",
		Title: "Validate User Experience", Priority: string(types.PriorityP2), EstimatedTime: estTimeQuickTask,
		ConfigGate: GateValidation, IntentGate: GateAllowAll, GenerateCondition: CondAlways,
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}}, UISurfaceOnly: true, MainSession: true},
	{Type: TypeDocConsolidate, Key: "consolidate-specs", ID: "T-specs-consolidate",
		Title: "Consolidate Specs", Priority: string(types.PriorityP2), EstimatedTime: "20min",
		ConfigGate: GateConsolidateSpecs, IntentGate: GateAllowAll, GenerateCondition: CondAlways, Mode: "breakdown",
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}}},
	{Type: TypeDocDrift, Key: "quick-drift-detection", ID: "T-quick-doc-drift",
		Title: "Detect Spec Drift", Priority: string(types.PriorityP2), EstimatedTime: estTimeQuickTask,
		ConfigGate: GateConsolidateSpecs, IntentGate: GateAllowAll, GenerateCondition: CondAlways, Mode: "quick",
		DependsOn: []DepRef{{Resolve: ResolveLastRunTestOrBusiness}}},
}

// Two-phase validation

// escapeHatchLimit is the maximum allowed number of post-generation injection functions.
const escapeHatchLimit = 5

func init() {
	if err := ValidatePipelineRegistry(); err != nil {
		panic("pipeline registry validation failed: " + err.Error())
	}
}

// ValidatePipelineRegistry performs Phase 1 (static) validation of the pipeline registry.
// Runs at init-time; panics on failure with actionable error messages.
func ValidatePipelineRegistry() error {
	nodeIDs := make(map[string]int)
	for i, node := range PipelineRegistry {
		if node.GenerateCondition == nil {
			return fmt.Errorf("node %q (index %d): GenerateCondition must be non-nil, use CondAlways for unconditional generation", node.ID, i)
		}
		if err := validatePlaceholders(node, i); err != nil {
			return err
		}
		for _, dep := range node.DependsOn {
			if dep.Resolve == nil && dep.Ref != "" {
				if _, ok := nodeIDs[dep.Ref]; !ok {
					if !idExistsInRegistry(dep.Ref) {
						return fmt.Errorf("node %q (index %d): DependsOn.Ref %q does not match any registry node ID", node.ID, i, dep.Ref)
					}
				}
			}
			if dep.Resolve != nil && isResolveIfGenerated(dep.Resolve) {
				refID := extractResolveIfGeneratedID(dep.Resolve)
				if refID != "" {
					if _, ok := nodeIDs[refID]; !ok {
						return fmt.Errorf("node %q (index %d): ResolveIfGenerated(%q) references a node not yet declared", node.ID, i, refID)
					}
				}
			}
		}
		nodeIDs[node.ID] = i
	}
	if err := validateExpandedIDsUnique(); err != nil {
		return err
	}
	if escapeHatchCount := 0; escapeHatchCount > escapeHatchLimit {
		return fmt.Errorf("escape hatch count %d exceeds limit %d", escapeHatchCount, escapeHatchLimit)
	}
	if err := validateOrderingInvariants(); err != nil {
		return err
	}
	return nil
}

func validatePlaceholders(node PipelineNode, idx int) error {
	hasSurfaceKey := strings.Contains(node.ID, "{surface-key}") || strings.Contains(node.Key, "{surface-key}")
	hasSurfaceType := strings.Contains(node.ID, "{surface-type}") || strings.Contains(node.Key, "{surface-type}")
	switch node.Expansion {
	case "per-surface-key":
		if !hasSurfaceKey {
			return fmt.Errorf("node %q (index %d): Expansion is per-surface-key but ID/Key lacks {surface-key}", node.ID, idx)
		}
	case "per-surface-type":
		if !hasSurfaceType {
			return fmt.Errorf("node %q (index %d): Expansion is per-surface-type but ID/Key lacks {surface-type}", node.ID, idx)
		}
	default:
		if hasSurfaceKey || hasSurfaceType {
			return fmt.Errorf("node %q (index %d): Expansion is empty but ID/Key contains placeholders", node.ID, idx)
		}
	}
	return nil
}

func idExistsInRegistry(id string) bool {
	for _, node := range PipelineRegistry {
		if node.ID == id {
			return true
		}
		if strings.Contains(node.ID, "{") && matchIDToTemplate(id, node.ID) {
			return true
		}
	}
	return false
}

func matchIDToTemplate(id, template string) bool {
	parts := strings.SplitN(template, "{", 2)
	if len(parts) < 2 {
		return id == template
	}
	return strings.HasPrefix(id, parts[0]) && len(id) > len(parts[0])
}

func isResolveIfGenerated(fn DepResolveFunc) bool {
	return extractResolveIfGeneratedID(fn) != ""
}

func extractResolveIfGeneratedID(fn DepResolveFunc) string {
	allIDs := make([]string, 0, len(PipelineRegistry))
	for _, node := range PipelineRegistry {
		allIDs = append(allIDs, node.ID)
	}
	ctx := &GenContext{AllGenerated: allIDs}
	result := fn(ctx)
	if len(result) == 1 {
		ctx2 := &GenContext{AllGenerated: nil}
		if fn(ctx2) == nil {
			return result[0]
		}
	}
	return ""
}

func validateExpandedIDsUnique() error {
	testSurfaces := map[string]string{
		".": "api", "api": "api", "cli": "cli", "tui": "tui", "web": "web", "mobile": "mobile",
	}
	seen := make(map[string]string)
	for _, node := range PipelineRegistry {
		for _, t := range expandNode(node, testSurfaces, nil) {
			if prev, dup := seen[t.ID]; dup {
				return fmt.Errorf("duplicate expanded ID %q produced by both %q and %q", t.ID, prev, node.Type)
			}
			seen[t.ID] = node.Type
		}
	}
	singleSurfaces := map[string]string{".": "api"}
	singleSeen := make(map[string]string)
	for _, node := range PipelineRegistry {
		for _, t := range expandNode(node, singleSurfaces, nil) {
			if prev, dup := singleSeen[t.ID]; dup {
				return fmt.Errorf("duplicate expanded ID %q (single-surface) produced by both %q and %q", t.ID, prev, node.Type)
			}
			singleSeen[t.ID] = node.Type
		}
	}
	return nil
}

func validateOrderingInvariants() error {
	var hasUpstreamProducer bool
	var hasRunTestProducer bool
	for i, node := range PipelineRegistry {
		for _, dep := range node.DependsOn {
			if dep.Resolve == nil {
				continue
			}
			if isSameResolver(dep.Resolve, ResolveUpstream) && !hasUpstreamProducer {
				return fmt.Errorf("node %q (index %d): uses ResolveUpstream but no upstream producer before it", node.ID, i)
			}
			if isSameResolver(dep.Resolve, ResolveLastRunTest) && !hasRunTestProducer {
				return fmt.Errorf("node %q (index %d): uses ResolveLastRunTest but no TypeTestRun before it", node.ID, i)
			}
		}
		if node.Expansion == "" || node.Expansion == "per-surface-key" {
			hasUpstreamProducer = true
		}
		if node.Type == TypeTestRun {
			hasRunTestProducer = true
		}
	}
	return nil
}

func isSameResolver(a, b DepResolveFunc) bool {
	return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)
}

// validateGeneratedTasks performs Phase 2 (dynamic) validation.
func validateGeneratedTasks(tasks []AutoGenTaskDef) error {
	idSet := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		idSet[t.ID] = true
	}
	for _, t := range tasks {
		for _, dep := range t.Dependencies {
			if !idSet[dep] && strings.HasPrefix(dep, "T-") {
				return fmt.Errorf("generated task %q depends on %q which is not in the generated task set", t.ID, dep)
			}
		}
	}
	return checkNoCycles(tasks)
}

func checkNoCycles(tasks []AutoGenTaskDef) error {
	idSet := make(map[string]bool, len(tasks))
	for _, t := range tasks {
		idSet[t.ID] = true
	}
	forward := make(map[string][]string)
	inDegree := make(map[string]int)
	for _, t := range tasks {
		inDegree[t.ID] = 0
	}
	for _, t := range tasks {
		for _, dep := range t.Dependencies {
			if idSet[dep] {
				forward[dep] = append(forward[dep], t.ID)
				inDegree[t.ID]++
			}
		}
	}
	queue := make([]string, 0)
	for _, t := range tasks {
		if inDegree[t.ID] == 0 {
			queue = append(queue, t.ID)
		}
	}
	visited := 0
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		visited++
		for _, next := range forward[curr] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	if visited < len(tasks) {
		return fmt.Errorf("circular dependency detected among %d pipeline tasks (topological sort visited %d of %d)", len(tasks)-visited, visited, len(tasks))
	}
	return nil
}
