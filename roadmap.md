# Roadmap

## Ideas

- possibly have kernels for different languages like Jupyter?

### V3

#### Components

##### Input data

- Variables that are input and displayed to the user
	- Displayed in windows that can be moved and combined by the user
- Arrays displayed to user in excel fashion
- Displayed in a tree fashion? (on the side?)
- Accessible to the analysis and application loop
- Include matrix library if possible?
- Data can be tied to files
	- csv, excel, json (, toml, yaml?)
- API interface with other applications?

##### Formulas

- chunks of go code
- can reference variables and other formulas
	- Pass by value
- Ran in parallel
	- Dependency management

##### Application function

- To be implemented later
- Can be used to run an application using the input data
- When to run?
	- On loop?
	- At button press of user?
	- Provide options to user?

### V2

#### Components

##### Input data

- Variables that are input and displayed to the user
	- Displayed in windows that can be moved and combined by the user
- Arrays displayed to user in excel fashion
- Displayed in a tree fashion (on the side?)
- Accessible to the analysis and application loop
- Include matrix library if possible?
- Data can be tied to files
	- csv, excel, json (, toml, yaml?)
- API interface with other applications?

##### Analysis function

- Run every time the input data or analysis loop changes

##### Application function

- To be implemented later
- Can be used to run an application using the input data
- When to run?
	- On loop?
	- At button press of user?
	- Provide options to user?

### V1

- Cells treated as variables
- Functions can be defined separately
- Main loop
- Use # for cell references
	- Function definitions may not use cell references

- output portable exports using js and html

- import all math and science library
- Make worker server a constant file
- add worker functions as a separate file in the same package

