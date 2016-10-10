# goworkon
A humble attempt at getting the functionality of virtualenv gor go.

## Requirements
A working GO version >= 1.4

## Usage
Creating environments:
``
goworkon create anenv --go-version 1.7 
``

This will create:

* A GOPATH for *anenv*
* A config for *anenv*
* If there is no 1.7 in $HOME/.local/share/go/ it will be checked out and built.

Switching to an environment:
``
goworkon switch anenv
``

This will:

* export PATH as $PATH:/this/go/version/bin/:$GOPATH/bin
* export GOPATH as $HOME/.local/share/goworkon/gopath

Config for this environment:
``
gowrokon add ./thisfolder github.com/foo/bar
``

Will set *thisfolder* as *github.com/foo/bar* in the environment

``
goworkon remove tighub.com/foo/bar
``

Will undo the previous command.

Updating a Go version:

``
goworkon update --go-version 1.7 --update-envs
``

Will update the go version in use to 1.7.latest and then rebuild envs using 1.7

Setting build steps:

``
goworkon build-steps wathever you please here semicolon separated
``

Will set the steps there, a special case

`` 
goworkon build-steps _1;_2;wathever you please semicolon separated;_3
``

Will replace ``_#`` for that step number, the same will apply to ``remove-steps``

``
goworkon build-steps --clear
``

Will delete the steps

``
goworkon build-steps
``

Will print the steps properly numbered like so:

``
_1: do something
_2: do something else
``

## TODO:
* Add a base config for the software that will
 * Define a default environment.
* Write switch, that changes the env variables:
 * PATH: backs it up and sets a new one.
 * GOPATH: sets GOPATH to the right place
* Write Create, that creates an empty set of configs and GOPATH.
* Write Update, that updates go versions.
* Write a set command for environments that allows:
 * set build commands
 * add local folders/repos as go paths ie: ``goworkon add ./thisfolder github.com/foo/bar`` Will set *thisfolder* as *github.com/foo/bar* in the environment.
* Write a set command for global that allows to set config such as:
 * set default environment.
* Write a way to set the default env to be used in .bashrc (or any .shellrc)
* Write a rebuild command for the current env
