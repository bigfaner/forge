# Design System: Apple

Premium white space, SF Pro, cinematic imagery. Consumer-grade elegance.

## Visual Theme & Atmosphere

Radical simplicity. Content is king — UI chrome disappears. Large, high-quality imagery dominates. Typography is bold and architectural. Every screen feels like a billboard: one message, maximum impact.

## Color Palette

| Role | Value | Usage |
|------|-------|-------|
| Background | #ffffff | Primary surface |
| Background Alt | #f5f5f7 | Alternating sections |
| Surface | #ffffff | Cards |
| Border | #d2d2d7 | Dividers |
| Text Primary | #1d1d1f | Headings |
| Text Secondary | #6e6e73 | Body |
| Text Tertiary | #86868b | Captions |
| Accent | #0071e3 | Links, CTAs |
| Accent Hover | #0077ed | Hover states |
| Accent Light | #e8f4fd | Light accent bg |
| Success | #30d158 | Positive |
| Error | #ff3b30 | Destructive |

## Typography

| Role | Font | Weight | Size | Tracking |
|------|------|--------|------|----------|
| Display | SF Pro Display | 700 | 56-80px | -0.02em |
| Heading 1 | SF Pro Display | 600 | 40-48px | -0.015em |
| Heading 2 | SF Pro Display | 600 | 28-32px | -0.01em |
| Heading 3 | SF Pro Display | 600 | 21-24px | normal |
| Body | SF Pro Text | 400 | 17px | normal |
| Body Small | SF Pro Text | 400 | 14px | normal |
| Caption | SF Pro Text | 500 | 12px | 0.02em |

Line height: 1.47 for body (17px), 1.08 for display. SF Pro is system font on Apple devices; use -apple-system, 'Segoe UI', Roboto as fallback stack.

## Components

### Buttons
- **Primary**: Accent (#0071e3), white text, rounded-xl (980px pill or 12px rounded), px-5 py-3
- **Secondary**: Ghost — accent text, no border, hover underline
- **CTA Large**: px-8 py-4, font-size 17px, weight 600
- **Hover**: Brighten 8%, no shadows, 200ms ease

### Cards
- Background: #ffffff, no border, rounded-2xl (20px)
- Padding: 24-32px. Shadow: 0 2px 12px rgba(0,0,0,0.06)
- Hover: scale(1.01), shadow softens

### Inputs
- Background: #f5f5f7, no border, rounded-lg (10px)
- Focus: ring 0 0 0 4px rgba(0,113,227,0.2), bg #fff
- Height: 44px, padding: 0 16px

### Navigation
- Semi-transparent: rgba(255,255,255,0.72) with backdrop-blur(20px)
- Slim: 44px height. Logo left, links center (weight 400, 12px uppercase tracking), CTA right
- Sticky, border-bottom 1px solid rgba(0,0,0,0.08)

## Layout

- Max content width: 980px (narrower = more premium)
- Grid: 12-column, 20px gap
- Section padding: 100-120px vertical (generous)
- Hero: full-width with max-text-width 680px centered

## Depth & Elevation

| Level | Shadow | Usage |
|-------|--------|-------|
| 0 | none | Flat sections |
| 1 | 0 2px 8px rgba(0,0,0,0.04) | Cards |
| 2 | 0 4px 16px rgba(0,0,0,0.08) | Modals |
| 3 | 0 8px 32px rgba(0,0,0,0.12) | Overlays |

Real depth comes from layering full-bleed photography with text, not from shadows.

## Do's and Don'ts

| Do | Don't |
|-----|------|
| Let imagery dominate | Use decorative borders or frames |
| Use extreme font sizes for display | Crowd hero sections with text |
| Keep navigation minimal | Add icons to nav links |
| Use one message per section | Mix multiple accent colors |
| Use system fonts (SF Pro) | Load custom display fonts |
| Embrace generous white space | Fill every pixel |

## Responsive Behavior

| Breakpoint | Behavior |
|------------|----------|
| >1068px | Full layout, max-width 980px |
| 734-1068px | Reduced padding, 2-col grids |
| <734px | Single column, stacked, larger touch targets |
