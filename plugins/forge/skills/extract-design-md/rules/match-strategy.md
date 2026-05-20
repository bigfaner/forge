# Match Strategy

Use `AskUserQuestion` to let the user choose a generation strategy:

| Option | Description |
|--------|-------------|
| Match closest built-in style, customize on top | Identify the closest built-in style (vercel/shadcn/tailwind-ui/stripe/apple), override differences with extracted actual tokens |
| Fully custom from web app extraction | Generate an independent DESIGN.md entirely from analysis results, no built-in style reference |

**If "match built-in" is chosen:**

Based on visual analysis, match against these characteristics to identify the closest built-in style:

| Built-in Style | Identifying Characteristics |
|---------------|----------------------------|
| Vercel | Black background, Geist font, no shadows, border depth |
| Shadcn | Zinc neutrals, CSS variables, dark mode support, Tailwind spacing |
| Tailwind UI | Indigo primary, white background, shadow-sm system, Inter font |
| Stripe | Purple gradient buttons, light gray background (#f6f9fc), weight-300 display |
| Apple | Pure white background, generous whitespace, SF Pro, rounded capsule buttons |

Read the corresponding built-in style file: `${CLAUDE_SKILL_DIR}/../ui-design/templates/styles/<name>.md`
