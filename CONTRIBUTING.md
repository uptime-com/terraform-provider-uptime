# Contributing Guidelines

## How Do I Submit A (Good) Bug Report or Feature Request

Please open a [GitHub issue](../../issues/new/choose) to report bugs or suggest features.

Please accurately fill out the appropriate GitHub issue form.

When filing an issue or feature request, help us avoid duplication and redundant effort. Check existing open or recently
closed issues first.

Detailed bug reports and requests are easier for us to work with. Please include the following in your issue:

* A reproducible test case or series of steps
* The version of provider being used
* Any modifications you've made, relevant to the bug
* Anything unusual about your environment or deployment
* Screenshots and code samples where illustrative and helpful

## How Do I Submit A (Good) Pull Request

We follow the [fork and pull model](https://opensource.guide/how-to-contribute/#opening-a-pull-request) for open source
contributions.

Tips for a faster merge:

* address one feature or bug per pull request
* make sure your commits are atomic, [addressing one change per commit](https://chris.beams.io/posts/git-commit/).
* provide tests for your changes
* make sure existing tests pass

## How do I run acceptance tests locally?

== Acceptance testing

In order to run acceptance tests Uptime.com API token is required.

    export TF_ACC=1
    export UPTIME_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    go test -v ./uptime/... -run ^TestAcc.+$

In order to get detailed HTTP requests/responses info, set `UPTIME_TRACE` environment variable to `1`.

    export UPTIME_TRACE=1

You can run those in Docker as well.

    cat > .env <<EOF
    TF_ACC=1
    UPTIME_TOKEN=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    EOF
    docker compose up testacc

## Licensing

See the [LICENSE file](/LICENSE) for our project's licensing.
