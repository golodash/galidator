# Welcome to GoLodash contributing guide

Thank you for investing your time in contributing to our project! :sparkles:. 

Read our [Code of Conduct](/CODE_OF_CONDUCT.md) to keep our community approachable and respectable.

In this guide you will get an overview of the contribution workflow from opening an issue, creating a PR, reviewing, and merging the PR.

Use the table of contents icon <img src="./images/table-of-contents.png" width="25" height="25" /> on the top left corner of this document to get to a specific section of this guide quickly.

## New contributor guide

To get an overview of the project, read the [README](/README.md).

## Contributor setup

Install requirements for contributing:
1. Please install python3 on your system. (It is needed for scripts to format your commit messages and handle [CHANGELOG.rst](/CHANGELOG.rst) file)
2. Install venv module for python (linux users only):
   - `sudo apt-get install python3-venv`
   - If you have python3.10, install it like: `sudo apt-get install python3.10-venv`
3. Run `install.py` script like:
   - `python3 ./.githooks/install.py`

### Issues

#### Create a new issue

1. **Bug**: If you spot a problem with a certain release or a problem, create a bug issue.
2. **Feature**: If you want to add a new functionality to the project, create a feature issue.

### Make changes

**Note**: This workflow is designed based on git flow but please read till the end before doing anything.
1. Pick or create a feature or bug issue.
2. Fork the repository.
3. Create a new branch.
4. Do your changes on that branch.
5. Debug and be sure about the changes.
6. Add documentation for new functions and variables and any other new things you provided.
7. Commit and push your changes. (see [commit message guidelines](#commit-message-guidelines))
8. Create a pull request and **mention** your issue number inside the pull request.

### Changelog

Mostly do not touch or change the file and let the script handle it.

### Commit message guidelines

There is a template for how to commit, this template is based on [conventional commit](https://www.conventionalcommits.org/en/v1.0.0/):

- **\<type>(\<scope>): \<subject>**

Samples:

```
docs(changelog): update changelog to beta.5
```

#### Type (Essential)

* docs: Documentation only changes
* feat: A new feature
* fix: A bug fix
* perf: A code change that improves performance
* refactor: A code change that neither fixes a bug nor adds a feature
* style: Changes that do not affect the meaning of the code (white-space, formatting, missing semi-colons, etc)
* test: Adding missing tests or correcting existing tests

#### Scope (Optional)

This section can be written in two formats:
1. (\<package>-\<function>)
2. (\<file,description...>)

**Note**: If you don't specify this part, remove parenthesis too.

#### Subject (Essential)

A brief description about what just happened.

## Versioning (Extra Section)

This is just a reminder for everyone to know what versioning system we are using.

Versioning in this project is based on semantic versioning:

v**Major**.**Minor**.**Patch**-**PreReleaseIdentifier**

Example:
- v1.4.0-beta.1

### Major Version

Signals backward-incompatible changes in a module’s public API. This release carries no guarantee that it will be backward compatible with preceding major versions.

### Minor Version

Signals backward-compatible changes to the module’s public API. This release guarantees backward compatibility and stability.

### Patch Version

Signals changes that don’t affect the module’s public API or its dependencies. This release guarantees backward compatibility and stability.

### Pre-release Version

Signals that this is a pre-release milestone, such as an alpha or beta. This release carries no stability guarantees.

### More information

For more information about versioning, read [this](https://go.dev/doc/modules/version-numbers).