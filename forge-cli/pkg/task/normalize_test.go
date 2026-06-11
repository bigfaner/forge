package task

import (
	"testing"
)

func TestNormalizeTaskMD(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "empty Hard Rules section between other sections",
			input: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Description
Some desc

## Hard Rules

## Implementation Notes
Some notes
`,
			want: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Description
Some desc

## Implementation Notes
Some notes
`,
		},
		{
			name: "Hard Rules with content is preserved",
			input: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Hard Rules
- MUST use just test-e2e
- MUST NOT use npx playwright test

## Implementation Notes
Some notes
`,
			want: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Hard Rules
- MUST use just test-e2e
- MUST NOT use npx playwright test

## Implementation Notes
Some notes
`,
		},
		{
			name: "no Hard Rules section unchanged",
			input: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Description
Some desc

## Implementation Notes
Some notes
`,
			want: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Description
Some desc

## Implementation Notes
Some notes
`,
		},
		{
			name: "empty Hard Rules at end of file",
			input: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Description
Some desc

## Hard Rules
`,
			want: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Description
Some desc
`,
		},
		{
			name: "empty Hard Rules with blank lines before next heading",
			input: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Hard Rules

## Description
Some desc
`,
			want: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Description
Some desc
`,
		},
		{
			name: "no frontmatter unchanged",
			input: `# Plain markdown

## Hard Rules

## Notes
`,
			want: `# Plain markdown

## Hard Rules

## Notes
`,
		},
		{
			name: "idempotent",
			input: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Hard Rules

## Notes
`,
			want: `---
id: "1.1"
title: "Test"
---
# 1.1: Test

## Notes
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeTaskMD([]byte(tt.input))

			// Idempotency check: running twice should produce same result
			got2 := NormalizeTaskMD(got)
			if string(got2) != string(got) {
				t.Errorf("NormalizeTaskMD not idempotent:\nfirst:  %q\nsecond: %q", got, got2)
			}

			if string(got) != tt.want {
				t.Errorf("NormalizeTaskMD()\ngot:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}
