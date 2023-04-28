# Architecture

## Controller

- Controls data passing
- calls listeners to update app as needed
- Kernel and gui listen and trigger events on the controller

### Details

- Outputs are versioned
	- When object is updated, version is incremented
	- New output map created for that version
	- Known output from non-dependents are copied from the previous
		- Last compatable version for an object output is recorded
	- New version number is reported so workers can stop producing outdated output
		- Maybe report which ones need to stop/can continue?
	- Output of updated object sent to all dependents with version number
	- Object reply with new data to fill out new map
	- Does dependency versioning also need to be tracked?

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
	- OutputChanged
	- MetaDataChanged
		- CodeUpdate
- Add listener functions
	- Listen for all variables
	- Listen by variable id
		- Pairs well with delete

### Kernel

- objects represented as an ID
- Kernels may keep an internal representation of variables as necessary

- Object ID
  - Allows it to communicate across separate devices

## Objects

- Arrays displayed to user in excel fashion
- Displayed in a tree fashion? (on the side?)

### Interfaces

- Object
	- Close
	- Name?
- Consumer Interface
	- Required
		- Eval
	- Provided
		- Dependencies
			- Listens to producer events
- Producer
	- Required
		- Output
	- Provided
		- Dependents?

### Types

- Formula (consumer + producer)
- selector (multiplexer)
- validator (filter)
- collector?

- How much of this can be covered by specialized objects vs formulas

### Inputs

- Producer

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

- Consumer + Producer

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

### Outputs

- Producers

### How to access objects

- Give direct access to the object?
	- Every new object will require a function
- Access all fields from controller
- Set fields through the controller and query fields from the object?

### How to extend objects

- Interfaces?
- anonymous functions?
- embedded objects?

### How to key objects

- The name?
	- Keeping track of name will be critical
- The object pointer?
- An ID variable?

### How to keep track of names

- Is it necessary
	- Makes it easy to verify name uniqueness
- In a map?
- Iterate through variables?

## Display

- Waterfall?
- Inputs at top flow data down

## Language server

- Instead of trying to make a decent editor, interface with them as a language server


