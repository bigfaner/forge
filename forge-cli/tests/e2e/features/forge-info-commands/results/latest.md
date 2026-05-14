=== RUN   TestTC_001_ConfigInitInteractiveSetup
    forge_info_commands_cli_test.go:22: requires interactive stdin: multi-step prompt with project-type, profiles, capabilities selections
--- SKIP: TestTC_001_ConfigInitInteractiveSetup (0.00s)
=== RUN   TestTC_002_ConfigInitReconfigurePrompt
    forge_info_commands_cli_test.go:27: requires interactive stdin and pre-existing .forge/config.yaml to trigger reconfigure prompt
--- SKIP: TestTC_002_ConfigInitReconfigurePrompt (0.00s)
=== RUN   TestTC_003_ConfigInitReconfigureAccepted
    forge_info_commands_cli_test.go:32: requires interactive stdin and pre-existing .forge/config.yaml with 'y' confirmation input
--- SKIP: TestTC_003_ConfigInitReconfigureAccepted (0.00s)
=== RUN   TestTC_004_ConfigGetProjectType
--- PASS: TestTC_004_ConfigGetProjectType (0.02s)
=== RUN   TestTC_005_ConfigGetCapabilitiesArrayOutput
--- PASS: TestTC_005_ConfigGetCapabilitiesArrayOutput (0.01s)
=== RUN   TestTC_006_ConfigGetMissingKey
--- PASS: TestTC_006_ConfigGetMissingKey (0.01s)
=== RUN   TestTC_007_ForgeConfigStructFields
    forge_info_commands_cli_test.go:71: struct field validation requires Go source inspection, not CLI invocation; covered by unit tests
--- SKIP: TestTC_007_ForgeConfigStructFields (0.00s)
=== RUN   TestTC_008_ProposalListAllProposals
--- PASS: TestTC_008_ProposalListAllProposals (0.03s)
=== RUN   TestTC_009_ProposalSlugDetailView
--- PASS: TestTC_009_ProposalSlugDetailView (0.02s)
=== RUN   TestTC_010_ProposalCreatedDateFromFrontmatter
--- PASS: TestTC_010_ProposalCreatedDateFromFrontmatter (0.02s)
=== RUN   TestTC_011_ProposalPRDColumnChecksPrdSpec
--- PASS: TestTC_011_ProposalPRDColumnChecksPrdSpec (0.02s)
=== RUN   TestTC_012_ProposalFeatureColumnReadsManifestStatus
--- PASS: TestTC_012_ProposalFeatureColumnReadsManifestStatus (0.01s)
=== RUN   TestTC_013_FeatureListAllFeatures
--- PASS: TestTC_013_FeatureListAllFeatures (0.02s)
=== RUN   TestTC_014_FeatureListProgressFromIndexJSON
    forge_info_commands_cli_test.go:195: requires manual setup: feature with known task counts in index.json for precise progress assertion
--- SKIP: TestTC_014_FeatureListProgressFromIndexJSON (0.00s)
=== RUN   TestTC_015_FeatureListScoresFromFrontmatter
--- PASS: TestTC_015_FeatureListScoresFromFrontmatter (0.01s)
=== RUN   TestTC_016_FeatureStatusDetailView
--- PASS: TestTC_016_FeatureStatusDetailView (0.01s)
=== RUN   TestTC_017_FeatureNoArgsKeepsExistingBehavior
--- PASS: TestTC_017_FeatureNoArgsKeepsExistingBehavior (0.01s)
=== RUN   TestTC_018_LessonListAllLessons
--- PASS: TestTC_018_LessonListAllLessons (0.03s)
=== RUN   TestTC_019_LessonCategoryFromFilenamePrefix
--- PASS: TestTC_019_LessonCategoryFromFilenamePrefix (0.02s)
=== RUN   TestTC_020_LessonNameDetailView
--- PASS: TestTC_020_LessonNameDetailView (0.03s)
=== RUN   TestTC_021_InitCreatesForgeDir
    forge_info_commands_cli_test.go:325: requires clean project state (no .forge/ directory); destructive to run against real project
--- SKIP: TestTC_021_InitCreatesForgeDir (0.00s)
=== RUN   TestTC_022_InitGeneratesCLAUDEmd
    forge_info_commands_cli_test.go:330: requires clean project state (no CLAUDE.md); destructive to run against real project
--- SKIP: TestTC_022_InitGeneratesCLAUDEmd (0.00s)
=== RUN   TestTC_023_InitAppendsGitignoreWithDedup
    forge_info_commands_cli_test.go:335: requires isolated project state; modifies .gitignore which is destructive
--- SKIP: TestTC_023_InitAppendsGitignoreWithDedup (0.00s)
=== RUN   TestTC_024_InitGitignoreDedupSkipsExisting
    forge_info_commands_cli_test.go:340: requires .gitignore with pre-existing forge entries; modifies .gitignore
--- SKIP: TestTC_024_InitGitignoreDedupSkipsExisting (0.00s)
=== RUN   TestTC_025_InitAppendsJustfileRecipes
    forge_info_commands_cli_test.go:345: requires isolated project state; modifies justfile which is destructive
--- SKIP: TestTC_025_InitAppendsJustfileRecipes (0.00s)
=== RUN   TestTC_026_InitJustfileDedupSkipsExisting
    forge_info_commands_cli_test.go:350: requires justfile with pre-existing claude recipe; modifies justfile
--- SKIP: TestTC_026_InitJustfileDedupSkipsExisting (0.00s)
=== RUN   TestTC_027_InitRunsConfigInitWhenNoConfig
    forge_info_commands_cli_test.go:355: requires clean project state (no .forge/config.yaml) and interactive stdin
--- SKIP: TestTC_027_InitRunsConfigInitWhenNoConfig (0.00s)
=== RUN   TestTC_028_InitSkipsExistingFiles
    forge_info_commands_cli_test.go:360: requires pre-existing .forge/, CLAUDE.md, .forge/config.yaml; complex setup
--- SKIP: TestTC_028_InitSkipsExistingFiles (0.00s)
=== RUN   TestTC_029_InitResultReportFormat
    forge_info_commands_cli_test.go:365: requires clean project state; destructive to run against real project
--- SKIP: TestTC_029_InitResultReportFormat (0.00s)
=== RUN   TestTC_030_ResolveScopeReadsConfigDirectly
    forge_info_commands_cli_test.go:374: ResolveScope() is an internal function, not a CLI command; covered by unit tests
--- SKIP: TestTC_030_ResolveScopeReadsConfigDirectly (0.00s)
=== RUN   TestTC_031_ResolveScopeMissingConfigReturnsEmpty
    forge_info_commands_cli_test.go:379: ResolveScope() is an internal function, not a CLI command; covered by unit tests
--- SKIP: TestTC_031_ResolveScopeMissingConfigReturnsEmpty (0.00s)
=== RUN   TestTC_032_JustfileHasNoProjectTypeRecipe
    forge_info_commands_cli_test.go:391: cannot locate justfile for migration check
--- SKIP: TestTC_032_JustfileHasNoProjectTypeRecipe (0.00s)
PASS
ok  	forge-cli/tests/e2e/features/forge-info-commands	0.819s
