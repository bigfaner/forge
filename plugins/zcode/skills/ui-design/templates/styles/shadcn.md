# Design System: Shadcn

Neutral zinc palette, Radix + Tailwind primitives. Functional minimalism.

## Visual Theme & Atmosphere

Clean and utilitarian. Zero decoration — every visual element serves a function. Surfaces are flat with subtle 1px borders. The aesthetic is "invisible UI": users focus on content, not chrome. Light mode default with seamless dark mode toggle.

## Color Palette

Defined as CSS variables for light/dark mode switching.

### Light Mode

| Role | Variable | Value |
|------|----------|-------|
| Background | --background | #ffffff |
| Foreground | --foreground | #09090b (zinc-950) |
| Card | --card | #ffffff |
| Card Foreground | --card-foreground | #09090b |
| Primary | --primary | #18181b (zinc-900) |
| Primary Foreground | --primary-foreground | #fafafa |
| Secondary | --secondary | #f4f4f5 (zinc-100) |
| Secondary Foreground | --secondary-foreground | #18181b |
| Muted | --muted | #f4f4f5 |
| Muted Foreground | --muted-foreground | #71717a (zinc-500) |
| Accent | --accent | #f4f4f5 |
| Accent Foreground | --accent-foreground | #18181b |
| Destructive | --destructive | #ef4444 |
| Border | --border | #e4e4e7 (zinc-200) |
| Input | --input | #e4e4e7 |
| Ring | --ring | #18181b |

### Dark Mode

| Role | Variable | Value |
|------|----------|-------|
| Background | --background | #09090b |
| Foreground | --foreground | #fafafa |
| Card | --card | #09090b |
| Primary | --primary | #fafafa |
| Primary Foreground | --primary-foreground | #18181b |
| Secondary | --secondary | #27272a (zinc-800) |
| Muted | --muted | #27272a |
| Muted Foreground | --muted-foreground | #a1a1aa (zinc-400) |
| Border | --border | #27272a |
| Ring | --ring | #d4d4d8 (zinc-300) |

## Typography

| Role | Font | Weight | Size | Line Height |
|------|------|--------|------|-------------|
| H1 | Inter / System | 700 | 36px (2.25rem) | 1.2 |
| H2 | Inter / System | 600 | 30px (1.875rem) | 1.25 |
| H3 | Inter / System | 600 | 24px (1.5rem) | 1.3 |
| H4 | Inter / System | 600 | 20px (1.25rem) | 1.35 |
| Body | Inter / System | 400 | 16px (1rem) | 1.5 |
| Small | Inter / System | 500 | 14px (0.875rem) | 1.5 |
| Muted | Inter / System | 400 | 14px | 1.5 |
| Mono | JetBrains Mono | 400 | 14px | 1.5 |

Font stack: `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif`

## Components

### Buttons
- **Default**: bg primary, text primary-foreground, rounded-md (6px), h-10, px-4 py-2, text-sm font-medium
- **Secondary**: bg secondary, text secondary-foreground
- **Outline**: border input, bg transparent, hover bg accent
- **Ghost**: transparent, hover bg accent
- **Destructive**: bg destructive, text white
- **Sizes**: sm(h-9 px-3), default(h-10 px-4), lg(h-11 px-8), icon(h-10 w-10)
- **Hover**: opacity shift, no shadows. 150ms transition.

### Cards
- bg card, text card-foreground, rounded-lg (8px), border 1px solid var(--border)
- Padding: 24px (p-6). No shadow.
- Variant: **CardHeader** (p-6 pb-0), **CardContent** (p-6 pt-0), **CardFooter** (p-6 pt-0)

### Inputs
- bg transparent, border 1px solid var(--input), rounded-md (6px)
- h-10, px-3 py-2, text-sm
- Focus: ring 2px var(--ring), ring-offset 2px var(--background)
- Placeholder: var(--muted-foreground)
- Disabled: opacity-50, cursor-not-allowed

### Badges
- rounded-full (9999px), px-2.5 py-0.5, text-xs font-medium
- Variants: default(bg-primary), secondary(bg-secondary), outline(border), destructive(bg-destructive)

### Table
- Clean grid, header bg-muted, row hover bg-muted/50
- Borders between rows only (border-b), no outer border

### Dialog / Modal
- Overlay: bg-black/80, backdrop-blur-sm
- Content: bg-background, rounded-lg, shadow-lg, max-w-lg
- Animation: fade-in + scale(0.95→1)

## Layout

- No enforced max-width — fluid by default
- Spacing scale: 4px base (Tailwind default)
  - 1 = 4px, 2 = 8px, 3 = 12px, 4 = 16px, 6 = 24px, 8 = 32px, 12 = 48px
- Common patterns: `p-6` (24px), `gap-4` (16px), `space-y-6` between sections

## Depth & Elevation

No shadow system. Depth is achieved through:
- Background opacity variation (muted vs background vs card)
- Border visibility
- Overlay with backdrop-blur for modals/popovers

## Do's and Don'ts

| Do | Don't |
|-----|------|
| Use zinc neutrals exclusively | Mix warm (amber) and cool (blue) tones |
| Keep borders 1px solid | Use thick or double borders |
| Use rounded-md for inputs, rounded-lg for cards | Mix border radius values freely |
| Support dark mode via CSS variables | Hard-code color values |
| Use Tailwind spacing scale | Invent custom spacing values |
| Use subtle opacity for disabled states | Use display:none for disabled |
| Use `font-medium` (500) for emphasis | Go above weight 700 |

## Responsive Behavior

| Breakpoint | Tailwind | Behavior |
|------------|----------|----------|
| <640px | sm | Stack everything, full-width cards |
| 640-768px | md | 2-col grids where appropriate |
| 768-1024px | lg | Sidebar nav possible |
| >1024px | xl | Full layout |

## CSS Variable Setup

```css
:root {
  --background: 0 0% 100%;
  --foreground: 240 10% 3.9%;
  --card: 0 0% 100%;
  --card-foreground: 240 10% 3.9%;
  --primary: 240 5.9% 10%;
  --primary-foreground: 0 0% 98%;
  --secondary: 240 4.8% 95.9%;
  --secondary-foreground: 240 5.9% 10%;
  --muted: 240 4.8% 95.9%;
  --muted-foreground: 240 3.8% 46.1%;
  --accent: 240 4.8% 95.9%;
  --accent-foreground: 240 5.9% 10%;
  --destructive: 0 84.2% 60.2%;
  --border: 240 5.9% 90%;
  --input: 240 5.9% 90%;
  --ring: 240 5.9% 10%;
  --radius: 0.5rem;
}

.dark {
  --background: 240 10% 3.9%;
  --foreground: 0 0% 98%;
  --card: 240 10% 3.9%;
  --primary: 0 0% 98%;
  --primary-foreground: 240 5.9% 10%;
  --secondary: 240 3.7% 15.9%;
  --muted: 240 3.7% 15.9%;
  --muted-foreground: 240 5% 64.9%;
  --accent: 240 3.7% 15.9%;
  --destructive: 0 62.8% 30.6%;
  --border: 240 3.7% 15.9%;
  --input: 240 3.7% 15.9%;
  --ring: 240 4.9% 83.9%;
}
```
