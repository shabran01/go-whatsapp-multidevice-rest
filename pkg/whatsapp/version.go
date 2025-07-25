package whatsapp

type clientVersion struct {
	Major int
	Minor int
	Patch int
	Build int
}

var version = clientVersion{
	Major: 2,
	Minor: 2412,
	Patch: 54,
	Build: 0,
}
