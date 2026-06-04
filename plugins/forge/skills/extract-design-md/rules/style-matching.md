# Style Matching Reference

Match analysis results against built-in style characteristics to identify the closest match.

## Web Built-in Styles

| Style | Background | Primary Accent | Font | Key Differentiator |
|-------|-----------|---------------|------|--------------------|
| Vercel | #000000 (pure black) | #0070f3 (blue) | Geist Sans | Monochrome dark, no shadows, border depth via surface brightness |
| Shadcn | #ffffff / CSS variables | zinc-900 (#18181b) | Inter / System | Zinc neutrals, CSS variable system, dark mode toggle, no shadow |
| Tailwind UI | #ffffff | #4f46e5 (indigo-600) | Inter | Indigo primary, shadow-sm system, alternating white/slate sections |
| Stripe | #f6f9fc (cool gray) | #635bff → #a259ff (purple gradient) | Inter | Gradient buttons, weight-300 display, fintech precision |
| Apple | #ffffff (pure white) | #0071e3 (blue) | SF Pro | Generous whitespace, cinematic imagery, rounded capsule buttons |

## TUI Built-in Themes

| Theme | Background | Color Depth | Character Set | Density |
|-------|-----------|-------------|---------------|---------|
| modern-dark-tui | Dark gray (xterm-235) | 256-color (xterm-256) | Box-drawing + block elements (Unicode) | Compact |
| minimal-ascii-tui | Default terminal | 16-color (standard ANSI) | Pure ASCII only | Loose |

## Matching Decision Guide

When multiple styles share characteristics, use these tiebreakers in order:

1. **Background color**: Exact or near-exact match is strongest signal (Vercel = pure black, Apple = pure white, Stripe = cool gray #f6f9fc)
2. **Accent color family**: Blue (Vercel/Apple) vs Purple (Stripe) vs Indigo (Tailwind UI) vs Zinc neutral (Shadcn)
3. **Font family**: Geist = Vercel, SF Pro = Apple, CSS variables = Shadcn, Inter with indigo = Tailwind UI, Inter with gradient = Stripe
4. **Shadow/elevation system**: No shadows = Vercel/Shadcn, shadow-sm = Tailwind UI, subtle soft shadows = Stripe/Apple
5. **Dark mode**: CSS variable light/dark toggle = Shadcn, dark-only = Vercel, light-only = Stripe/Apple/Tailwind UI

After identifying the closest style, read the corresponding style file from the ui-design skill: `ui-design/templates/styles/<name>.md` (resolve relative to the skills parent directory).
