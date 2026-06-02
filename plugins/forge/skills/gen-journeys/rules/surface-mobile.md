# Surface: Mobile

Mobile surface 适用于原生或跨平台移动应用程序（Android、iOS）。这是 Forge 的 **best-effort** 支持级别 -- 只生成 Maestro YAML 骨架和 deep link 测试，复杂场景标记 `manual-only`。

**Test type**: 移动端端到端测试 (Mobile E2E Test). Test type: 移动端端到端测试，通过 Maestro YAML / 设备自动化验证 UI 元素可见性、用户操作响应和 deep link 导航。Generated test code MUST use `@mobile-e2e` tags. This is one of the two surfaces where "e2e" terminology is correct (the other being Web).

## Detection Signals

| Signal | File Pattern | Dependency Pattern | Exclusion |
|--------|-------------|-------------------|-----------|
| Android | `AndroidManifest.xml` exists in project tree | Android UI framework dependencies (Activity, Fragment, Compose) in `build.gradle` or `build.gradle.kts` | None (AndroidManifest.xml is a strong Mobile signal) |
| iOS | `*.xcodeproj` or `*.xcworkspace` directory exists | UIKit, SwiftUI, or other iOS UI framework in source imports | None (xcodeproj is a strong Mobile signal) |
| React Native | `package.json` with `react-native` in dependencies | `react-native` + mobile-specific plugins | Not web (react-native targets mobile, not browser) |
| Flutter | `pubspec.yaml` exists with `flutter` dependency | Flutter framework in `pubspec.yaml` | None (Flutter is a strong Mobile signal) |

**Confidence Levels**:

- **High**: Platform-specific project file (`AndroidManifest.xml` or `*.xcodeproj`) + UI framework dependencies
- **Medium**: Cross-platform framework detected (React Native, Flutter) without native project files visible
- **Low**: Only mobile-adjacent files (e.g., `pubspec.yaml` without Flutter dependency)

**Disambiguation Rules**:

1. React Native vs web: `react-native` in dependencies = mobile. `react` without `react-native` = web. Do not confuse the two.
2. React Native Web: If both `react-native` and `react-dom` are present, check the primary target. If `react-native-web` is used for web, classify based on the primary user interaction model.
3. Progressive Web App (PWA): If the project is primarily a website with PWA capabilities, classify as web, not mobile.

## General Testing Principles

1. **Best-effort scope**: Mobile testing in Forge is intentionally limited. Focus on:
   - **Maestro YAML skeletons**: Generate structural test skeletons that developers fill in with specific selectors
   - **Deep link testing**: Verify that deep links navigate to the correct screen with correct parameters
   - **App lifecycle basics**: Launch, background, foreground
2. **Mark complex scenarios `manual-only`**: Any scenario involving:
   - Complex gesture sequences (multi-touch, swipe patterns)
   - Biometric authentication
   - Camera/AR interactions
   - Push notification handling
   - Platform-specific permissions flows

   These should be marked `manual-only` in generated tests, with a description of what manual verification is needed.

3. **No hard dependency on Maestro CLI**: Generated Maestro YAML should be syntactically valid even when the Maestro CLI is not installed. The test generation step does not require Maestro to be present.

4. **Platform-aware assertions**: When possible, note platform-specific behavior differences (Android vs iOS) in generated test descriptions.

## Test Strategy Guidance

**Test Level Emphasis**: Journey skeleton + deep link (minimal Contract testing)

Mobile testing in Forge prioritizes Journey-level test skeletons over Contract tests. The reasoning:

- Mobile UI automation is inherently fragile (selector changes, OS updates, device fragmentation)
- Contract-level testing for mobile is better served by platform-native tools (XCTest, Espresso) that developers write manually
- Forge's value-add is generating the structural skeleton and deep link coverage, not replacing native testing tools

**Execution Model**: Maestro YAML

- Generate Maestro-compatible YAML test files
- Use Maestro's `assertVisible`, `tapOn`, `inputText` commands
- Deep link tests use Maestro's `openLink` command
- Each generated test includes a header comment indicating whether it's a complete test or a skeleton requiring manual completion

**Environment Readiness Checks**:

| Check | How to Verify |
|-------|--------------|
| Maestro CLI installed (optional) | `which maestro` returns 0 |
| Android emulator or iOS simulator available | `maestro devices` lists at least one device |
| App binary available | APK or IPA exists at expected path |
| Deep link scheme registered | Android: `intent-filter` in AndroidManifest; iOS: URL scheme in Info.plist |

**Note**: Environment readiness checks are informational only. Missing Maestro CLI does not block test generation -- it only blocks test execution.

## Required Outcome Reference

Mobile surface does not have mandatory derived Outcomes in the same way as CLI/API/Web due to its best-effort nature. Instead, use these common patterns as reference:

**Common Mobile test patterns**:

- **app-launch**: App starts and displays the expected initial screen. Maestro: `assertVisible` on root element.
- **deep-link-navigation**: Deep link opens the app to the correct screen with correct content. Maestro: `openLink` + `assertVisible` on target screen element.
- **deep-link-invalid**: Invalid deep link shows error state or falls back to home screen.
- **app-background-foreground**: App resumes correctly after being backgrounded. State is preserved or gracefully restored.

**Patterns to mark `manual-only`**:

- Complex form submission with validation
- Multi-step wizard/flow completion
- Offline mode behavior
- Push notification interaction
- Permission grant/deny flows
- Biometric authentication
- Camera/photo capture
