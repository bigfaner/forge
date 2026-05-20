# Page Exploration Rules

## Element Extraction

1. Get page title
2. Extract elements from snapshot with filters:
   - Exclude elements already in `layout.elements` (matched by `role + name`)
   - `role` in the full ARIA role set:
     {button, link, heading, textbox, checkbox, radio, combobox, tab, dialog, alert, navigation, search, form, menuitem, switch,
      table, row, cell, columnheader, rowheader,
      grid, listbox, option, list, listitem,
      tooltip, progressbar, meter, slider, spinbutton,
      status, log, marquee, timer,
      img, separator, group, region,
      feed, article, figure, caption}
   - `name` is non-empty
3. For each element, record full attributes:
   - Common: `{ role, name }`
   - heading: additionally record `level`
   - textbox/combobox: additionally record `label` (associated label text) and `placeholder`

## Dynamic State Exploration

Explore dynamic states triggered by multiple interaction types:

**1. Click triggers**: For elements with role=button/tab/disclosure and non-empty name:

```
ab('click @eN')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (compare with base snapshot)
ab('press Escape')  // or ab('click @close_btn') to reset
```

**2. Hover triggers**: For elements with tooltip or aria-describedby attributes:

```
ab('hover @eN')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (tooltip content)
ab('move 0 0')  // move mouse away to reset
```

**3. Scroll triggers**: For elements with role=feed, role=list with overflow, or scrollable containers:

```
ab('scroll @eN down')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (lazy-loaded content)
ab('scroll @eN up')  // scroll back to reset
```

**4. Form submission triggers**: For elements with role=form:

```
// Fill required fields with test data
ab('fill @required_field "test value"')
ab('click @submit_button')
ab('wait --load networkidle')
state_snapshot = abJson('snapshot -i')
// Extract new elements (validation results, success/error messages)
ab('press Escape')  // or navigate back to reset
```

For each trigger type:
1. Compare state snapshot with base snapshot, extract new elements
2. Record as `states` entry: `{ name, trigger: "<elementID>", elements: [...] }`
3. `trigger` references the trigger element's E-NNN ID (e.g. `"E-002"`)
4. In-state elements also receive E-NNN IDs
5. Reset to base state before exploring next trigger

> **Note**: `@eN` is agent-browser CLI's element reference syntax, used only during sitemap generation. Generated test scripts (`*.spec.ts`) must NOT use `@eN`; they must use Playwright Locator API.
