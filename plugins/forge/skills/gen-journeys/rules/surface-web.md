# Surface: Web (Web User Interface)

Web surface 适用于基于浏览器的 Web 应用程序（React、Vue、Svelte 等）。测试重点是用户交互流程、状态转换、可访问性和浏览器自动化。

**Test type**: Web 端到端测试 (Web E2E Test). Test type: Web 端到端测试，通过浏览器自动化验证 DOM 元素可见性、用户操作响应和页面状态转换。Generated test code MUST use `@web-e2e` tags. This is one of the two surfaces where "e2e" terminology is correct (the other being Mobile).

## Detection Signals

| Signal | File Pattern | Dependency Pattern | Exclusion |
|--------|-------------|-------------------|-----------|
| React SPA | `package.json` exists at project root | `react` in `dependencies` + browser DOM entry (`document.getElementById`, `createRoot`, or equivalent in source) | None (React is a strong Web signal) |
| Vue SPA | `package.json` exists at project root | `vue` in `dependencies` + browser DOM entry (`createApp`, `mount('#app')`, or equivalent) | None (Vue is a strong Web signal) |
| Svelte SPA | `package.json` exists at project root | `svelte` in `dependencies` or `devDependencies` + browser DOM entry | None (Svelte is a strong Web signal) |

**Confidence Levels**:

- **High**: `package.json` with frontend framework + browser DOM entry point + typical frontend file structure (`src/`, `pages/`, `components/`)
- **Medium**: `package.json` with frontend framework but no clear DOM entry point detected
- **Low**: Only frontend-adjacent files detected (e.g., HTML files without framework dependency)

**Disambiguation Rules**:

1. If `package.json` contains both a frontend framework and a server framework (`express`/`fastify`), this is likely a full-stack application. Detect the **user-facing surface**: if the primary user interaction is through a browser, classify as web. If the server is purely an API backend consumed by other clients, classify as api.
2. `ink` in Node.js dependencies is NOT web -- it renders to the terminal (classify as tui).
3. Server-side rendering (SSR) frameworks like Next.js, Nuxt, SvelteKit: classify as web because the user interaction model is browser-based.
4. Static site generators (Astro, Hugo, etc.): if the output is interactive pages with JavaScript, classify as web. If purely static content, this may not need automated testing through this pipeline.

## General Testing Principles

1. **Browser automation**: Web tests use browser automation frameworks (Playwright, Cypress, Selenium, etc.). The specific framework is defined by the project's Convention file, not by this surface rule.
2. **User-centric assertions**: Test from the user's perspective -- what they see and interact with. Avoid asserting internal component state or implementation details.
3. **State transitions**: Verify that UI state changes correctly in response to user actions:
   - Form submissions update displayed data
   - Navigation changes the visible page/view
   - Loading states appear during async operations
4. **Accessibility**: Test that interactive elements are reachable via keyboard navigation and that ARIA labels are present for dynamic content.
5. **Async handling**: Web tests must account for:
   - Network request latency (use appropriate wait strategies, not fixed timeouts)
   - Animation completion (wait for elements to become stable before asserting)
   - Client-side routing (wait for page transition to complete)

## Test Strategy Guidance

**Test Level Emphasis**: Balanced 50/50 (Contract 50% / Journey smoke 50%)

Web applications benefit equally from Contract-level tests (individual component/interaction behavior) and Journey smoke tests (end-to-end user workflows). The visual and interactive nature of Web makes both levels important.

**Execution Model**: Browser automation

- Launch a browser instance (headless by default) controlled by the Convention-defined framework
- Each test navigates to the application and interacts with DOM elements
- Use `data-testid` or accessible selectors for element targeting
- Clean up browser state between tests (clear cookies, localStorage, etc.)

**Environment Readiness Checks**:

| Check | How to Verify |
|-------|--------------|
| Dev server starts | `npm run dev` / equivalent starts and responds on expected port |
| Browser automation framework installed | `npx playwright install` or equivalent completes |
| Application loads | HTTP GET to dev server root returns 200 |
| Test database seeded | Required test data is available |

**Why balanced 50/50**: Unlike CLI (where Contract tests are highly reliable due to subprocess isolation), Web Journey tests provide unique value by validating the full rendering pipeline, client-side routing, and browser-specific behaviors that Contract tests alone cannot catch.

## Required Outcome Reference

**Mandatory derived Outcomes** (must be considered for every Web Journey):

- **validation-error**: User submits a form with invalid data. Example: required field left empty, email format invalid, numeric field has non-numeric input. Assert: error message displayed near the relevant field, form is not submitted, user can correct and retry.
- **session-expired**: User's session has expired during an active workflow. Example: user fills out a long form, session times out, user submits. Assert: appropriate redirect or message shown, unsaved data is either preserved or user is warned about data loss, login flow is accessible from the expired state.

**Additional common Web boundary Outcomes**:

- **network-error**: API request fails due to network issues. Assert: error message displayed, retry option available, no data loss.
- **loading-state**: Async operation in progress. Assert: loading indicator visible, UI remains responsive (or shows appropriate blocking state).
- **navigation-guard**: User attempts to leave a page with unsaved changes. Assert: confirmation dialog shown, user choice respected.
- **responsive-layout**: Page layout adapts to viewport size. Assert: content remains accessible and usable at different breakpoints.
- **concurrent-edit**: Another user has modified the same resource. Assert: conflict notification shown, merge or overwrite option provided.
