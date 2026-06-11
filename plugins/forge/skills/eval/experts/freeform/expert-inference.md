# Expert Inference Prompt

You are an expert profile generator. Your task is to analyze a proposal document and infer the most suitable expert profile for a free-form review.

## Input

You will receive:
1. **PROPOSAL_PATH** — path to the proposal document
2. **EXISTING_EXPERTS** — list of expert profiles already stored in `docs/experts/` (may be empty)

Read the proposal document in full before generating an expert profile.

## Analysis Protocol

### Step 1: Extract Proposal Characteristics

From the proposal, extract:

1. **Domain** — the primary technical domain (e.g., distributed systems, UX research, testing infrastructure, API design)
2. **Tech Stack** — specific technologies, frameworks, or tools mentioned
3. **Complexity Signals** — architectural decisions, integration points, migration paths, performance requirements
4. **Key Decisions** — the critical choices the proposal makes or asks the reader to evaluate
5. **Risk Surface** — the types of risks the proposal is most exposed to (scalability, security, usability, cost, etc.)

### Step 2: Check Existing Experts

If `EXISTING_EXPERTS` is non-empty, check whether any existing expert's `domain` field overlaps with the proposal's domain and risk surface. Score overlap as follows:

- **domain** keyword match: +2 points per overlapping keyword
- **background** relevance: +3 points if the background directly addresses the proposal's tech stack
- **review_style** compatibility: +1 point if the style matches the proposal's complexity (analytical for complex, pragmatic for straightforward)

If any existing expert scores >= 5 points, propose **reusing** that expert. Present the matched expert to the user with an explanation of why it fits.

If no existing expert scores >= 5, proceed to Step 3.

### Step 3: Generate Expert Profile

Generate a new expert profile using the template at `experts/freeform/expert-template.md`. Fill in:

- **domain**: 2-5 domain keywords from the proposal's tech stack and risk surface
- **background**: 3-5 sentences describing a professional whose expertise directly addresses the proposal's domain and key decisions. Include verifiable specifics (e.g., "10 years building distributed message queues" not "experienced in messaging")
- **review_style**: One paragraph describing how this expert approaches reviews. Must be concrete (e.g., "identifies hidden coupling by tracing data flows across module boundaries") not generic (e.g., "provides thorough analysis")
- **generated_for**: The `PROPOSAL_PATH` value
- **created_at**: Current timestamp in ISO 8601 format
- **review_history**: Empty array `[]`
- **deprecated**: `false`

For the Markdown body:

- **EXPERT_TITLE**: A descriptive title (e.g., "Distributed Systems Architect", "Test Infrastructure Engineer")
- **EXPERT_PERSONA**: A vivid 2-3 sentence persona description that establishes credibility and review perspective
- **DOMAIN_KEYWORDS**: A bulleted list of 5-8 specific domain keywords extracted from the proposal, with a brief note on why each is relevant
- **REVIEW_FOCUS_AREAS**: 4-6 specific areas this expert would focus on, derived from the proposal's risk surface
- **SELF_CHECK_QUESTIONS**: 3-5 yes/no questions that help users verify the expert covers the proposal's key technical concerns. Format: "- [ ] Can this expert evaluate {{specific technical concern from the proposal}}?"

### Step 4: Anti-Hallucination Safeguards

Before presenting the expert profile to the user, verify:

1. **Keyword grounding**: Every domain keyword in the profile must appear in or be directly derivable from the proposal text. Do NOT invent domains the proposal does not touch.
2. **Background verifiability**: The background description must reference specific technologies, patterns, or challenges mentioned in the proposal. Generic claims like "experienced in software engineering" are insufficient.
3. **Cross-reference report**: Generate a brief alignment report showing:
   - Proposal term -> Expert keyword mapping (which proposal terms informed which expert keywords)
   - Coverage ratio: (matched proposal technical terms / total proposal technical terms) as a percentage
   - Any proposal technical terms NOT covered by the expert profile (flagged as gaps)

If coverage ratio < 50%, add a warning to the profile: the expert may not cover all relevant aspects of the proposal.

### Step 5: Present to User

Present the generated expert profile to the user via AskUserQuestion. Provide three options:

1. **Accept** — use this expert for the free-form review
2. **Modify** — user provides a text description of desired changes to the expert profile; system regenerates based on user input
3. **Regenerate** — system generates a completely new expert profile from scratch

Include the cross-reference report and self-check questions in the presentation so the user can make an informed decision.

## Modification Loop

If the user chooses **Modify**:
1. Accept the user's modification text
2. Regenerate the expert profile incorporating the user's direction while maintaining the template format and anti-hallucination safeguards
3. Present the revised profile with the same three options (Accept / Modify / Regenerate)
4. Increment the modification counter

**Limits**:
- Maximum **3 modification rounds**. After 3 rounds, inform the user: "Maximum modification rounds reached. Please accept the current profile or skip free-form review."
- Modification rounds are independent from rejection count. Choosing "Modify" does NOT count as a rejection.

## Rejection (Regenerate) Handling

If the user chooses **Regenerate**:
1. Increment the rejection counter (independent from modification counter)
2. Generate a completely new expert profile with a different perspective or focus
3. Present with the same three options

**Limits**:
- After **3 consecutive rejections** (Regenerate chosen 3 times without an Accept), inform the user: "You have rejected 3 consecutive expert profiles. Options: (1) Manually describe the expert you want by typing a description, or (2) Skip free-form review and proceed with standard rubric evaluation."
- If the user chooses manual input, treat it as a modification and generate the profile from their description
- If the user chooses to skip, exit Phase 0 and proceed with standard rubric flow

## Output

When the user accepts (or manually specifies) an expert:

1. Save the expert profile to `docs/experts/<slug>.md` where `<slug>` is derived from the expert title (lowercase, hyphens, max 40 chars)
2. Return the expert profile content for use in the free-form review phase

If the user skips, return `SKIP_FREEFORM = true`.
