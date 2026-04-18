# Design System: Stripe

Purple gradients, weight-300 elegance. Fintech precision with warm details.

## Visual Theme & Atmosphere

Clean, confident, and weightless. Light surfaces with deep purple accents. Gradients are used purposefully — never decorative. Typography feels light (weight 300) even in headings, creating an airy sophistication.

## Color Palette

| Role | Value | Usage |
|------|-------|-------|
| Background | #f6f9fc | Page background |
| Surface | #ffffff | Cards, containers |
| Surface Alt | #f0f4fa | Alternating sections |
| Border | #e3e8ee | Dividers, card borders |
| Text Primary | #1a1f36 | Headings |
| Text Secondary | #697386 | Body, descriptions |
| Text Tertiary | #8792a2 | Captions, placeholders |
| Accent | #635bff | Primary actions, links |
| Accent Hover | #7a73ff | Hover states |
| Gradient Start | #635bff | Gradient left |
| Gradient End | #a259ff | Gradient right |
| Success | #30d158 | Confirmations |
| Warning | #f5a623 | Alerts |
| Error | #df1b41 | Errors, destructive |

## Typography

| Role | Font | Weight | Size |
|------|------|--------|------|
| Display | Inter | 300 | 48-56px |
| Heading 1 | Inter | 600 | 36-40px |
| Heading 2 | Inter | 600 | 24-28px |
| Heading 3 | Inter | 500 | 20px |
| Body | Inter | 400 | 16px |
| Body Small | Inter | 400 | 14px |
| Caption | Inter | 500 | 12px |
| Mono | JetBrains Mono | 400 | 14px |

Line height: 1.6 for body, 1.25 for headings. Generous paragraph spacing (24px).

## Components

### Buttons
- **Primary**: Solid gradient (635bff → a259ff), white text, rounded-lg (8px), px-5 py-2.5
- **Secondary**: White bg, border #e3e8ee, text #1a1f36
- **Ghost**: No border, text accent, hover underline
- **Hover**: Gradient brightens, slight scale(1.02), 200ms ease

### Cards
- Background: #ffffff, border: 1px solid #e3e8ee, rounded-xl (12px)
- Padding: 32px. Subtle shadow: 0 2px 8px rgba(0,0,0,0.04)
- Hover: shadow deepens, border softens

### Inputs
- Background: #ffffff, border: 1px solid #d8dee6, rounded-md (6px)
- Focus: border #635bff, ring 0 0 0 3px rgba(99,91,255,0.15)
- Height: 40px, padding: 0 12px

### Navigation
- Clean white bar, no backdrop blur needed on light
- Logo left, links center (weight 500, 14px), CTA right
- Active: purple text, no underline

## Layout

- Max content width: 1168px
- Grid: 12-column, 32px gap
- Section padding: 80px vertical
- Feature cards: 3 columns with 24px gap

## Depth & Elevation

| Level | Shadow | Usage |
|-------|--------|-------|
| 0 | none | Flat surfaces |
| 1 | 0 1px 3px rgba(0,0,0,0.06) | Cards |
| 2 | 0 4px 12px rgba(0,0,0,0.08) | Dropdowns |
| 3 | 0 12px 32px rgba(0,0,0,0.12) | Modals |

## Do's and Don'ts

| Do | Don't |
|-----|------|
| Use weight 300 for display text | Use bold for hero headings |
| Apply gradients to buttons only | Put gradients on text or backgrounds |
| Use generous white space | Crowd components |
| Keep purple as the sole accent | Mix warm and cool accent colors |
| Use soft shadows | Use hard shadows or borders for depth |

## Responsive Behavior

| Breakpoint | Behavior |
|------------|----------|
| >1080px | Full layout, 3-col grids |
| 720-1080px | 2-col grids, compact nav |
| <720px | Single column, hamburger, stacked CTAs |

## Signature Patterns

- **Gradient buttons**: 635bff→a259ff on primary CTAs, white text
- **Code blocks with terminal aesthetic**: Dark bg (#1a1f36), JetBrains Mono, purple syntax highlights
- **Illustration accents**: Abstract geometric shapes in gradient purple tones
- **Icon row features**: 4-col grid, icon + heading + description per cell
