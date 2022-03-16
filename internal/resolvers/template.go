package resolvers

import "github.com/deref/exo/internal/template"

type TemplateResolver struct {
	Q *RootResolver
	template.TemplateDescription
}

func (r *RootResolver) AllTemplates() []*TemplateResolver {
	return templateDescriptionsToTemplateResolvers(r, template.GetTemplateDescriptions())
}

func templateDescriptionsToTemplateResolvers(r *RootResolver, descriptions []template.TemplateDescription) []*TemplateResolver {
	resolvers := make([]*TemplateResolver, len(descriptions))
	for i, descrption := range descriptions {
		resolvers[i] = &TemplateResolver{
			Q:                   r,
			TemplateDescription: descrption,
		}
	}
	return resolvers
}
