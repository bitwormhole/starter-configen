package gen1

type templateFactory interface {
	createMainTemplate(ctx *gen1Context) mainTemplate
	createImportListTemplate(ctx *gen1Context) importListTemplate
	createStructListTemplate(ctx *gen1Context) structListTemplate
}

type templateFactoryImpl struct {
}

func (o *templateFactoryImpl) _Impl() templateFactory {
	return o
}

func (o *templateFactoryImpl) createMainTemplate(ctx *gen1Context) mainTemplate {
	return (&mainTemplateImpl{}).init(ctx)
}

func (o *templateFactoryImpl) createImportListTemplate(ctx *gen1Context) importListTemplate {
	return (&importListTemplateImpl{}).init(ctx)
}

func (o *templateFactoryImpl) createStructListTemplate(ctx *gen1Context) structListTemplate {
	return (&structListTemplateImpl{}).init(ctx)
}
