// Package merges renders the merge request collection
package merges

type MergesController struct {
	repo   MergesRepository
	render MergesRenderer
}

func NewMergesController(
	repo MergesRepository,
	render MergesRenderer,
) *MergesController {

	render.MakeBinding(repo)
	return &MergesController{
		repo:   repo,
		render: render,
	}
}

func (m *MergesController) Run() error {
	return m.render.Run()
}
