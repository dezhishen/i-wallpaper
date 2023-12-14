package provider

type sourceType int

const (
	FileSource sourceType = iota
	UrlSource
)

type ImageSource struct {
	Type   sourceType
	Source string
}

type Provider interface {
	GetTody() (*ImageSource, error)
	Random() (*ImageSource, error)
}
