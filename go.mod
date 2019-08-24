module github.com/tanema/promptui

require (
	github.com/alecthomas/gometalinter v3.0.0+incompatible
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/client9/misspell v0.3.4
	github.com/gordonklaus/ineffassign v0.0.0-20180909121442-1003c8bd00dc
	github.com/juju/ansiterm v0.0.0-20180109212912-720a0952cc2a
	github.com/manifoldco/promptui v0.3.2
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db
	github.com/stretchr/testify v1.2.2
	github.com/tsenart/deadcode v0.0.0-20160724212837-210d2dc333e9
	golang.org/x/lint v0.0.0-20181026193005-c67002cb31c3
)

// This version of kingpin is incompatible with the released version of
// gometalinter until the next release of gometalinter, and possibly until it
// has go module support, we'll need this exclude, and perhaps some more.
//
// After that point, we should be able to remove it.
exclude gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c
