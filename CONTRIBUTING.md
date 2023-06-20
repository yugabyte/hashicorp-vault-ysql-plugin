# Contributing to Vault

**First:** if you're unsure or afraid of _anything_, just ask or submit the
issue or pull request anyways. You won't be yelled at for giving it your best
effort. The worst that can happen is that you'll be politely asked to change
something. We appreciate any sort of contributions, and don't want a wall of
rules to get in the way of that. 

## Issues

This section will cover what we're looking for in terms of reporting issues.

By addressing all the points we're looking for, it raises the chances we can
quickly merge or address your contributions.

### Reporting an Issue

* Make sure you test against the latest released version. It is possible we
  already fixed the bug you're experiencing. Even better is if you can test
  against the `master` branch, as the bugs are regularly fixed but new versions
  are only released every few months.

* Provide steps to reproduce the issue, and if possible include the expected 
  results as well as the actual results. Please provide text, not screen shots!

* If you experienced a panic, please create a [gist](https://gist.github.com)
  of the *entire* generated crash log for us to look at. Double check
  no sensitive items were in the log.

* Respond as promptly as possible to any questions made by the YugabyteDB
  team to your issue.

## Pull requests

When submitting a PR you should reference an existing issue. If no issue already exists, 
please create one. This can be skipped for trivial PRs like fixing typos.

Creating an issue in advance of working on the PR can help to avoid duplication of effort, 
e.g. maybe we know of existing related work. Or it may be that we can provide guidance 
that will help with your approach.

Your pull request should have a description of what it accomplishes, how it does so,
and why you chose the approach you did.  PRs should include unit tests that validate
correctness and the existing tests must pass.  Follow-up work to fix tests
does not need a fresh issue filed.

Someone will do a first pass review on your PR making sure it follows the guidelines 
in this document.  If it doesn't we'll mark the PR incomplete and ask you to follow
up on the missing requirements.