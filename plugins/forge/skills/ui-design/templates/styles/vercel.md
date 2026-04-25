# Design System: Vercel

Black and white precision, Geist font. Developer platform aesthetic.

## Visual Theme & Atmosphere

Monochrome precision. Every element earns its place. Dark surfaces with stark white text create a cinematic, terminal-like atmosphere. Negative space is a first-class element.

## Color Palette

| Role | Value | Usage |
|------|-------|-------|
| Background | #000000 | Primary surface |
| Surface | #111111 | Cards, elevated containers |
| Surface Alt | #171717 | Secondary backgrounds |
| Border | #262626 | Subtle dividers |
| Text Primary | #ededed | Headings, body |
| Text Secondary | #888888 | Descriptions, metadata |
| Accent | #0070f3 | Links, CTAs, interactive |
| Accent Hover | #3291ff | Hover states |
| Success | #50e3c2 | Positive feedback |
| Error | #ee0000 | Destructive actions |

## Typography

| Role | Font | Weight | Size |
|------|------|--------|------|
| Display | Geist Sans | 700 | 48-64px |
| Heading 1 | Geist Sans | 700 | 36-40px |
| Heading 2 | Geist Sans | 600 | 28-32px |
| Heading 3 | Geist Sans | 600 | 20-24px |
| Body | Geist Sans | 400 | 16px |
| Caption | Geist Sans | 400 | 14px |
| Mono | Geist Mono | 400 | 14px |

Line height: 1.5 for body, 1.2 for headings. Tight letter-spacing on display text.

## Components

### Buttons
- **Primary**: Solid accent (#0070f3), white text, rounded-md (6px), px-4 py-2
- **Secondary**: Ghost — transparent bg, border #333, white text
- **Hover**: Brighten 15%, no shadows. Smooth 150ms transition.

### Cards
- Background: #111111, border: 1px solid #262626, rounded-xl (12px)
- Padding: 24px. No box-shadow by default.
- Hover: border brightens to #444, slight translate-y(-1px)

### Inputs
- Background: #111111, border: 1px solid #333, rounded-lg (8px)
- Focus: border becomes accent (#0070f3), subtle ring
- Text: Geist Sans, 14px

### Navigation
- Horizontal top bar, bg #000/80 with backdrop-blur
- Logo left, links center, CTA right
- Active link: white text + 1px bottom accent line

## Layout

- Max content width: 1200px
- Grid: 12-column, 24px gap
- Section padding: 96px vertical (desktop), 64px (mobile)
- Card grid: 3 columns (desktop), 1 column (mobile)

## Depth & Elevation

No drop shadows. Depth conveyed through:
- Surface brightness (#000 > #111 > #171)
- Border opacity
- Z-index layers: base(0), card(10), nav(20), modal(30), toast(40)

## Do's and Don'ts

| Do | Don't |
|----|-------|
| Use pure black for backgrounds | Add colored tints to surfaces |
| Rely on spacing for hierarchy | Use drop shadows for elevation |
| Use accent color sparingly | Apply gradients to text |
| Keep hover effects subtle | Use more than 2 font weights |
| Use Geist Mono for code | Mix serif and sans-serif |

## Responsive Behavior

| Breakpoint | Behavior |
|------------|----------|
| >1024px | Full layout, 3-col grids |
| 768-1024px | 2-col grids, reduced padding |
| <768px | Single column, stacked navigation, hamburger menu |
