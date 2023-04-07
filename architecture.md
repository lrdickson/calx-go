# Architecture

## Controller

- Controls data passing
- calls listeners to update app as needed
- Kernel and gui listen and trigger events on the controller

### Details

- Variable as interface with implementations
	- formula
	- file
	- network
		- pubsub
		- server
			- http
			- etc
- Should variable be a separate package?
- Variable access functions
	- IterVariables
		- takes function that returns bool
		- true = continue
	- GetVariable
- Add event listeners as functions
- Events/Functions
	- Add
	- Delete
	- Rename
	- DataUpdate
	- Formula
		- CodeUpdate
- Add listener functions
	- Listen for all variables
	- Listen by variable name?
		- Pairs well with delete

## Variables

- Arrays displayed to user in excel fashion
- Displayed in a tree fashion? (on the side?)

### Inputs

- Options
- manual input
- Files
	- csv, excel, json (, toml, yaml?)
	- file import (saves to project)
	- File read (does not save to project)
- Network - subscribe
  - libp2p
  - mqtt
- User can select 1 or more options?
- manual input as default?
- validator formula?
- Or should the user just write a separate formula object as a validator?

### Formulas

- Run chunks of code
- can reference variables and other formulas
	- Pass by value
- Ran in parallel
	- Dependency management
- 0 or more output options?
- Inverse of input options
- continuous running options?
	- can make http requests, etc
	- out/put/post function?
- Include matrix library if possible?

