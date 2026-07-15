# Database Rule

When changing persistence:
- Update models deliberately.
- Consider migration impact immediately.
- Keep database access behind service methods or dedicated query layers.
- Use `GetGorm()` safely and handle errors.
- Do not hide schema changes inside unrelated refactors.
- If SQLC is in use, keep generated code in sync with query changes.
