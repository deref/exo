package resolvers

type ReconciliationResolver struct {
	Component *ComponentResolver
	Job       *TaskResolver
}
