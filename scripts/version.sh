#!/bin/bash

# Version management script for LogSpan
# Usage: ./scripts/version.sh [command] [options]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Get current version from git tags
get_current_version() {
    git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"
}

# Parse version string
parse_version() {
    local version="$1"
    # Remove 'v' prefix if present
    version=${version#v}

    IFS='.' read -ra VERSION_PARTS <<< "$version"
    MAJOR=${VERSION_PARTS[0]:-0}
    MINOR=${VERSION_PARTS[1]:-0}
    PATCH=${VERSION_PARTS[2]:-0}

    echo "$MAJOR $MINOR $PATCH"
}

# Calculate next version
next_version() {
    local current_version="$1"
    local version_type="$2"

    read -r MAJOR MINOR PATCH <<< "$(parse_version "$current_version")"

    case $version_type in
        major)
            MAJOR=$((MAJOR + 1))
            MINOR=0
            PATCH=0
            ;;
        minor)
            MINOR=$((MINOR + 1))
            PATCH=0
            ;;
        patch)
            PATCH=$((PATCH + 1))
            ;;
        *)
            log_error "Invalid version type: $version_type. Use: major, minor, or patch"
            exit 1
            ;;
    esac

    echo "v${MAJOR}.${MINOR}.${PATCH}"
}

# Show current version
show_current() {
    local current=$(get_current_version)
    log_info "Current version: $current"

    # Show what the next versions would be
    echo ""
    log_info "Next versions would be:"
    echo "  Patch: $(next_version "$current" "patch")"
    echo "  Minor: $(next_version "$current" "minor")"
    echo "  Major: $(next_version "$current" "major")"
}

# Validate version format
validate_version() {
    local version="$1"
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
        log_error "Invalid version format: $version. Expected format: vX.Y.Z"
        return 1
    fi
    return 0
}

# Check if version already exists
version_exists() {
    local version="$1"
    git tag -l | grep -q "^${version}$"
}

# Create a new version tag
create_version() {
    local version_type="$1"
    local current=$(get_current_version)
    local new_version=$(next_version "$current" "$version_type")

    log_info "Current version: $current"
    log_info "New version: $new_version"

    # Check if version already exists
    if version_exists "$new_version"; then
        log_error "Version $new_version already exists!"
        exit 1
    fi

    # Confirm with user
    echo ""
    read -p "Create version $new_version? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_warning "Version creation cancelled"
        exit 0
    fi

    # Create and push tag
    log_info "Creating tag $new_version..."
    git tag -a "$new_version" -m "Release $new_version"

    log_info "Pushing tag to origin..."
    git push origin "$new_version"

    log_success "Version $new_version created and pushed successfully!"
}

# List all versions
list_versions() {
    log_info "All versions:"
    git tag -l | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V || log_warning "No versions found"

    echo ""
    log_info "Pre-release versions:"
    git tag -l | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+-' | sort -V || log_warning "No pre-release versions found"
}

# Delete a version (with confirmation)
delete_version() {
    local version="$1"

    if [[ -z "$version" ]]; then
        log_error "Version not specified"
        exit 1
    fi

    if ! validate_version "$version"; then
        exit 1
    fi

    if ! version_exists "$version"; then
        log_error "Version $version does not exist"
        exit 1
    fi

    # Confirm with user
    echo ""
    log_warning "This will delete version $version locally and remotely!"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        log_warning "Version deletion cancelled"
        exit 0
    fi

    # Delete local tag
    log_info "Deleting local tag $version..."
    git tag -d "$version"

    # Delete remote tag
    log_info "Deleting remote tag $version..."
    git push origin ":refs/tags/$version"

    log_success "Version $version deleted successfully!"
}

# Show help
show_help() {
    cat << EOF
Version Management Script for LogSpan

Usage: $0 [command] [options]

Commands:
    current                 Show current version and next version options
    list                   List all versions (releases and pre-releases)
    create <type>          Create a new version tag
                          Types: major, minor, patch
    delete <version>       Delete a version tag (local and remote)
    help                   Show this help message

Examples:
    $0 current             # Show current version
    $0 list                # List all versions
    $0 create patch        # Create next patch version
    $0 create minor        # Create next minor version
    $0 create major        # Create next major version
    $0 delete v1.0.0       # Delete version v1.0.0

Notes:
    - This script works with git tags
    - Version format: vX.Y.Z (semantic versioning)
    - Tags are pushed to origin automatically
    - Use with caution in production environments
EOF
}

# Main script logic
main() {
    cd "$PROJECT_ROOT"

    # Check if we're in a git repository
    if ! git rev-parse --git-dir > /dev/null 2>&1; then
        log_error "Not in a git repository"
        exit 1
    fi

    case "${1:-help}" in
        current)
            show_current
            ;;
        list)
            list_versions
            ;;
        create)
            if [[ -z "$2" ]]; then
                log_error "Version type not specified. Use: major, minor, or patch"
                exit 1
            fi
            create_version "$2"
            ;;
        delete)
            if [[ -z "$2" ]]; then
                log_error "Version not specified"
                exit 1
            fi
            delete_version "$2"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown command: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"