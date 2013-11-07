fat
===

Fat fields and structs for go

Idea
---

have one struct and use it:

  - to validate data
  - to set the values
  - to have default values
  - to allow further useful things 3rd party libraries could do with structs, like:
    -  fill with values from posted forms
    -  use them as placeholders in templates (with escaping)
    -  construct urls based on them
    -  to database operations, like creating a table o querying

a fat field is a field that knows about its name, its validation, its type and
the name of its struct, all without resorting to field tags and without code duplication

The third party libraries could allow operations for single fat fields or more of them
by associating the field path (which includes the field name and the struct type name)
with special field tags in a registry, so that every struct needs to be registered once.

Then the library could offer top level functions to return specialized objects for a given
field that knows its path and therefor those fat structs might be extended from the outside
with arbitrary data.

Examples
--------

see example directory

