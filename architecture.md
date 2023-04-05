# Architecture

## Controller

- Core of all data being passed around
- Event driven
- Kernel and gui listen and trigger events on the controller

### Interfaces

- Go API
- Pubsub frameworks
	- libp2p
- HTTP
- QUIC

## Listeners

### Input data

- Variables that are input and displayed to the user
	- Displayed in windows that can be moved and combined by the user
- Arrays displayed to user in excel fashion
- Displayed in a tree fashion? (on the side?)
- Accessible to the analysis and application loop
- Include matrix library if possible?
- Data can be tied to files
	- csv, excel, json (, toml, yaml?)
- API interface with other applications?

### Kernel

- Run chunks of code
- can reference variables and other formulas
	- Pass by value
- Ran in parallel
	- Dependency management

