# goworkon
A humble attempt at getting the functionality of virtualenv for go.

## Requirements
A working GO version >= 1.4

## Rationale
I often find myself working in different projects during the day, some might use different 
versions of go and most definitely all are better off having separate GOPATHS.

What I am attempting here is to automate a way to switch between GOPATHS and
to get new golang versions.

## Usage
###First run will ask you for a working goroot, so it can compile go versions.

####Creating environments:
``
goworkon --go-version 1.7 create envname gopathlocation
``

This will create:

* A GOPATH for *envname* 
* A config for *envname* (in $HOME/.local/share/goworkon/configs/envname.json
* If there is no 1.7 in $HOME/.local/share/goworkon/install/ it will be checked out and built.

####Switching to an environment:
``
. goactivate envname
``

This will (most likely only in bash):

* export PATH as $PATH:/this/go/version/bin/:$GOPATH/bin
* export GOPATH as $HOME/.local/share/goworkon/gopath
* export PS1 as $PS1(envname)$ (this requires PS1 to be exported)

as an alternative you can

``
goworkon switch envname
``

and you will see the variables that need to be set printed for you
to write whatever suits your shell.

####Un-switching
``
. goactivate
``

will return the environment to its former state.

as an alternative 

``
goworkon switch
``

###Listing
``
goworkon list
``

Will print a detailed list of environments.

#### Default environment

``
goworkong set default <envname>
. goactivate <envname>
``

Will set envname as the default environment, this environment will not change
PS1 variable and can be switched in .bashrc its intended to be used as the general one.


####Setting global settings
``
goworkon set <attribute> [value]
``

will set ``<attribute>`` settings value if its a valid one (you will get an error otherwise)
if value is not passed the attribute will be blanked.

####Setting environment settings
``
goworkon set envname@<attribute> [value]
``

will set ``<attribute>`` environment config value if its a valid one (you will get an error otherwise)
if value is not passed the attribute will be blanked.

Interesting settings:

* ``goworkon set envname@globalbin "true"`` sets a flag in the project that makes its $GOPATH/bin included
in all envs, very useful for tools that you build but want to keep separate.


####Updating a Go version:

``
goworkon update --go-version 1.7 [envname]
``

Will update the go version in use to 1.7.latest and then (in next release) rebuild envs using 1.7

``
goworkon update envname
``

Will udpate the env to the latest available version of go.

``
goworkon update --go-version 1.7
``

Will update all projects using 1.7.x to the latest version of 1.7

## To be implemented.

####TESTS
This was an attempt to replace bash scripts I was using for this
so I sort of just coded it in a couple of sittings so it needs
extensive tests

####Debug output
Debug log level should be setable and proper information should be added
to the logging.


####Setting build steps:

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

####More ideas:

* Write a rebuild command for the current env
