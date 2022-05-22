package models

import "Yaratam/internal/domain"

type File struct {
	Path string `db:"paths"`
}

func (f File) Domain() *domain.File {
	return &domain.File{
		Path: f.Path,
	}
}

type Files []*File

func (ff Files) Domain() []*domain.File {
	dd := make([]*domain.File, 0)
	for _, v := range ff {
		dd = append(dd, v.Domain())
	}

	return dd
}

type Path struct {
	ID          int    `db:"id"`
	DisplayName string `db:"display_name"`
}

func (p Path) Domain() *domain.Path {
	return &domain.Path{
		ID:          p.ID,
		DisplayName: p.DisplayName,
	}
}

type Paths []*Path

func (pp Paths) Domain() []*domain.Path {
	dd := make([]*domain.Path, 0)
	for _, v := range pp {
		dd = append(dd, v.Domain())
	}
	return dd
}
