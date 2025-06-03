# Release Management Guide

This document describes the release process for the LogSpan project.

## Overview

LogSpan follows [Semantic Versioning](https://semver.org/) and uses automated release workflows to ensure consistent and reliable releases.

## Version Format

- **Release versions**: `vX.Y.Z` (e.g., `v1.0.0`, `v1.2.3`)
- **Pre-release versions**: `vX.Y.Z-TYPE.TIMESTAMP.SHA` (e.g., `v1.1.0-beta.20231027123456.abc1234`)

### Semantic Versioning Rules

- **MAJOR** (`X`): Incompatible API changes
- **MINOR** (`Y`): New functionality in a backwards compatible manner
- **PATCH** (`Z`): Backwards compatible bug fixes

## Release Types

### 1. Production Releases

Production releases are created from the `main` branch and follow semantic versioning.

#### Automatic Release (Recommended)

Releases are automatically triggered when code is pushed to the `main` branch:

```bash
# The release workflow will automatically:
# 1. Run all tests and quality checks
# 2. Determine the next version (patch increment by default)
# 3. Generate changelog
# 4. Create GitHub release
# 5. Update documentation
```

#### Manual Release

You can also trigger a manual release with a specific version type:

1. Go to GitHub Actions → Release workflow
2. Click "Run workflow"
3. Select version type: `patch`, `minor`, or `major`
4. Click "Run workflow"

#### Local Release Management

Use the version management script for local operations:

```bash
# Show current version and next version options
./scripts/version.sh current

# List all versions
./scripts/version.sh list

# Create a new version locally (for testing)
./scripts/version.sh create patch
./scripts/version.sh create minor
./scripts/version.sh create major
```

### 2. Pre-releases

Pre-releases are created from the `develop` branch for testing purposes.

#### Automatic Pre-release

Pre-releases are automatically triggered when code is pushed to the `develop` branch:

```bash
# Push to develop branch
git checkout develop
git push origin develop

# This will create a pre-release version like:
# v1.1.0-beta.20231027123456.abc1234
```

#### Manual Pre-release

You can also trigger a manual pre-release:

1. Go to GitHub Actions → Pre-release workflow
2. Click "Run workflow"
3. Select pre-release type: `alpha`, `beta`, or `rc`
4. Click "Run workflow"

## Release Process

### For Maintainers

#### 1. Prepare for Release

1. **Ensure all tests pass**:
   ```bash
   go test ./...
   ```

2. **Run quality checks**:
   ```bash
   golangci-lint run
   gosec ./...
   ```

3. **Update documentation** if needed

4. **Review and merge all pending PRs**

#### 2. Create Release

##### Option A: Automatic Release (Recommended)

1. Merge all changes to `main`:
   ```bash
   git checkout main
   git pull origin main
   git push origin main
   ```

2. The release workflow will automatically:
   - Run all tests and quality checks
   - Increment patch version
   - Generate changelog
   - Create GitHub release
   - Update CHANGELOG.md

##### Option B: Manual Release with Specific Version

1. Use GitHub Actions manual trigger:
   - Go to Actions → Release workflow
   - Click "Run workflow"
   - Select version type (`major`, `minor`, `patch`)
   - Click "Run workflow"

#### 3. Post-Release

1. **Verify the release**:
   - Check GitHub releases page
   - Verify pkg.go.dev updates
   - Test installation: `go get github.com/zentooo/logspan@vX.Y.Z`

2. **Announce the release** (if significant):
   - Update project documentation
   - Notify users through appropriate channels

### For Contributors

#### 1. Development Workflow

1. **Create feature branch**:
   ```bash
   git checkout -b feature/your-feature
   ```

2. **Develop and test**:
   ```bash
   # Make changes
   go test ./...
   golangci-lint run
   ```

3. **Create Pull Request** to `develop` branch

4. **Test with pre-release**:
   - After merge to `develop`, a pre-release will be created
   - Test the pre-release version

#### 2. Testing Pre-releases

```bash
# Install pre-release version
go get github.com/zentooo/logspan@v1.1.0-beta.20231027123456.abc1234

# Test your application with the pre-release
go test ./...
```

## Changelog Management

### Automatic Changelog Generation

The release workflow automatically generates changelog entries based on commit messages since the last release.

### Manual Changelog Updates

You can manually update `CHANGELOG.md` following the [Keep a Changelog](https://keepachangelog.com/) format:

```markdown
## [vX.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security improvements
```

## Commit Message Guidelines

To ensure good changelog generation, follow these commit message conventions:

```bash
# Feature additions
feat: add new logging middleware
feat(context): implement auto-flush functionality

# Bug fixes
fix: resolve memory leak in context logger
fix(formatter): handle nil context values

# Documentation
docs: update API documentation
docs(readme): add installation instructions

# Refactoring
refactor: simplify middleware pipeline
refactor(logger): optimize performance

# Tests
test: add benchmark tests
test(middleware): improve test coverage

# Chores
chore: update dependencies
chore(ci): improve workflow performance
```

## Troubleshooting

### Common Issues

#### 1. Release Workflow Fails

**Problem**: Release workflow fails during execution

**Solutions**:
- Check GitHub Actions logs for specific error
- Ensure all tests pass locally
- Verify no duplicate version tags exist
- Check repository permissions

#### 2. Version Already Exists

**Problem**: Trying to create a version that already exists

**Solutions**:
```bash
# Check existing versions
./scripts/version.sh list

# Delete problematic version (if safe)
./scripts/version.sh delete vX.Y.Z

# Or increment to next available version
```

#### 3. Changelog Not Updated

**Problem**: CHANGELOG.md not updated after release

**Solutions**:
- Check if release workflow completed successfully
- Manually update CHANGELOG.md if needed
- Ensure proper commit message format

### Emergency Procedures

#### Rollback a Release

If a release has critical issues:

1. **Create hotfix**:
   ```bash
   git checkout -b hotfix/critical-fix
   # Fix the issue
   git commit -m "fix: critical issue description"
   ```

2. **Create patch release**:
   ```bash
   # Merge hotfix to main
   git checkout main
   git merge hotfix/critical-fix
   git push origin main
   # This will trigger automatic patch release
   ```

3. **Mark problematic release** (if needed):
   - Edit GitHub release to mark as problematic
   - Add warning in release notes

## Security Considerations

### Security Releases

For security-related releases:

1. **Follow responsible disclosure**
2. **Create security advisory** on GitHub
3. **Use patch version** for security fixes
4. **Update security documentation**

### Access Control

- Only maintainers can trigger manual releases
- All releases require passing security scans
- Release artifacts are signed and verified

## References

- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [GitHub Releases](https://docs.github.com/en/repositories/releasing-projects-on-github)
- [Go Modules](https://golang.org/ref/mod)