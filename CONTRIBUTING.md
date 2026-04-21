# Contributing to MediFlow

Thank you for contributing to MediFlow.

## Contribution Principles

- Keep changes focused and easy to review
- Preserve implementation quality, readability, and operational clarity
- Update documentation when behavior or setup changes
- Add or adjust tests when code paths change

## Local Development

1. Copy `.env.example` to `.env`
2. Start the stack with `docker compose up --build`
3. Make changes in `frontend/` or `backend/`
4. Run relevant validation before opening a pull request

## Recommended Validation

Frontend:

```bash
cd frontend
npm install
npm run lint
npm run build
```

Backend:

```bash
cd backend
go test ./...
```

## Pull Request Expectations

- Use a clear title and describe the user or system impact
- Link related issues when applicable
- Include screenshots for visible UI changes
- Mention environment or migration changes explicitly
- Avoid bundling unrelated refactors into the same PR

## Commit Quality

- Prefer intentional, descriptive commit messages
- Keep commits logically grouped
- Do not commit secrets, credentials, or local-only configuration

## Documentation

If your change affects setup, architecture, APIs, or workflows, update the relevant documentation in the same pull request.
