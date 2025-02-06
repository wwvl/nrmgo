# Contributing Guide

Thank you for your interest in the NRMGO project! We welcome all forms of contributions, including but not limited to: feature improvements, bug fixes, documentation enhancements, etc.

## Code of Conduct

This project follows the [Contributor Covenant](https://www.contributor-covenant.org/version/2/0/code_of_conduct/). By participating in this project, you agree to abide by its terms.

## Getting Started

1. Fork this repository
2. Clone your fork
   ```bash
   git clone https://github.com/your-wwvl/nrmgo.git
   ```
3. Add upstream repository
   ```bash
   git remote add upstream https://github.com/wwvl/nrmgo.git
   ```
4. Create development branch
   ```bash
   git checkout -b feature/your-feature
   ```

## Development Process

1. Ensure your code meets our coding standards

   ```bash
   make lint
   ```

2. Pre-commit checks

   ```bash
   # Run pre-release checks
   ./scripts/pre-release.ps1 -fix
   ```

3. Commit code

   ```bash
   git add .
   git commit -m "feat: your feature description"
   git push origin feature/your-feature
   ```

4. Create Pull Request

## Commit Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code formatting (changes that do not affect code execution)
- `refactor`: Code changes that neither fixes a bug nor adds a feature
- `perf`: Performance improvements
- `test`: Adding or modifying tests
- `chore`: Changes to build process or auxiliary tools

Examples:

```bash
feat: add automatic retry mechanism
fix: resolve configuration file path issue
docs: update installation instructions
style: format code
refactor: restructure error handling logic
perf: optimize caching mechanism
test: add concurrency tests
chore: update build scripts
```

## Branch Management

- `main`: Main branch for releases
- `develop`: Development branch
- `feature/*`: Feature branches
- `fix/*`: Bug fix branches
- `release/*`: Release branches

## Release Process

1. Version Number Convention

   - Major version: Incompatible API changes
   - Minor version: Backward-compatible functionality additions
   - Patch version: Backward-compatible bug fixes

2. Release Steps

   ```bash
   # Prepare release
   ./scripts/pre-release.ps1

   # Create release
   ./scripts/release.ps1

   # Verify release
   ./scripts/verify-release.ps1
   ```

## Development Guidelines

1. Code Quality

   - Keep code simple and clear
   - Add necessary comments
   - Follow Go coding standards
   - Use meaningful variable and function names

2. Documentation

   - Update relevant documentation
   - Add code examples
   - Document important changes

3. Performance
   - Mind resource usage
   - Avoid unnecessary memory allocations
   - Consider concurrency safety

## Feedback

- Submit issues via [GitHub Issues](https://github.com/wwvl/nrmgo/issues)
- Discuss features in [GitHub Discussions](https://github.com/wwvl/nrmgo/discussions)
