# End-to-end tests

This directory holds `footloose` end-to-end tests. All commands given in this
README are assuming being run from the root of the repository.

The prerequisites to run to the tests are:

- `docker` installed on the machine with no container running. This
limitation can be lifted once we can select `footloose` containers better
([#17][issue-17]).
- `footloose` in the path.

[issue-17]: https://github.com/dlespiau/footloose/issues/17

## Running the tests

To run all tests:

```console
go tests -v ./tests
```

To exclude long running tests (useful to smoke test a change before a longer
run in CI):

```console
go tests -short -v ./tests
```

To run a specific test:

```console
go test -v -run TestEndToEnd/test-create-delete-centos7
```

Remember that the `-run` argument is a regex so it's possible to select a
subset of the tests with this:

```console
go test -v -run TestEndToEnd/test-create-delete
```

This will match `test-create-delete-centos7`, `test-create-delete-fedora29`,
...
