# Contributing to LoginRadius CLI

:+1::tada: First off, thanks for taking the time to contribute! :tada::+1:

The following is a set of guidelines for contributing to LoginRadius CLI, which is hosted in the [LoginRadius Organization](https://github.com/loginradius) on GitHub. These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

## Code of Conduct

This project and everyone participating in it are governed by the [LoginRadius Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behaviour to [support@loginradius.com](mailto:support@loginradius.com).

## I don't want to read this whole thing I just have a question!!!

> **Note:** Please don't file an issue to ask a question. You'll get faster results by using the resources below.

* [Discuss on Community Page](https://community.loginradius.com/)
* [Ask On StackOverflow](https://stackoverflow.com/questions/ask/?tags=loginradius)

## How Can I Contribute?

### Reporting Bugs

* **Do not open up a GitHub issue if the bug is a security vulnerability
  in LoginRadius CLI**, and instead to refer to our [security policy](https://www.loginradius.com/security-policy).

* **Ensure the bug was not already reported** by searching on GitHub under [Issues](https://github.com/loginradius/lr-cli/issues).

* If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/loginradius/lr-cli/issues/new) and provide the following information by filling in [the template](https://github.com/loginradius/lr-cli/.github/blob/master/.github/ISSUE_TEMPLATE/bug_report.md).

* **If the problem wasn't triggered by a specific action**, describe what you were doing before the problem happened and share more information using the guidelines below.

#### Provide more context by answering these questions:

* **Did the problem start happening recently** (e.g. after updating to a new version of LoginRadius CLI) or was this always a problem?
* If the problem started happening recently, **can you reproduce the problem in an older version of LoginRadius CLI?** What's the most recent version in which the problem doesn't happen? You can download older versions of LoginRadius CLI from [the releases page](https://github.com/loginradius/lr-cli/releases).
* **Can you reliably reproduce the issue?** If not, provide details about how often the problem happens and under which conditions it normally happens.

#### Did you write a patch that fixes a bug?

* Open a new GitHub pull request with the patch.

* Ensure the PR description clearly describes the problem and solution. Include the relevant issue number if applicable.

#### Did you fix whitespace, format code, or make a purely cosmetic patch?

Changes that are cosmetic and do not add anything substantial to the stability, functionality, or testability of LoginRadius CLI will generally not be accepted.

### Suggesting Enhancements

* Suggest your change in the [LoginRadius community page](https://community.loginradius.com/) and start writing code.

* Do not open an issue on GitHub until you have collected positive feedback about the change. GitHub issues are primarily intended for bug reports and fixes.

### Your First Code Contribution

Unsure where to begin contributing to LoginRadius? You can start by looking through these `beginner` and `help-wanted` issues:

* [Beginner issues][beginner] - issues that should only require a few lines of code, and a test or two.
* [Help wanted issues][help-wanted] - issues which should be a bit more involved than `beginner` issues.

### Pull Requests

The process described here has several goals:

- Maintain LoginRadius's quality
- Fix problems that are important to users
- Engage the community in working toward the best possible authentication experience.
- Enable a sustainable system for LoginRadius's maintainers to review contributions

Please follow these steps to have your contribution considered by the maintainers:

1. Follow all instructions in [the template](PULL_REQUEST_TEMPLATE.md)
2. Follow the [styleguides](#styleguides)
3. After you submit your pull request, verify that all [status checks](https://help.github.com/articles/about-status-checks/) are passing <details><summary>What if the status checks are failing?</summary>If a status check is failing, and you believe that the failure is unrelated to your change, please leave a comment on the pull request explaining why you believe the failure is unrelated. A maintainer will re-run the status check for you. If we conclude that the failure was a false positive, then we will open an issue to track that problem with our status check suite.</details>

While the prerequisites above must be satisfied before having your pull request reviewed, the reviewer(s) may ask you to complete additional design work, tests, or other changes before your pull request can be ultimately accepted.

## Styleguides

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line
* When only changing documentation, include `[ci skip]` in the commit title
* Consider starting the commit message with an applicable emoji:
    * :art: `:art:` when improving the format/structure of the code
    * :racehorse: `:racehorse:` when improving performance
    * :non-potable_water: `:non-potable_water:` when plugging memory leaks
    * :memo: `:memo:` when writing docs
    * :penguin: `:penguin:` when fixing something on Linux
    * :apple: `:apple:` when fixing something on macOS
    * :checkered_flag: `:checkered_flag:` when fixing something on Windows
    * :bug: `:bug:` when fixing a bug
    * :fire: `:fire:` when removing code or files
    * :green_heart: `:green_heart:` when fixing the CI build
    * :white_check_mark: `:white_check_mark:` when adding tests
    * :lock: `:lock:` when dealing with security
    * :arrow_up: `:arrow_up:` when upgrading dependencies
    * :arrow_down: `:arrow_down:` when downgrading dependencies
    * :shirt: `:shirt:` when removing linter warnings

Thanks! :heart: :heart: :heart:

LoginRadius Team





