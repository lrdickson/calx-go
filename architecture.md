# Architecture

## Controller

- Core of all data being passed around
- Event driven
- Kernel and gui listen and trigger events on the controller

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

