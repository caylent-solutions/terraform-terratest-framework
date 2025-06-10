# CI/CD Pipeline

This document describes the CI/CD pipeline for the Terraform Terratest Framework.

## Pull Request Workflow

When a pull request is opened against the main branch, the following automated checks are performed:

1. **Code checkout**: The repository is checked out with full history
2. **Merge simulation**: A simulated merge is performed to ensure the PR can be merged cleanly
3. **Dependency installation**: System and project dependencies are installed
4. **Linting**: Code is checked for style and quality issues
5. **CLI build**: The tftest CLI tool is built
6. **Unit tests**: All unit tests are run
7. **Functional tests**: All functional tests are run
8. **Test coverage**: Code coverage is checked (minimum 40% required)
9. **Code owners identification**: Code owners are identified based on changed files
10. **Slack notification**: A notification is sent to Slack requesting code review

The PR workflow helps ensure that all code changes meet quality standards before being merged.

## Main Branch Validation Workflow

When code is pushed to the main branch (typically through a merged PR), the following workflow is triggered:

1. **Validation**: Similar to the PR workflow, the code is validated through:
   - System dependency installation
   - Project dependency installation
   - Linting
   - CLI building
   - Unit testing
   - Functional testing
   - Test coverage checking

2. **CodeQL Analysis**: Security scanning is performed using GitHub's CodeQL

3. **Manual QA Approval**: After automated checks pass, a manual approval is required from the QA team

4. **Release Process**: Once QA approval is granted, the release process begins:
   - Next version is computed using semantic-release
   - Changelog is generated
   - VERSION file is updated
   - Changes are committed to a release branch
   - A PR is created and automatically merged
   - A Git tag is created for the new version

## Release Strategy

The project follows semantic versioning with automated version determination:

- **Major version bump** (x.0.0): Breaking changes
  - Commits with `BREAKING CHANGE:` prefix
  - Commits with an exclamation mark (`!`): `feat!:`, `fix!:`, etc.

- **Minor version bump** (0.x.0): New features
  - Commits with `feat:` or `feature:` prefix

- **Patch version bump** (0.0.x): Bug fixes and other changes
  - Commits with `fix:`, `docs:`, `style:`, `refactor:`, `test:`, `chore:`, `ci:`, `build:`, `perf:` prefixes

## Manual Release Process

In exceptional circumstances when the automated pipeline cannot be used, a manual release can be performed using:

```bash
make release-manual [TYPE=major|minor|patch]
```

This should only be used when the automated pipeline is unavailable.