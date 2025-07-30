package repositories

type PasteRepository interface {
	Create(dto )
	UpdateByTitle(title string)
}
