package eOverdracht.receiver

import rego.v1

# Only owner can update the pet's information
# Ownership information is provided as part of OPA's input
default allow := false

allow if {
	input.request.method = "POST"
}