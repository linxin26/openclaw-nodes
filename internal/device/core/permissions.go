package core

type PermissionState string

const (
	PermissionGranted       PermissionState = "granted"
	PermissionDenied        PermissionState = "denied"
	PermissionRestricted    PermissionState = "restricted"
	PermissionNotSupported  PermissionState = "not_supported"
	PermissionNotApplicable PermissionState = "not_applicable"
)
