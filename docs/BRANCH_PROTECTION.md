# Branch Protection Configuration Guide

## Main Branch Protection (main)
- Branch name pattern: `main`
- Protection rules:
  - ✅ Require a pull request before merging
  - ✅ Require approvals: 1
  - ✅ Dismiss stale PR approvals when new commits are pushed
  - ✅ Require review from code owners
  - ✅ Require status checks to pass before merging
    - test
    - lint
    - build (ubuntu-latest)
    - build (windows-latest)  
    - build (macos-latest)
  - ✅ Require branches to be up to date before merging
  - ✅ Require conversation resolution before merging
  - ❌ Allow deletions (protect main branch)
  - ❌ Allow force pushes (protect main branch)
  - Bypass restrictions:
    - @Ledermayer
    - @LederWorks/hippodamus-maintainers

## Release Branch Protection (release/*)
- Branch name pattern: `release/*`
- Protection rules:
  - ✅ Require a pull request before merging
  - ✅ Require approvals: 1
  - ✅ Require review from code owners
  - ✅ Allow deletions (teams can clean up)
  - ✅ Allow force pushes (for hotfixes)
  - Bypass restrictions:
    - @Ledermayer
    - @LederWorks/hippodamus-maintainers
  - Who can delete:
    - @LederWorks/hippodamus-maintainers
    - @LederWorks/hippodamus-devops-team

## Feature Branch Protection (feature/*)
- Branch name pattern: `feature/*`
- Protection rules:
  - ✅ Allow deletions (clean up after merge)
  - ✅ Allow force pushes (for rebasing)
  - ❌ No other restrictions (development branches)
  - Who can delete:
    - Branch creators
    - @LederWorks/hippodamus-maintainers

## Hotfix Branch Protection (hotfix/*)
- Branch name pattern: `hotfix/*`
- Protection rules:
  - ✅ Require a pull request before merging
  - ✅ Require approvals: 1
  - ✅ Allow deletions (clean up after merge)
  - ✅ Allow force pushes (for urgent fixes)
  - Who can delete:
    - @LederWorks/hippodamus-maintainers
    - @LederWorks/hippodamus-devops-team

## Team Permissions Summary

### hippodamus-maintainers
- Role: Admin
- Can: Delete any branch, bypass all protections
- Members: Core maintainers

### hippodamus-devops-team  
- Role: Maintain
- Can: Delete release/hotfix branches, manage CI/CD
- Members: DevOps engineers

### hippodamus-core-team
- Role: Write
- Can: Create PRs, review code
- Cannot: Delete protected branches

## Manual Configuration Steps

1. Go to GitHub.com → LederWorks/hippodamus → Settings → Branches
2. Create each branch protection rule above
3. Go to Settings → Collaborators and teams
4. Ensure teams have correct role assignments
5. Test branch deletion permissions
