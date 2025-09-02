# Contributing to DebussyOps

Thank you for your interest in contributing to **DebussyOps**!  
We welcome contributions that improve the project’s functionality, documentation, and overall quality.  
This document outlines the guidelines to ensure a smooth and consistent contribution process.

---

## Getting Started

1. **Fork the repository**
   - Click the **Fork** button on the top right of this repo.
   - Clone your fork locally:
     ```bash
     git clone https://github.com/<your-username>/DebussyOps.git
     cd DebussyOps
     ```

2. **Set up the upstream remote**
   ```bash
   git remote add upstream https://github.com/rcallaby/DebussyOps.git
    ```

3. **Create a feature branch**

   ```bash
   git checkout -b feature/my-new-feature
   ```

---

## Development Guidelines

* **Code style**

  * Follow standard **Go conventions** (since this is a Go project).
  * Keep code modular and well-documented using comments where necessary.
  * Run `go fmt ./...` before committing.

* **Commits**

  * Write clear, descriptive commit messages.
  * Use the following convention where possible:

    ```
    <type>: <short description>
    ```

    Examples:

    * `fix: resolve panic in scheduler`
    * `feat: add support for YAML config`
    * `docs: improve usage section in README`

* **Testing**

  * All new features should include unit tests where applicable.
  * Run tests before submitting a PR:

    ```bash
    go test ./...
    ```

* **Documentation**

  * Update `README.md` or inline documentation if your changes affect usage.
  * Keep examples minimal but clear.

---

## Issues and Pull Requests

* **Before opening a PR**

  * Ensure your branch is up-to-date with `main`:
  * Verify all tests pass.

* **Issues**

  * Use issues to report bugs, request features, or propose improvements.
  * Provide as much context as possible (steps to reproduce, environment, expected behavior).

* **Pull Requests**

  * Keep PRs focused and atomic (one feature or fix per PR).
  * Reference related issues (e.g., `Closes #12`).
  * Include a brief description of what the change does and why it’s needed.

---

## Contribution Checklist

Before submitting, please make sure you:

* [ ] Followed the code style (`go fmt ./...`)
* [ ] Added or updated tests if applicable
* [ ] Updated relevant documentation
* [ ] Verified that all tests pass (`go test ./...`)
* [ ] Squashed commits into logical units if necessary

---

## Questions?

If you have questions start a relevant issue.
We value clarity, collaboration, and clean code — let’s keep DebussyOps growing together!



