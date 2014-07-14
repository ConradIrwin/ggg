The start of an experimental port of [Pry](https://pryrepl.org) to golang.

# Status

This is currently vaporware, but I've done the following things:

* [x] get the current stack pointer (x86 only)
* [x] load the DWARF data for the current process (macho only)
* [x] use dwarf data to compute the canonical frame address
* [x] lookup the address of variables by name
* [ ] create interface values from DWARF types
* [ ] call functions
* [ ] add some UI

Much help appreciated!
