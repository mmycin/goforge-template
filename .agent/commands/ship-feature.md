# Ship Feature

Use this command for end-to-end feature delivery.

Workflow:
1. Read `.agent/commands/context.md`.
2. Identify which skill applies under `.agent/skills/`.
3. Implement the smallest coherent change first.
4. Update wiring, config, migrations, and tests in the same pass.
5. Verify with targeted tests before broad test runs.

Feature checklist:
- Service logic
- Handler or transport updates
- Route or registration updates
- Database or config changes
- Tests
- Final change summary with risks
