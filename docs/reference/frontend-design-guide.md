# Frontend design guide

The Flower frontend aims to feel like a focused planning board: calm surfaces, dense-but-readable type, and colour used to signal story type and state rather than decoration.

## Colour palette

Flower’s visual identity is petal magenta, stem green, and warm cream paper.

### Brand colours

- **Primary (Bloom):** `#C43B6E`
  - Usage: primary buttons, current-iteration highlights, key links.
  - Hover `#d24d7d`, active `#a3325a`.
- **Secondary (Stem):** `#2F7D4A`
  - Usage: accepted / done states, secondary buttons.
- **Accent (Pollen):** `#E8A317`
  - Usage: started / in-progress indicators, estimate chips.

### Neutral / utility

- **Text primary / “black”:** `#1C1917` (use wherever pure black would normally appear).
- **Surface base:** `#FFFFFF` for cards and the board; pair with `#E8DFD6` borders.
- **Paper background:** `#FBF7F2`.
- **Borders / dividers:** `#E8DFD6` over paper, `#D9D0C7` over white cards.

### Status colours

- **Accepted / success:** `#2F7D4A`
- **Rejected / destructive:** `#B42318`
- **Started / in progress:** `#E8A317`
- **Delivered:** `#3A6EA5`

## Typography

- **Headings:** Fraunces or a high-contrast serif. Semi-bold.
- **Body:** Inter. Regular (400) for story titles and panel copy.
- **Meta / labels:** Medium (500), often uppercase and muted.

## UI components

### Buttons

- **Primary:** `rounded-full bg-bloom px-4 py-2 text-sm font-medium text-white hover:bg-bloom/90`
- **Secondary:** `rounded-full border border-paper-300 px-4 py-2 text-sm font-medium text-ink hover:bg-paper-50`
- **Destructive:** `rounded-full bg-red-700 px-4 py-2 text-sm font-medium text-white hover:bg-red-800`
- **Disabled:** always add `disabled:cursor-not-allowed disabled:opacity-50`

### Cards / stories

Stories are the primary object on screen. They should look like tickets on paper, not enterprise table rows.

- **Structure:** `rounded-lg border border-paper-300 bg-white`
- **Hover:** `hover:border-ink/20`
- **Selected:** `border-bloom ring-2 ring-bloom/30`

### Panels

Icebox, backlog, current, and done are full-height columns.

- **Column:** `flex h-full min-w-[22rem] flex-col border-r border-paper-300 bg-paper`
- **Column header:** `border-b border-paper-300 px-4 py-3 text-xs font-medium uppercase tracking-wider text-ink-700`

### Form inputs

- **Text input:** `rounded-lg border border-paper-300 px-3 py-2 text-sm focus:border-bloom focus:outline-none focus:ring-1 focus:ring-bloom`
- **Labels:** `text-sm font-medium text-ink-800`

## Spacing and motion

- Card padding: `p-3` for stories, `p-6` for settings cards.
- Section gaps: `gap-6`.
- Animation duration: `200ms` to `300ms`, `ease-out`.

## Iconography

Use Lucide icons, stroke width 2px, matching nearby text colour.
