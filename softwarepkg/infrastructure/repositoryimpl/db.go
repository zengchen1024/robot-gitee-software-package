package repositoryimpl

import "github.com/opensourceways/robot-gitee-software-package/softwarepkg/infrastructure/postgresql"

type dbClient interface {
	Insert(filter, result interface{}) error
	Count(filter interface{}) (int, error)
	GetRecords(filter, result interface{}, p postgresql.Pagination, sort []postgresql.SortByColumn) error
	GetRecord(filter, result interface{}) error
	UpdateRecord(filter, update interface{}) error
	DeleteRecord(filter interface{}) error
	IsRowNotFound(err error) bool
	IsRowExists(err error) bool
}
