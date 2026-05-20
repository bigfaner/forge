# Design Style Selection

Select design style based on the platform identified from the Navigation Architecture.

## For Web / Mobile Platform

Select design style by priority:

### Priority 1: User-provided DESIGN.md

Check the following locations; use directly if found:

- Project root `DESIGN.md`
- Feature directory `docs/features/<slug>/ui/style.md`

If the user specifies a custom path, prioritize that.

### Priority 2: Built-in Web/Mobile Styles

If no user-provided DESIGN.md, let the user choose from 5 built-in styles:

| Style           | Vibe                     | Best for                                     |
| --------------- | ------------------------ | -------------------------------------------- |
| **Vercel**      | Black & white minimal, developer-tool feel | Developer platforms, CLI tools, technical docs |
| **Shadcn**      | Zinc neutral, functional minimalism | SaaS, admin panels, tool applications        |
| **Tailwind UI** | Indigo primary, professional warmth | General SaaS, marketing pages, enterprise    |
| **Stripe**      | Purple gradients, light elegance | Fintech, brand sites, payment products       |
| **Apple**       | Generous whitespace, image-driven, premium | Consumer products, brand sites, marketing     |

Use `AskUserQuestion` tool for user selection with brief descriptions.

Built-in style files located at: `templates/styles/{vercel,shadcn,tailwind-ui,stripe,apple}.md`

### Priority 3: Clone More Styles from Repo

If the 5 built-in styles are insufficient, clone additional styles from the awesome-design-md repo:

```bash
# Clone to a temp directory (not into project)
git clone --depth 1 git@github.com:VoltAgent/awesome-design-md.git /tmp/awesome-design-md

# Then use npx to fetch a specific site's DESIGN.md:
npx getdesign@latest add <site-name>
```

> Note: The repo's DESIGN.md files are hosted externally at getdesign.md. The `npx getdesign@latest` CLI fetches them.

## For TUI Platform

TUI has its own design style selection, separate from web/mobile.

### Priority 1: User-provided DESIGN.md

Same check as web: project root `DESIGN.md` or feature directory `docs/features/<slug>/ui/style.md`.

If a user-provided DESIGN.md is found, use it directly as the TUI design system. Skip theme selection.

### Priority 2: Built-in TUI Themes

If no user-provided DESIGN.md, let the user choose from 2 built-in TUI themes:

| Theme              | Vibe                                         | Best for                                      |
| ------------------ | -------------------------------------------- | --------------------------------------------- |
| **Modern Dark**    | 256-color, box-drawing chars, dark bg, high contrast, compact density | Most CLI tools, developer dashboards, monitoring apps |
| **Minimal ASCII**  | 16-color, pure ASCII chars, default terminal bg, loose density, max compatibility | Minimal tools, SSH sessions, legacy terminals, CI output |

Use `AskUserQuestion` tool for user selection with brief descriptions.

Built-in TUI theme files located at: `templates/styles/{modern-dark-tui,minimal-ascii-tui}.md`

## Using the Selected Style

Inline the selected style content as the `Design System` section in `ui-design.md`. All subsequent component designs must follow the style's color, typography, and component specifications.

## Multi-Platform Features

When the PRD declares multiple platforms (e.g., `platform: web, tui`):
- Select a style for EACH platform independently (web/mobile style + TUI theme)
- Each platform will produce its own `ui-design-{platform}.md` file
- Output file naming:
  - Single platform (web): `ui/ui-design.md`
  - Single platform (tui): `ui/ui-design-tui.md`
  - Multi-platform: `ui/ui-design-web.md` + `ui/ui-design-tui.md`
