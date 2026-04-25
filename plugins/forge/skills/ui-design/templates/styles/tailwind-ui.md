# Design System: Tailwind UI

Professional warmth, indigo accent, scene-driven. The gold standard for SaaS UI.

## Visual Theme & Atmosphere

Professional and approachable. Clean surfaces with enough color to feel alive — sectioned backgrounds alternate between white and soft gray/indigo tints. Typography is generous and readable. Illustrations and icons add warmth without clutter. Feels like a well-run company.

## Color Palette

### Brand Colors (Indigo)

| Role | Value | Usage |
|------|-------|-------|
| Primary 50 | #eef2ff | Lightest accent bg |
| Primary 100 | #e0e7ff | Accent surfaces |
| Primary 200 | #c7d2fe | Borders on accent |
| Primary 400 | #818cf8 | Inline highlights |
| Primary 500 | #6366f1 | Default accent |
| Primary 600 | #4f46e5 | Buttons, links |
| Primary 700 | #4338ca | Hover / pressed |
| Primary 800 | #3730a3 | Dark accent |
| Primary 900 | #312e81 | Darkest accent |

### Neutral (Slate)

| Role | Value | Usage |
|------|-------|-------|
| Background | #ffffff | Primary surface |
| Background Alt | #f8fafc (slate-50) | Alternating sections |
| Background Warm | #f9fafb (gray-50) | Subtle warmth sections |
| Surface | #ffffff | Cards |
| Border | #e2e8f0 (slate-200) | Dividers |
| Border Dark | #cbd5e1 (slate-300) | Emphasized borders |
| Text Primary | #0f172a (slate-900) | Headings |
| Text Secondary | #475569 (slate-600) | Body |
| Text Tertiary | #94a3b8 (slate-400) | Captions |

### Semantic

| Role | Value | Usage |
|------|-------|-------|
| Success | #059669 (emerald-600) | Positive actions |
| Warning | #d97706 (amber-600) | Alerts |
| Error | #dc2626 (red-600) | Destructive / errors |

## Typography

| Role | Font | Weight | Size | Line Height |
|------|------|--------|------|-------------|
| Display | Inter | 800 | 48-60px | 1.1 |
| H1 | Inter | 700 | 36-48px | 1.2 |
| H2 | Inter | 600 | 30px | 1.25 |
| H3 | Inter | 600 | 24px | 1.3 |
| H4 | Inter | 500 | 20px | 1.4 |
| Body | Inter | 400 | 16px | 1.625 (1.625) |
| Body Small | Inter | 400 | 14px | 1.5 |
| Caption | Inter | 500 | 12px | 1.5 |
| Mono | JetBrains Mono | 400 | 14px | 1.5 |

Font stack: `Inter, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif`

Generous tracking on uppercase text: `letter-spacing: 0.05em`. Body line-height 1.625 for comfortable reading.

## Components

### Buttons
- **Primary**: bg primary-600, text white, rounded-lg (8px), px-4 py-2.5, text-sm font-medium
- **Secondary**: bg white, border slate-300, text slate-700, shadow-sm
- **Soft**: bg primary-50, text primary-700, no border
- **Ghost**: transparent, text slate-600, hover bg slate-50
- **Sizes**: xs(px-2.5 py-1.5 text-xs), sm(px-3 py-2), default(px-4 py-2.5), lg(px-5 py-3)
- **Hover**: primary → bg primary-700, secondary → bg slate-50, shadow transitions 150ms
- **Icon buttons**: square aspect ratio, rounded-lg

### Cards
- bg white, border 1px solid slate-200, rounded-xl (12px)
- Shadow: 0 1px 3px rgba(0,0,0,0.1)
- Hover (interactive): shadow 0 4px 6px rgba(0,0,0,0.07)
- Padding: 24px (p-6), card header with border-bottom for sections

### Inputs
- bg white, border slate-300, rounded-md (6px), shadow-sm
- h-10, px-3, text-sm
- Focus: border primary-500, ring 2px primary-200
- Labels: text-sm font-medium text-slate-700, mb-1.5

### Navigation
- White bg, border-b slate-200, h-16
- Logo left, nav links center (text-sm font-medium text-slate-600), CTA right
- Mobile: hamburger with slide-in panel, bg white
- Active link: text primary-600 + bottom border

### Badges
- rounded-full, px-2.5 py-0.5, text-xs font-medium
- Colors: green(bg emerald-50 text emerald-700), blue, red, yellow, gray

## Layout

- Max content: 1280px (wider than Apple, narrower than Vercel)
- Grid: 12-column, 24-32px gap
- Section padding: 80px vertical (py-20)
- Feature sections: alternating white / slate-50 backgrounds
- CTA sections: primary-600 bg with white text for contrast

## Depth & Elevation

| Level | Shadow | Usage |
|-------|--------|-------|
| 0 | none | Flat sections |
| 1 | 0 1px 2px rgba(0,0,0,0.05) | Inputs |
| 2 | 0 1px 3px rgba(0,0,0,0.1) | Cards |
| 3 | 0 4px 6px -1px rgba(0,0,0,0.1) | Dropdowns |
| 4 | 0 10px 15px -3px rgba(0,0,0,0.1) | Modals |

## Do's and Don'ts

| Do | Don't |
|-----|------|
| Use indigo as the sole accent color | Mix multiple accent hues |
| Alternate section backgrounds | Keep every section the same bg |
| Use shadow-sm on inputs/cards | Skip shadows entirely (looks flat) |
| Use weight 600 for headings | Go below 500 for any heading |
| Use slate (cool gray) for neutrals | Mix warm gray and cool gray |
| Add subtle ring on focus | Use only border-color change for focus |

## Responsive Behavior

| Breakpoint | Tailwind | Behavior |
|------------|----------|----------|
| <640px | sm | Single column, stacked nav, full-width CTAs |
| 640-768px | md | 2-col grids |
| 768-1024px | lg | Sidebar nav, 3-col feature grids |
| >1024px | xl | Full layout, 4-col grids possible |

## Signature Patterns

- **Alternating sections**: white → slate-50 → white, creates visual rhythm
- **Gradient CTAs**: primary-600 to primary-700, rounded-lg, white text
- **Icon + heading combos**: Heroicon outline 24px + heading, for feature blocks
- **Stat blocks**: Big number (text-4xl font-bold) + label (text-sm text-slate-600)
- **Testimonial cards**: Quote + avatar + name + role, bg white with shadow
