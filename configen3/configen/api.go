package configen

type Process interface {
	Run(store *Context) error
}
