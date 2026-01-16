package domain

type SnippetRepository interface {
	List() ([]*Snippet, error)
	FindByID(id int) (*Snippet, error)
	FindByCategory(categoryID int) ([]*Snippet, error)
	FindByTag(tagID int) ([]*Snippet, error)
	FindByLanguage(language string) ([]*Snippet, error)
	Search(value string) ([]*Snippet, error)
	Create(s *Snippet) error
	Update(s *Snippet) error
	Delete(id int) error
}

type CategoryRepository interface {
	List() ([]*Category, error)
	FindByID(id int) (*Category, error)
	FindByName(name string) (*Category, error)
	Create(c *Category) error
	Update(c *Category) error
	Delete(id int) error
}

type TagRepository interface {
	List() ([]*Tag, error)
	FindByID(id int) (*Tag, error)
	FindByName(name string) (*Tag, error)
	Create(t *Tag) error
	Update(t *Tag) error
	Delete(id int) error
}
