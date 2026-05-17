---
type: mobile
conventions:
  - testing-mobile.md
---

# Mobile Test Script Generation Instructions

Type-specific Steps for **Mobile** (touch, gestures, screen transitions) test script generation. Loaded by the dispatcher when interface detection identifies Mobile-type test cases.

## Classification Indicators

Classify test cases as **Mobile** when they involve any of:

- Touch interactions (tap, double-tap, long-press)
- Gestures (swipe, pinch, drag, rotate)
- Screen transitions and navigation flows
- Accessibility labels and resource IDs
- App lifecycle events (background, foreground, terminate)
- Platform-specific UI components (bottom sheets, native dialogs, permissions)
- Push notifications, deep links

## Reconnaissance Strategy

Mobile reconnaissance discovers the project's app structure, screen definitions, and navigation routes from source code and configuration files.

### Search Commands

Run these searches to discover Mobile interface details. Adapt file extensions to the project's language and framework.

| Target | Grep Command | What It Finds |
|--------|-------------|---------------|
| App manifest (Android) | `grep -rn "package=" --include='AndroidManifest.xml' .` | App package identifier |
| App manifest (iOS) | `grep -rn "CFBundleIdentifier" --include='*.plist' .` | iOS bundle identifier |
| React Native screens | `grep -rn "Stack.Screen\\|navigation.navigate" --include='*.tsx' --include='*.jsx' .` | Screen definitions and navigation calls |
| Flutter routes | `grep -rn "Navigator.push\\|GetPage\\|routes:" --include='*.dart' .` | Route definitions |
| Maestro flow configs | `find . -name '*.yaml' -exec grep -l 'appId\\|launchApp' {} \\;` | Existing Maestro test files |
| Accessibility labels | `grep -rn "accessibilityLabel\\|testID\\|accessibilityIdentifier" --include='*.tsx' --include='*.jsx' --include='*.swift' --include='*.kt' .` | Element identifiers for test targeting |
| Entry point | `grep -rn "AppRegistry\\|main()" --include='*.tsx' --include='*.jsx' --include='*.dart' --include='*.swift' .` | App entry point |
| App configuration | `grep -rn '"appId"\\|"bundleId"\\|"applicationId"' --include='*.json' --include='*.gradle' --include='*.yaml' .` | App identifiers in config files |

### Reconnaissance Procedure

1. **Detect mobile framework**: Run the grep commands above. Identify whether the project uses React Native, Flutter, native Android, or native iOS.
2. **Extract app identifier**: Find the app bundle identifier (Android package name or iOS bundle ID). This is required for Maestro's `appId` field.
3. **Map screen definitions**: Extract screen names and navigation routes from the source code.
4. **Discover accessibility labels**: Find `testID`, `accessibilityLabel`, or resource ID values used in the app's components. These become element selectors in test scripts.
5. **Locate existing Maestro flows**: If the project already has Maestro YAML files, analyze their structure for conventions and patterns.

## Fact Table Required Keys

After reconnaissance, the Fact Table must contain at least these Mobile-specific entries for the completeness gate to pass:

| Key Pattern | Description | Example |
|-------------|-------------|---------|
| `MOBILE_APP_ID` | App bundle identifier (package name) | `MOBILE_APP_ID` = `com.example.myapp` |
| `MOBILE_FRAMEWORK` | Detected mobile framework | `MOBILE_FRAMEWORK` = `react-native` |
| `MOBILE_SCREEN_*` | At least one screen name entry | `MOBILE_SCREEN_LOGIN` = `Login` |

**Minimum requirement**: `MOBILE_APP_ID` must be non-UNKNOWN. Without the app identifier, Maestro cannot launch the app and no tests can execute. If `MOBILE_APP_ID` is UNKNOWN, skip Mobile test generation and emit a WARNING.

**Completeness gate rule** (from SKILL.md Step 1.5): If all required keys for Mobile are UNKNOWN, do NOT generate Mobile tests. Individual UNKNOWN keys for screens or accessibility labels are acceptable -- generate tests for discovered screens only.

## Verification Method

Before generating Mobile test scripts, confirm the project is actually a mobile application. A web app or CLI tool does not need Mobile test scripts.

Run these checks in order -- first success is sufficient:

| Check | Command | Pass Condition |
|-------|---------|----------------|
| Maestro config | `find . -name '*.yaml' -exec grep -l 'appId' {} \\;` | At least one Maestro flow file found |
| React Native | `grep -rn "react-native" package.json` | `react-native` in dependencies |
| Flutter | `ls pubspec.yaml && grep "flutter" pubspec.yaml` | Flutter project detected |
| Android manifest | `find . -name 'AndroidManifest.xml'` | Android manifest exists |
| iOS bundle | `find . -name 'Info.plist' -exec grep -l "CFBundleIdentifier" {} \\;` | iOS bundle identifier found |
| Expo | `grep -rn "expo" app.json` | Expo project detected |

**If all checks fail**: The project is not a mobile application. Skip Mobile test generation and emit a WARNING suggesting the user verify the project is a mobile app.

## Generation Patterns

Mobile test cases translate to executable Maestro YAML flows. Follow the active strategy's `generate.md` for the Maestro-specific template structure (flow skeleton, env variables, lifecycle hooks).

### Flow Skeleton

Each Maestro flow follows this structure (from the Maestro strategy template):

```yaml
appId: ${MOBILE_APP_ID}
name: TC-NNN Description
# Traceability: TC-NNN -> {PRD Source}
env:
  KEY: ${VALUE}
onFlowStart:
  - launchApp
onFlowEnd:
  - killApp
commands:
  - ...
```

### Touch and Gesture Simulation

Map test case step descriptions to Maestro commands:

| Test Case Step Pattern | Maestro Command | Example |
|------------------------|-----------------|---------|
| "Tap the X button" | `tapOn` | `tapOn: { text: "Login" }` |
| "Double-tap the item" | `tapOn` (double) | `tapOn: { text: "Item", repeat: 2 }` |
| "Long-press the card" | `longPressOn` | `longPressOn: { text: "Card" }` |
| "Swipe left" | `swipe` | `swipe: { direction: LEFT }` |
| "Scroll to element" | `scroll` | `scroll: { to: { text: "Footer" } }` |
| "Enter text in field" | `tapOn` + `inputText` | `tapOn: { text: "Email" }` then `inputText: user@example.com` |
| "Press back button" | `back` | `back` |
| "Launch app" | `launchApp` | `launchApp` |
| "Kill app" | `killApp` | `killApp` |

### Element Location via Accessibility Labels

Maestro selects elements using a priority order (from the Maestro strategy's `generate.md`):

| Priority | Method | Example |
|----------|--------|---------|
| 1 | Text match | `tapOn: { text: "Login" }` |
| 2 | Resource ID | `tapOn: { id: "com.app:id/button" }` |
| 3 | Index | `tapOn: { index: 0 }` |
| 4 | Relative | `tapOn: { below: { text: "Label" } }` |

When the test case references an accessibility label or resource ID discovered during reconnaissance, use it directly. When only a textual description is available, use text match.

### Screen Transition Assertions

After navigation steps, assert the expected screen is visible:

```yaml
commands:
  - tapOn: { text: "Login" }
  # ... login steps ...
  - assertVisible: { text: "Dashboard" }
```

Map test case Expected fields to Maestro assertions:

| Expected Pattern | Maestro Command |
|------------------|-----------------|
| "X is visible" | `assertVisible: { text: "X" }` |
| "X is not visible" | `assertNotVisible: { text: "X" }` |

### App Lifecycle Events

Map lifecycle-related test case steps to Maestro hooks and commands:

| Lifecycle Step | Maestro Mapping |
|----------------|-----------------|
| App launch (start of test) | `onFlowStart: [launchApp]` |
| App kill (end of test) | `onFlowEnd: [killApp]` |
| Send app to background | Not directly supported -- use `killApp` + `launchApp` as approximation |
| Screenshot for evidence | `takeScreenshot: TC-NNN` |

### Auth Pattern

For test cases requiring authentication, use Maestro's `env` section with credential variables:

```yaml
env:
  USERNAME: ${TEST_USERNAME}
  PASSWORD: ${TEST_PASSWORD}
commands:
  - tapOn: { text: "Email" }
  - inputText: ${USERNAME}
  - tapOn: { text: "Password" }
  - inputText: ${PASSWORD}
  - tapOn: { text: "Sign in" }
```

## Mobile Antipattern Guards

Beyond the generic 6 antipattern guards in the main SKILL.md, Mobile-specific generation must avoid these additional patterns:

### 1. Device-Dependent Tests Without Fixture Isolation

**Pattern**: Hardcoding device-specific values (screen dimensions, pixel coordinates, device model names) in test assertions.

**Why harmful**: Tests pass on one device but fail on others with different screen sizes. CI runners may use emulators with different configurations than local devices.

**Instead**: Use Maestro's device-agnostic selectors (`text`, `id`, `index`). Never reference pixel coordinates or device model strings. Let Maestro handle device abstraction.

### 2. Hardcoded Wait Durations

**Pattern**: Using `extendedWait` with fixed time durations (e.g., `extendedWait: { time: 5000 }`) to handle asynchronous loading.

**Why harmful**: Makes tests slow (over-waiting on fast devices) and flaky (under-waiting on slow devices). Duration-based waits do not adapt to network conditions or device performance.

**Instead**: Use `assertVisible` which waits automatically for the element to appear. Maestro's built-in assertion polling handles timing without explicit durations.

### 3. Tests Requiring Physical Device Capabilities

**Pattern**: Writing tests that depend on hardware sensors (camera, fingerprint reader, accelerometer, GPS) or physical device features that emulators may not support.

**Why harmful**: Tests fail in CI environments that use emulators without sensor simulation. Creates a gap between CI and local development environments.

**Instead**: Limit test scripts to interactions supported by Maestro's emulator/simulator environment (tap, swipe, text input, assertions). If a test case requires physical device capabilities, note it in a comment and mark it as requiring a physical device run -- do not generate an assertion that will fail in CI.

## Output

Mobile test scripts are written as Maestro YAML flows to `tests/e2e/features/<feature>/` following the strategy's template naming convention. Each flow includes a traceability comment linking back to the source test case ID.
