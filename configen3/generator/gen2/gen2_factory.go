package gen2

type templateFactory interface {
	createMainTemplate(ctx *gen2Context) mainTemplate
	createComItemInConfigFuncTemplate(ctx *gen2Context) comItemInConfigFuncTemplate
	createComItemAfterConfigFuncTemplate(ctx *gen2Context) comItemAfterConfigFuncTemplate
	createComFieldGetterTemplate(ctx *gen2Context) comFieldGetterTemplate
}

type templateFactoryImpl struct {
}

func (o *templateFactoryImpl) _Impl() templateFactory {
	return o
}

func (o *templateFactoryImpl) loadDefaultImports(ctx *gen2Context) {

	const noHash = true

	ib := &ctx.importsBuilder
	// ib.AddPackagePath("strings", noHash)
	// ib.AddPackagePath("text/template", noHash)
	ib.AddPackagePath("github.com/bitwormhole/starter/application", noHash)
	ib.AddPackagePath("github.com/bitwormhole/starter/application/config", noHash)
	ib.AddPackagePath("github.com/bitwormhole/starter/lang", noHash)
	ib.AddPackagePath("github.com/bitwormhole/starter/util", noHash)

}

func (o *templateFactoryImpl) createMainTemplate(ctx *gen2Context) mainTemplate {
	o.loadDefaultImports(ctx)
	return (&mainTemplateImpl{}).init(ctx)
}

func (o *templateFactoryImpl) createComItemInConfigFuncTemplate(ctx *gen2Context) comItemInConfigFuncTemplate {
	return (&comItemInConfigFuncTemplateImpl{}).init(ctx)
}

func (o *templateFactoryImpl) createComItemAfterConfigFuncTemplate(ctx *gen2Context) comItemAfterConfigFuncTemplate {
	return (&comItemAfterConfigFuncTemplateImpl{}).init(ctx)
}

func (o *templateFactoryImpl) createComFieldGetterTemplate(ctx *gen2Context) comFieldGetterTemplate {
	return (&comFieldGetterTemplateImpl{}).init(ctx)
}
