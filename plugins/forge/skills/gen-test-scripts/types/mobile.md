---
type: mobile
conventions:
  - testing-mobile.md
---

# Mobile Type Instructions

Type-specific generation instructions for **Mobile** (touch, gestures, screen transitions) test scripts. Loaded by the dispatcher when interface detection identifies Mobile-type test cases.

**Test type**: 移动端端到端测试 (Mobile E2E Test). See `docs/reference/test-type-model.md` for the authoritative definition. Generated test code MUST use `@mobile-e2e` tags. This is one of the two surfaces where "e2e" terminology is correct (the other being Web).

This file defines two zones:

- **Golden Rules**: Framework-agnostic constraints that govern all Mobile test generation. These rules are mandatory and cannot be overridden by Convention files.
- **Reconnaissance Hints**: Discovery-only guidance for extracting Mobile-specific information from source code. Hints inform the Fact Table, not the generated code directly.

## Golden Rules

Mobile-specific constraints beyond the universal principles in `_shared.md` (Isolation, Determinism, Timeout Protection, Idempotency, Resource Cleanup).

### App State Reset

Every test must clean application state before execution. No test may depend on or leave behind state from a previous test.

**Constraint**: Before each test executes, the application must be returned to a clean baseline state. Acceptable methods include: kill and clear application data, or uninstall and reinstall. The specific mechanism is defined by the active Convention, not by this rule.

**Rationale**: Mobile applications accumulate persistent state (local storage, cached data, user preferences, authentication tokens) that causes order-dependent failures. A test that passes in isolation fails when run after another test that modified the app state.

**Antipattern guard**: Tests that assume the app is in a specific state without enforcing it, or tests that rely on data created by a previous test.

### Permission Handling

System permission dialogs (camera, location, notifications, contacts) must be handled as a pre-step, not as an ad-hoc dismissal during test execution.

**Constraint**: Permissions must be pre-authorized before test execution, or system permission dialogs must be explicitly handled as a setup step at the beginning of the test. No test may assume permissions are already granted, and no test may dismiss a permission dialog mid-flow without declaring it as a pre-step.

**Rationale**: System permission dialogs are non-deterministic -- they appear on first request and not thereafter, or reappear after app reinstalls. Unhandled permission dialogs block test execution and cause timeout failures.

**Antipattern guard**: Tests that trigger a permission-requiring action without first handling the permission dialog, or tests that assume a specific permission state.

### Deep Link Pattern

Opening the application via URL scheme is a supported navigation entry point. Tests may use deep links as an alternative to in-app navigation.

**Constraint**: Deep links (URL schemes that open the application to a specific screen) are a valid navigation method for test setup. Tests using deep links must declare the URL scheme and target path explicitly. The deep link mechanism is provided by the Convention, not by this rule.

**Rationale**: Deep links provide a fast, deterministic way to navigate to specific application states without traversing the full UI flow. This is especially valuable for tests that target screens deep in the navigation hierarchy.

**Antipattern guard**: Tests that navigate through 5+ intermediate screens to reach a target when a deep link could skip directly to it.

### Element Location Strategy

Element selection follows an abstract priority chain, not a framework-specific selector syntax:

1. **Accessibility ID**: Use accessibility labels, test IDs, or accessibility identifiers from Fact Table entries
2. **Resource ID**: Use platform resource identifiers (Android resource ID, iOS accessibility identifier)
3. **Text content**: Use visible text as a last resort for element identification

**Constraint**: Tests must prefer stable, unique identifiers (accessibility labels, resource IDs) over visible text. Text-based element location is fragile -- it breaks when copy changes and fails when the same text appears in multiple elements.

**Hard Rule**: Never guess accessibility label or resource ID values. Every identifier locator must come from a Fact Table entry. If the Fact Table has no entry for the needed element, use text-based location as a fallback.

### Touch and Gesture Principles

Touch interactions follow framework-agnostic action categories:

| Action Category | Description | Examples |
|----------------|-------------|----------|
| Tap | Brief contact with a single element | Single tap, double-tap, tap-and-hold |
| Swipe | Directional gesture across the screen | Swipe left/right/up/down, swipe to scroll |
| Input | Text entry into focused fields | Type text, clear field, paste |
| Navigate | Move between screens or application states | Navigate forward, navigate back, open deep link |
| System | Platform-level interactions | Launch application, terminate application, send to background |

**Constraint**: Generated test code expresses interactions using these abstract categories. The Convention translates these categories into framework-specific commands. Golden Rules do not contain any framework-specific command syntax.

**Antipattern guard**: Tests that use pixel coordinates for touch interactions instead of element-based targeting. Coordinates are device-dependent and break across screen sizes.

### Screen Transition Assertions

After navigation steps, assert the expected screen is visible before proceeding. Navigation is asynchronous on mobile -- the next screen is not immediately available.

**Constraint**: Every navigation action must be followed by an assertion confirming the target screen is visible. The assertion waits for the target element to appear within a timeout window. The timeout value comes from the Convention, not from the test code.

**Rationale**: Mobile navigation involves animations, network loading, and state transitions that are inherently asynchronous. Proceeding without confirming the target screen causes subsequent interactions to target elements from the previous screen, producing false failures.

**Antipattern guard**: Navigation actions without a follow-up screen assertion, or fixed-duration waits between navigation and interaction.

### Application Lifecycle

Tests must declare application lifecycle events explicitly. The application must be in a known state before and after each test.

**Constraint**: Each test must define its application lifecycle: launch at start, terminate at end. The specific lifecycle mechanism is defined by the Convention. No test may assume the application is already running.

**Rationale**: Application lifecycle state affects test behavior. A test that assumes the app is running will fail if the app was terminated by a previous test's cleanup. Explicit lifecycle declaration makes each test independently executable.

**Antipattern guard**: Tests that do not launch the application at the start, or tests that do not clean up the application state at the end.

## Reconnaissance Hints

<!-- Discovery-only: information extracted here populates the Fact Table. These hints do not directly guide code generation. -->

Mobile reconnaissance discovers the project's app structure, screen definitions, and navigation routes from source code and configuration files.

### Search Targets

| Target | What It Finds | Discovery Method |
|--------|---------------|------------------|
| App manifest (Android) | App package identifier | Search for `package=` in `AndroidManifest.xml` |
| App manifest (iOS) | iOS bundle identifier | Search for `CFBundleIdentifier` in `.plist` files |
| React Native screens | Screen definitions and navigation calls | Search for `Stack.Screen` or `navigation.navigate` in `.tsx`/`.jsx` files |
| Flutter routes | Route definitions | Search for `Navigator.push`, `GetPage`, or `routes:` in `.dart` files |
| Accessibility labels | Element identifiers for test targeting | Search for `accessibilityLabel`, `testID`, `accessibilityIdentifier` in `.tsx`/`.jsx`/`.swift`/`.kt` files |
| Entry point | App entry point | Search for `AppRegistry` or `main()` in `.tsx`/`.jsx`/`.dart`/`.swift` files |
| App configuration | App identifiers in config | Search for `appId`, `bundleId`, `applicationId` in `.json`/`.gradle`/`.yaml` files |
| Existing test flows | Test framework and conventions | Search for test configuration files (`.yaml`, `.json`, feature files) |

### Reconnaissance Procedure

1. **Detect mobile framework**: Identify whether the project uses React Native, Flutter, native Android, or native iOS from file structure and dependency files.
2. **Extract app identifier**: Find the app bundle identifier (Android package name or iOS bundle ID). This is required for launching the application in tests.
3. **Map screen definitions**: Extract screen names and navigation routes from the source code.
4. **Discover accessibility labels**: Find `testID`, `accessibilityLabel`, or resource ID values used in the app's components. These become element identifiers in test scripts.
5. **Locate existing test configurations**: If the project already has mobile test configuration files, analyze their structure for conventions and patterns.

### Reference Example for Maestro

<!-- Reference example for Maestro -- not generation instructions -->

The following Maestro YAML patterns illustrate how the abstract principles above map to a specific framework. This section is for reconnaissance reference only -- it shows what existing Maestro flows look like, not how to generate them.

**Flow structure** typically follows: app identifier declaration, environment variables, lifecycle hooks, and a command sequence.

**Touch mapping** -- Maestro's command syntax for the abstract action categories:

| Action Category | Maestro Command |
|----------------|-----------------|
| Tap | `tapOn` |
| Long press | `longPressOn` |
| Swipe | `swipe` |
| Scroll to element | `scroll` |
| Text input | `tapOn` + `inputText` |
| Navigate back | `back` |
| Launch app | `launchApp` |
| Terminate app | `killApp` |

**Element location in Maestro** uses a priority: text match > resource ID > index > relative position.

**Screen assertion in Maestro**: `assertVisible` and `assertNotVisible` commands confirm screen state.

**App lifecycle in Maestro**: `onFlowStart: [launchApp]` and `onFlowEnd: [killApp]` hooks manage lifecycle.

## Classification Indicators

Classify test cases as **Mobile** when they involve any of:

- Touch interactions (tap, double-tap, long-press)
- Gestures (swipe, pinch, drag, rotate)
- Screen transitions and navigation flows
- Accessibility labels and resource IDs
- App lifecycle events (background, foreground, terminate)
- Platform-specific UI components (bottom sheets, native dialogs, permissions)
- Push notifications, deep links

## Fact Table Required Keys

After reconnaissance, the Fact Table must contain at least these Mobile-specific entries for the completeness gate to pass:

| Key Pattern | Description | Example |
|-------------|-------------|---------|
| `MOBILE_APP_ID` | App bundle identifier (package name) | `MOBILE_APP_ID` = `com.example.myapp` |
| `MOBILE_FRAMEWORK` | Detected mobile framework | `MOBILE_FRAMEWORK` = `react-native` |
| `MOBILE_SCREEN_*` | At least one screen name entry | `MOBILE_SCREEN_LOGIN` = `Login` |

**Minimum requirement**: `MOBILE_APP_ID` must be non-UNKNOWN. Without the app identifier, the test framework cannot launch the app and no tests can execute. If `MOBILE_APP_ID` is UNKNOWN, skip Mobile test generation and emit a WARNING.

**Completeness gate rule**: If all required keys for Mobile are UNKNOWN, do NOT generate Mobile tests. Individual UNKNOWN keys for screens or accessibility labels are acceptable -- generate tests for discovered screens only.

## Verification Method

Before generating Mobile test scripts, confirm the project is actually a mobile application. A web app or CLI tool does not need Mobile test scripts.

Run these checks in order -- first success is sufficient:

| Check | What to Look For |
|-------|-------------------|
| Mobile framework dependency | `react-native` in `package.json`, or Flutter in `pubspec.yaml`, or Expo in `app.json` |
| Native project files | `AndroidManifest.xml` or `Info.plist` with `CFBundleIdentifier` |
| Existing mobile test configs | Configuration files referencing mobile testing frameworks |
| App manifest identifiers | Application ID in build configuration files |

**If all checks fail**: The project is not a mobile application. Skip Mobile test generation and emit a WARNING suggesting the user verify the project is a mobile app.

## Mobile Antipattern Guards

Beyond the shared antipattern guards in `_shared.md` (Sleep-Based Waits, Hardcoded Configuration, Vacuous Assertions, Source-Code-Level Testing), Mobile-specific forbidden patterns:

| Pattern | Why Forbidden | What To Do Instead |
|---------|--------------|-------------------|
| **Pixel-coordinate-based interactions** | Coordinates are device-dependent and break across screen sizes and resolutions | Use element-based targeting (accessibility ID, resource ID, text) |
| **State leakage between tests** | Tests that pass in isolation fail when run in sequence due to accumulated app state | Clean app state before each test (see Golden Rule: App State Reset) |
| **Unhandled permission dialogs** | System dialogs block test execution and cause timeout failures | Pre-authorize permissions or handle dialogs as an explicit pre-step |
| **Physical device capability assumptions** | Tests fail in CI environments using emulators without sensor simulation | Limit test scripts to interactions supported by emulator/simulator environments; mark tests requiring physical devices explicitly |
| **Ad-hoc navigation without screen assertion** | Subsequent interactions target elements from the wrong screen, producing false failures | Assert target screen visibility after every navigation step (see Golden Rule: Screen Transition Assertions) |

## Output

Mobile test scripts are written to `tests/<journey>/` following the active Convention's file naming and structure. Each test file includes a traceability comment linking back to the source Contract step.

## Test Ratio Constraint (Best-Effort)

Mobile surface follows a **best-effort** strategy — not measured by Contract test ratio.

- **Output format**: Maestro YAML files (`.yaml` extension)
- **Skeleton structure**: Each generated Maestro YAML MUST contain at minimum:
  1. `appId` declaration (from Fact Table `MOBILE_APP_ID`)
  2. `onFlowStart: [launchApp]` lifecycle hook
  3. `onFlowEnd: [killApp]` lifecycle hook
  4. Command sequence for the Contract step's happy path actions
  5. `assertVisible` assertions for each expected screen state
- **Deep link tests**: For each Journey step that navigates to a specific screen, generate an additional Maestro YAML file that:
  1. Opens the app via URL scheme (e.g., `myapp://screen/detail`)
  2. Asserts the target screen is visible
  3. Verifies expected content on the target screen
- **File naming**: `tests/<journey>/step<N>_<action>.yaml` for step tests, `tests/<journey>/step<N>_<action>_deeplink.yaml` for deep link variants

### Manual-Only Marking

Complex scenarios that cannot be reliably automated in Maestro YAML MUST be marked as `manual-only`:

| Scenario | Why Manual |
|----------|-----------|
| Multi-finger gestures (pinch, rotate) | Maestro does not support complex multi-touch |
| Physical device sensors (accelerometer, camera) | Emulator/simulator cannot replicate |
| Biometric authentication (Face ID, fingerprint) | Requires physical device or special simulator config |
| System-level interactions (notifications, clipboard) | Unreliable in automated environments |

**Marking convention**: Generate a placeholder Maestro YAML with `manual-only` status:

```yaml
# MANUAL-ONLY: <reason>
# This test scenario requires <capability> which cannot be automated via Maestro.
# Manual test procedure: <description of what to test manually>

appId: ${MOBILE_APP_ID}
---
- launchApp
- assertVisible: "Home Screen"
# Remaining steps require manual execution
```

<HARD-RULE>
Mobile test generation MUST NOT fail the pipeline. Any generation issue results in a skeleton with `manual-only` markers or `gen-failed` annotation. Mobile tests are best-effort — incomplete coverage is acceptable.
</HARD-RULE>
