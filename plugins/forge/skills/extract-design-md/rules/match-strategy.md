# Match Strategy

Use `AskUserQuestion` to let the user choose a generation strategy:

| Option | Description |
|--------|-------------|
| Match closest built-in style, customize on top | Identify the closest built-in style (vercel/shadcn/tailwind-ui/stripe/apple), override differences with extracted actual tokens |
| Fully custom from web app extraction | Generate an independent DESIGN.md entirely from analysis results, no built-in style reference |

**If "match built-in" is chosen:**

Based on visual analysis, match against built-in style characteristics defined in `rules/style-matching.md` to identify the closest built-in style, then read the corresponding style file from the ui-design skill: `ui-design/templates/styles/<name>.md` (resolve relative to the skills parent directory)
