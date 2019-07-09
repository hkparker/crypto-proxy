crypto-proxy
============

This application creates a proxy in front of a cryptocurrency node and runs functions on messages as they are recieved.  Currently the only supported cryptocurrency is Bitcoin, and the only message supported is `filterload`, which is written to a file.

New currencies are added to the `parsers` directory, and new actions are added to the `actions` directory.
