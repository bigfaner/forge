# Sitemap Schema Reference

## Key Fields

| Field | Description |
|-------|-------------|
| `baseUrl` | Application base URL (e.g. `http://localhost:3456`) |
| `updatedAt` | Last update time (RFC3339 format) |
| `layout.name` | Layout component name (e.g. `AppLayout`) |
| `layout.wraps` | List of routes sharing this layout |
| `layout.elements[]` | Layout-level shared elements (sidebar, top nav, etc.), ID format `L-NNN` |
| `pages[].elements[].role` | Accessibility role (button, heading, etc.) |
| `pages[].elements[].name` | Accessibility name |
| `pages[].elements[].level` | Heading level (heading role only) |
| `pages[].elements[].label` | Associated label text (form elements like textbox only) |
| `pages[].elements[].placeholder` | Placeholder text (textbox only) |
| `pages[].states[]` | Dynamic states (modal, tab panel, dropdown, etc.) |
| `pages[].states[].trigger` | Trigger element ID (e.g. `"E-002"`) |
| `pages[].states[].elements` | Elements within the state (also with E-NNN IDs) |

## Layout Elements vs Page Elements

- `layout.elements`: Elements shared across all wrapped pages (navbar, sidebar, footer), ID format `L-NNN`
- `pages[].elements`: Elements unique to the current page, ID format `E-NNN`
- During test script generation, layout elements are available in any wrapped page

## Element ID Assignment Rules

- **Layout elements**: format `L-NNN`, globally unique, independent numbering space
- **Page elements**: format `E-NNN`, globally unique (base and states share ID space)
- First generation: start from `L-001` and `E-001` respectively
- Incremental update: new elements start from current max ID + 1
- IDs are never reused (deleted IDs are not recycled)

## Dynamic Route Parameterization

| URL segment pattern | Replace with | Example |
|---------------------|-------------|---------|
| Pure numeric | `:id` | `/tasks/42` -> `/tasks/:id` |
| UUID format | `:uuid` | `/orders/550e8400-...` -> `/orders/:uuid` |
| 32-char hex | `:hash` | `/files/a1b2c3d4e5f6...` -> `/files/:hash` |

During deduplication, routes with the same template keep only one entry. `layout.wraps` also uses template form.

Also merge manually added routes from existing sitemap that are not found by either source.

## Full Example

```json
{
	"baseUrl": "http://localhost:3456",
	"updatedAt": "2026-04-25T10:00:00Z",
	"layout": {
		"name": "AppLayout",
		"wraps": ["/", "/settings", "/tasks/:id"],
		"elements": [
			{ "id": "L-001", "role": "navigation", "name": "Main nav" },
			{ "id": "L-002", "role": "link", "name": "Dashboard" },
			{ "id": "L-003", "role": "link", "name": "Settings" },
			{ "id": "L-004", "role": "link", "name": "Profile" },
			{ "id": "L-005", "role": "button", "name": "Logout" }
		]
	},
	"pages": [
		{
			"route": "/login",
			"title": "Login",
			"elements": [
				{ "id": "E-030", "role": "heading", "name": "Sign In", "level": 1 },
				{ "id": "E-031", "role": "textbox", "name": "Username", "label": "Username", "placeholder": "Enter username" },
				{ "id": "E-032", "role": "textbox", "name": "Password", "label": "Password", "placeholder": "Enter password" },
				{ "id": "E-033", "role": "button", "name": "Login" }
			]
		},
		{
			"route": "/",
			"title": "Dashboard",
			"elements": [
				{ "id": "E-001", "role": "heading", "name": "Dashboard", "level": 1 },
				{ "id": "E-002", "role": "button", "name": "Add Task" },
				{ "id": "E-003", "role": "textbox", "label": "Search", "placeholder": "Search tasks..." }
			],
			"states": [
				{
					"name": "add-task-modal",
					"trigger": "E-002",
					"elements": [
						{ "id": "E-004", "role": "dialog", "name": "Add Task" },
						{ "id": "E-005", "role": "textbox", "label": "Task Title" },
						{ "id": "E-006", "role": "button", "name": "Save" },
						{ "id": "E-007", "role": "button", "name": "Cancel" }
					]
				}
			]
		},
		{
			"route": "/settings",
			"title": "Settings",
			"elements": [
				{ "id": "E-010", "role": "heading", "name": "Settings", "level": 1 },
				{ "id": "E-011", "role": "tab", "name": "General" },
				{ "id": "E-012", "role": "tab", "name": "Notifications" },
				{ "id": "E-013", "role": "button", "name": "Save" }
			]
		},
		{
			"route": "/tasks/:id",
			"title": "Task Detail",
			"elements": [
				{ "id": "E-020", "role": "heading", "name": "Task Detail", "level": 1 },
				{ "id": "E-021", "role": "button", "name": "Edit" },
				{ "id": "E-022", "role": "button", "name": "Delete" },
				{ "id": "E-023", "role": "link", "name": "Back" }
			],
			"states": [
				{
					"name": "delete-confirm",
					"trigger": "E-022",
					"elements": [
						{ "id": "E-024", "role": "dialog", "name": "Confirm Delete" },
						{ "id": "E-025", "role": "button", "name": "Confirm" },
						{ "id": "E-026", "role": "button", "name": "Cancel" }
					]
				}
			]
		}
	]
}
```
