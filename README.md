# helvar-go

A Go library for the HelvarNET protocol.

Documentation for the protocol:
https://aca.im/driver_docs/Helvar/HelvarNet-Overview.pdf

## Notice

Originally was forked from Christoffer's Tombre [`helvar-go`](https://github.com/chtombre/helvar-go). Almost all code has been redesigned, therefore I do not consider that I am obliged to comply with the Apache 2.0 license and this
lib will come with MIT license.

Key differences:
- fixed design flaw in constructing/serializing messages
- switched from unreliable network interaction sync model based on `sync.Map` and `sync.Cond` to some kind of connection pool based on channels
- added missing checks for errors everywhere
- added handy shortcuts functions
- added support for color and color temperature commands
- added tests and some tools for testing

## Contributing

There are still lots of things you can improve, PRs and issues are welcome.
