package annotator

type IAnnotator interface {
	Annotate(annotations map[string]string) map[string]string
}
