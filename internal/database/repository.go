package database

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository[T any, I any] struct {
	db     *gorm.DB
	entity T
}

func NewRepository[T any, I any](db *gorm.DB, entity T) *Repository[T, I] {
	return &Repository[T, I]{db: db, entity: entity}
}

func (r *Repository[T, I]) Create(t T) error {
	fmt.Println(r.db)
	return r.db.Create(t).Error
}

func (r *Repository[T, I]) Find(id I) error {
	return r.db.Find(r.entity, id).Error
}

func (r *Repository[T, I]) FindAll() error {
	return r.db.Find(r.entity).Error
}

func (r *Repository[T, I]) Update(t T) error {
	return r.db.Save(t).Error
}

func (r *Repository[T, I]) Delete(id I) error {
	return r.db.Delete(r.entity, id).Error
}
