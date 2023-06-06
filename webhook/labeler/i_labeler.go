package labeler

type ILabeler interface {
	Label(annotations map[string]string) map[string]string
}
