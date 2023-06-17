// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"example/bluebean-go/model"
)

func newFacility(db *gorm.DB, opts ...gen.DOOption) facility {
	_facility := facility{}

	_facility.facilityDo.UseDB(db, opts...)
	_facility.facilityDo.UseModel(&model.Facility{})

	tableName := _facility.facilityDo.TableName()
	_facility.ALL = field.NewAsterisk(tableName)
	_facility.ID = field.NewInt64(tableName, "Id")
	_facility.Name = field.NewString(tableName, "Name")
	_facility.Address = field.NewString(tableName, "Address")
	_facility.ImageURL = field.NewString(tableName, "ImageUrl")
	_facility.CityID = field.NewInt64(tableName, "CityId")
	_facility.CreatorID = field.NewInt64(tableName, "CreatorId")

	_facility.fillFieldMap()

	return _facility
}

type facility struct {
	facilityDo facilityDo

	ALL       field.Asterisk
	ID        field.Int64
	Name      field.String
	Address   field.String
	ImageURL  field.String
	CityID    field.Int64
	CreatorID field.Int64

	fieldMap map[string]field.Expr
}

func (f facility) Table(newTableName string) *facility {
	f.facilityDo.UseTable(newTableName)
	return f.updateTableName(newTableName)
}

func (f facility) As(alias string) *facility {
	f.facilityDo.DO = *(f.facilityDo.As(alias).(*gen.DO))
	return f.updateTableName(alias)
}

func (f *facility) updateTableName(table string) *facility {
	f.ALL = field.NewAsterisk(table)
	f.ID = field.NewInt64(table, "Id")
	f.Name = field.NewString(table, "Name")
	f.Address = field.NewString(table, "Address")
	f.ImageURL = field.NewString(table, "ImageUrl")
	f.CityID = field.NewInt64(table, "CityId")
	f.CreatorID = field.NewInt64(table, "CreatorId")

	f.fillFieldMap()

	return f
}

func (f *facility) WithContext(ctx context.Context) *facilityDo { return f.facilityDo.WithContext(ctx) }

func (f facility) TableName() string { return f.facilityDo.TableName() }

func (f facility) Alias() string { return f.facilityDo.Alias() }

func (f *facility) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := f.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (f *facility) fillFieldMap() {
	f.fieldMap = make(map[string]field.Expr, 6)
	f.fieldMap["Id"] = f.ID
	f.fieldMap["Name"] = f.Name
	f.fieldMap["Address"] = f.Address
	f.fieldMap["ImageUrl"] = f.ImageURL
	f.fieldMap["CityId"] = f.CityID
	f.fieldMap["CreatorId"] = f.CreatorID
}

func (f facility) clone(db *gorm.DB) facility {
	f.facilityDo.ReplaceConnPool(db.Statement.ConnPool)
	return f
}

func (f facility) replaceDB(db *gorm.DB) facility {
	f.facilityDo.ReplaceDB(db)
	return f
}

type facilityDo struct{ gen.DO }

func (f facilityDo) Debug() *facilityDo {
	return f.withDO(f.DO.Debug())
}

func (f facilityDo) WithContext(ctx context.Context) *facilityDo {
	return f.withDO(f.DO.WithContext(ctx))
}

func (f facilityDo) ReadDB() *facilityDo {
	return f.Clauses(dbresolver.Read)
}

func (f facilityDo) WriteDB() *facilityDo {
	return f.Clauses(dbresolver.Write)
}

func (f facilityDo) Session(config *gorm.Session) *facilityDo {
	return f.withDO(f.DO.Session(config))
}

func (f facilityDo) Clauses(conds ...clause.Expression) *facilityDo {
	return f.withDO(f.DO.Clauses(conds...))
}

func (f facilityDo) Returning(value interface{}, columns ...string) *facilityDo {
	return f.withDO(f.DO.Returning(value, columns...))
}

func (f facilityDo) Not(conds ...gen.Condition) *facilityDo {
	return f.withDO(f.DO.Not(conds...))
}

func (f facilityDo) Or(conds ...gen.Condition) *facilityDo {
	return f.withDO(f.DO.Or(conds...))
}

func (f facilityDo) Select(conds ...field.Expr) *facilityDo {
	return f.withDO(f.DO.Select(conds...))
}

func (f facilityDo) Where(conds ...gen.Condition) *facilityDo {
	return f.withDO(f.DO.Where(conds...))
}

func (f facilityDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *facilityDo {
	return f.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (f facilityDo) Order(conds ...field.Expr) *facilityDo {
	return f.withDO(f.DO.Order(conds...))
}

func (f facilityDo) Distinct(cols ...field.Expr) *facilityDo {
	return f.withDO(f.DO.Distinct(cols...))
}

func (f facilityDo) Omit(cols ...field.Expr) *facilityDo {
	return f.withDO(f.DO.Omit(cols...))
}

func (f facilityDo) Join(table schema.Tabler, on ...field.Expr) *facilityDo {
	return f.withDO(f.DO.Join(table, on...))
}

func (f facilityDo) LeftJoin(table schema.Tabler, on ...field.Expr) *facilityDo {
	return f.withDO(f.DO.LeftJoin(table, on...))
}

func (f facilityDo) RightJoin(table schema.Tabler, on ...field.Expr) *facilityDo {
	return f.withDO(f.DO.RightJoin(table, on...))
}

func (f facilityDo) Group(cols ...field.Expr) *facilityDo {
	return f.withDO(f.DO.Group(cols...))
}

func (f facilityDo) Having(conds ...gen.Condition) *facilityDo {
	return f.withDO(f.DO.Having(conds...))
}

func (f facilityDo) Limit(limit int) *facilityDo {
	return f.withDO(f.DO.Limit(limit))
}

func (f facilityDo) Offset(offset int) *facilityDo {
	return f.withDO(f.DO.Offset(offset))
}

func (f facilityDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *facilityDo {
	return f.withDO(f.DO.Scopes(funcs...))
}

func (f facilityDo) Unscoped() *facilityDo {
	return f.withDO(f.DO.Unscoped())
}

func (f facilityDo) Create(values ...*model.Facility) error {
	if len(values) == 0 {
		return nil
	}
	return f.DO.Create(values)
}

func (f facilityDo) CreateInBatches(values []*model.Facility, batchSize int) error {
	return f.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (f facilityDo) Save(values ...*model.Facility) error {
	if len(values) == 0 {
		return nil
	}
	return f.DO.Save(values)
}

func (f facilityDo) First() (*model.Facility, error) {
	if result, err := f.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Facility), nil
	}
}

func (f facilityDo) Take() (*model.Facility, error) {
	if result, err := f.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Facility), nil
	}
}

func (f facilityDo) Last() (*model.Facility, error) {
	if result, err := f.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Facility), nil
	}
}

func (f facilityDo) Find() ([]*model.Facility, error) {
	result, err := f.DO.Find()
	return result.([]*model.Facility), err
}

func (f facilityDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Facility, err error) {
	buf := make([]*model.Facility, 0, batchSize)
	err = f.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (f facilityDo) FindInBatches(result *[]*model.Facility, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return f.DO.FindInBatches(result, batchSize, fc)
}

func (f facilityDo) Attrs(attrs ...field.AssignExpr) *facilityDo {
	return f.withDO(f.DO.Attrs(attrs...))
}

func (f facilityDo) Assign(attrs ...field.AssignExpr) *facilityDo {
	return f.withDO(f.DO.Assign(attrs...))
}

func (f facilityDo) Joins(fields ...field.RelationField) *facilityDo {
	for _, _f := range fields {
		f = *f.withDO(f.DO.Joins(_f))
	}
	return &f
}

func (f facilityDo) Preload(fields ...field.RelationField) *facilityDo {
	for _, _f := range fields {
		f = *f.withDO(f.DO.Preload(_f))
	}
	return &f
}

func (f facilityDo) FirstOrInit() (*model.Facility, error) {
	if result, err := f.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Facility), nil
	}
}

func (f facilityDo) FirstOrCreate() (*model.Facility, error) {
	if result, err := f.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Facility), nil
	}
}

func (f facilityDo) FindByPage(offset int, limit int) (result []*model.Facility, count int64, err error) {
	result, err = f.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = f.Offset(-1).Limit(-1).Count()
	return
}

func (f facilityDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = f.Count()
	if err != nil {
		return
	}

	err = f.Offset(offset).Limit(limit).Scan(result)
	return
}

func (f facilityDo) Scan(result interface{}) (err error) {
	return f.DO.Scan(result)
}

func (f facilityDo) Delete(models ...*model.Facility) (result gen.ResultInfo, err error) {
	return f.DO.Delete(models)
}

func (f *facilityDo) withDO(do gen.Dao) *facilityDo {
	f.DO = *do.(*gen.DO)
	return f
}
