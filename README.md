# idl-conv

## Features

* Convert one or more Go type declarations to protobuf message declarations

## TODO

* dependency resoution

    * topologically sort types

* `.go` file parsing: find Go type declaration in given Go source file

* type locator: program automatically find target types in given locations
    
    * `idl-conv /path/to/go-project/src TypeToConvert`

* namesapce resolution: use embedded message declaration as namesapce

## Future

* rearrange field sequence to compact payload size

* support multiple programming languages

* support map struct code generation

* support multiple idl languages

