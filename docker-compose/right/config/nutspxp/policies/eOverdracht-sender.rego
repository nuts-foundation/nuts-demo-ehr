package eOverdracht.sender

import rego.v1

# Only owner can update the pet's information
# Ownership information is provided as part of OPA's input
default allow := false

# eOVerdracht expects the following data: path to [actions] mapping

# example when http methods are stored in nuts-pxp
allow if {
	input.request.method in input.external[input.request.path]
}