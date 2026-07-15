# Testing Rule

Every behavior change should come with verification.

Expectations:
- Run targeted tests for touched packages first.
- Run broader tests when startup, transport, or shared wiring changes.
- Prefer tests that exercise the public behavior, not just internals.
- Mention unrun tests explicitly in the final handoff.
- If a change is risky and tests are missing, say so plainly.
