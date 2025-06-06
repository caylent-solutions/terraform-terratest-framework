# Contributing to Terraform Terratest Framework

Thank you for your interest in contributing to the Terraform Terratest Framework! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### For External Contributors

If you're an external contributor (not a Caylent employee), please follow the standard open source fork and pull request workflow:

1. **Fork the Repository**:
   - Fork the repository to your GitHub account.
   - Clone your fork locally: `git clone https://github.com/YOUR-USERNAME/terraform-terratest-framework.git`

2. **Create a Branch**:
   - Create a branch for your changes: `git checkout -b feature/your-feature-name`

3. **Make Your Changes**:
   - Make your changes following the coding standards and guidelines.
   - Write tests for your changes.
   - Ensure all tests pass: `make test`

4. **Commit Your Changes**:
   - Use conventional commit messages:
     - `feat:` for new features
     - `fix:` for bug fixes
     - `docs:` for documentation changes
     - `test:` for test changes
     - `refactor:` for code refactoring
     - `chore:` for routine tasks, maintenance, etc.
   - Example: `feat: add support for AWS S3 bucket testing`

5. **Push Your Changes**:
   - Push your changes to your fork: `git push origin feature/your-feature-name`

6. **Submit a Pull Request**:
   - Go to the original repository and create a pull request from your branch.
   - Provide a clear description of your changes.
   - Reference any related issues.

7. **Review Process**:
   - Maintainers will review your PR and may request changes.
   - Once approved, your PR will be merged.

### For Caylent Contributors

If you're a Caylent employee, please follow the trunk-based development workflow:

1. **Clone the Repository**:
   - Clone the repository directly: `git clone https://github.com/caylent-solutions/terraform-terratest-framework.git`

2. **Create a Branch**:
   - Create a short-lived feature branch: `git checkout -b feature/your-feature-name`

3. **Make Your Changes**:
   - Make your changes following the coding standards and guidelines.
   - Write tests for your changes.
   - Ensure all tests pass: `make test`

4. **Commit Your Changes**:
   - Use conventional commit messages as described above.

5. **Push Your Changes**:
   - Push your changes to the repository: `git push origin feature/your-feature-name`

6. **Create a Pull Request**:
   - Create a PR to the main branch.
   - Get it reviewed by at least one team member.
   - Address any feedback.

7. **Merge Your PR**:
   - Once approved, merge your PR to the main branch.
   - Delete your feature branch after merging.

## Release Process

Releases are managed using semantic versioning and automated with the `make release` command:

```bash
# Create a release based on commit messages
make release

# Explicitly specify the version bump type
make release TYPE=major
make release TYPE=minor
make release TYPE=patch
```

The release process:
1. Determines the version bump type based on commit messages (or uses the specified type)
2. Updates the VERSION file
3. Updates version references in other files
4. Commits the changes
5. Creates a git tag
6. Provides instructions for pushing the changes and tag

## Development Guidelines

### Code Style

- Follow Go best practices and idiomatic Go.
- Use `make format` to format your code before committing.
- Use `make lint` to check for linting issues.

### Testing

- Write tests for all new features and bug fixes.
- Ensure all tests pass before submitting a PR: `make test`
- Check test coverage: `make test-coverage`

### Documentation

- Update documentation for any changes to functionality.
- Document new features, options, and behaviors.
- Keep the README and other documentation up to date.

## Getting Help

If you have questions or need help, please:
- Open an issue in the repository
- Contact the maintainers
- Refer to the documentation

Thank you for contributing to the Terraform Terratest Framework!