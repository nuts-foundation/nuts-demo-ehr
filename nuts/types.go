package nuts

// NutsOrganization models the credentialSubject of a NutsOrganizationCredential.
type NutsOrganization struct {
	ID      string              `json:"id"`
	Details OrganizationDetails `json:"organization"`
}

type OrganizationDetails struct {
	Name string `json:"name"`
	City string `json:"city"`
}
