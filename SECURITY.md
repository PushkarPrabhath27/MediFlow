# Security Policy

## Supported Scope

Security reports are welcomed for issues affecting:

- authentication and authorization
- tenant isolation
- sensitive configuration handling
- API access control
- infrastructure and dependency exposure

## Reporting a Vulnerability

Please report security issues privately. Do not open a public GitHub issue for a suspected vulnerability.

Include:

- a clear description of the issue
- reproduction steps or proof of concept
- affected components or files
- impact assessment if known

## Response Process

Reported issues will be triaged, validated, and addressed as quickly as possible. Coordinated disclosure is preferred.

## Secrets and Configuration

- Never commit production secrets
- Use `.env.example` only as a template
- Rotate any credential immediately if it is exposed
