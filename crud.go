package kitty

import (
	vd "github.com/bytedance/go-tagexpr/validator"
	"github.com/jinzhu/gorm"
)

//CRUDInterface ...
type CRUDInterface interface {
	Do(*SearchCondition, string, Context) (interface{}, error)
}

// SuccessCallback 执行成功后，回调。返回error后，回滚事务
type SuccessCallback func(*Structs, *gorm.DB) error

// CRUD配置
type config struct {
	strs   *Structs         //模型结构
	search *SearchCondition //查询条件
	db     *gorm.DB         //db
	ctx    Context          //上下文
	callbk SuccessCallback  //成功回调
}

type crud struct {
	*config
}

func newcrud(conf *config) *crud {
	return &crud{conf}
}
func (crud *crud) queryExpr() (interface{}, error) {
	var (
		s      = crud.strs
		search = crud.search
		db     = crud.db
		c      = crud.ctx
	)

	if err := vd.Validate(s.raw); err != nil {
		return nil, err
	}

	if err := getter(s, make(map[string]interface{}), db, c); err != nil {
		return nil, err
	}

	kittys := &kittys{
		ModelStructs: s,
		db:           db,
	}
	if err := kittys.parse(s); err != nil {
		return nil, err
	}

	var qry qry
	//	if len(kittys.kittys) > 1 {
	qry = evalJoin(s, kittys, search, db)
	//} else {
	//	qry = evalSimpleQry(s, kittys, search, db)
	//}
	return qry.prepare().QueryExpr(), nil
}

func (crud *crud) queryObj() (interface{}, error) {
	var (
		s      = crud.strs
		search = crud.search
		db     = crud.db
		c      = crud.ctx
		callbk = crud.callbk
	)

	if err := vd.Validate(s.raw); err != nil {
		return nil, err
	}

	if err := getter(s, make(map[string]interface{}), db, c); err != nil {
		return nil, err
	}

	Page := &Page{}
	if f, ok := s.FieldOk("Page"); ok {
		Page.Page = f.Value().(uint32)
	}
	if f, ok := s.FieldOk("Limit"); ok {
		Page.Limit = f.Value().(uint32)
	}
	if Page.Limit > 0 && Page.Page > 0 {
		search.Page = Page
	}

	kittys := &kittys{
		ModelStructs: s,
		db:           db,
	}
	if err := kittys.parse(s); err != nil {
		return nil, err
	}

	var qry qry
	//if len(kittys.kittys) > 1 {
	qry = evalJoin(s, kittys, search, db)
	//	} else {
	//	qry = evalSimpleQry(s, kittys, search, db)
	//}

	var (
		res interface{}
		err error
	)
	res, err = execqry(qry, kittys.multiResult)
	if err != nil || res == nil {
		return nil, err
	}

	if len(kittys.resultField) > 0 {
		if err = s.Field(kittys.resultField).Set(res); err != nil {
			return nil, err
		}
	}

	params := make(map[string]interface{})
	params["ms"] = s
	params["kittys"] = kittys
	if err = setter(s, params, db, c); err != nil {
		return nil, err
	}

	if callbk != nil {
		if err = callbk(s, db); err != nil {
			return nil, err
		}
	}

	if f, ok := s.FieldOk("Pages"); ok && search.Page != nil {
		f.Set(search.Page)
	}

	return s.raw, nil
}

// CreateObj ...
func (crud *crud) createObj() (interface{}, error) {
	var (
		s      = crud.strs
		search = crud.search
		db     = crud.db
		c      = crud.ctx
		callbk = crud.callbk
	)

	if err := vd.Validate(s.raw); err != nil {
		return nil, err
	}

	if err := getter(s, make(map[string]interface{}), db, c); err != nil {
		return nil, err
	}

	kittys := &kittys{
		ModelStructs: s,
		db:           db,
	}
	if err := kittys.parse(s); err != nil {
		return nil, err
	}

	qry := &simpleQuery{
		db:           db,
		ModelStructs: s,
		search:       search,
		Result:       kittys.master().structs,
	}
	for _, v := range kittys.kittys {
		if !v.Master {
			qry.Next = append(qry.Next, &simpleQuery{
				db:           db,
				ModelStructs: s,
				search:       &SearchCondition{},
				Result:       v.structs,
			})
		}
	}
	res, err := qry.create()
	if err != nil {
		return nil, err
	}
	for _, v := range kittys.kittys {
		f := s.Field(v.FieldName)
		f.Set(v.structs.raw)
	}

	if len(kittys.resultField) > 0 {
		if err = s.Field(kittys.resultField).Set(res); err != nil {
			return nil, err
		}
	}

	params := make(map[string]interface{})
	params["ms"] = s
	params["kittys"] = kittys
	if err = setter(s, params, db, c); err != nil {
		return nil, err
	}
	if callbk != nil {
		if err = callbk(s, db); err != nil {
			return nil, err
		}
	}

	return s.raw, nil
}

func (crud *crud) updateObj() (interface{}, error) {
	var (
		s      = crud.strs
		search = crud.search
		db     = crud.db
		c      = crud.ctx
		callbk = crud.callbk
	)

	if err := vd.Validate(s.raw); err != nil {
		return nil, err
	}

	if err := getter(s, make(map[string]interface{}), db, c); err != nil {
		return nil, err
	}

	kittys := &kittys{
		ModelStructs: s,
		db:           db,
	}
	if err := kittys.parse(s); err != nil {
		return nil, err
	}

	qry := &simpleQuery{
		db:           db,
		ModelStructs: s,
		search:       search,
		Result:       kittys.master().structs,
	}
	for _, v := range kittys.kittys {
		if !v.Master {
			qry.Next = append(qry.Next, &simpleQuery{
				db:           db,
				ModelStructs: s,
				search:       &SearchCondition{},
				Result:       v.structs,
			})
		}
	}

	if err := qry.update(); err != nil {
		return nil, err
	}
	params := make(map[string]interface{})
	params["ms"] = s
	params["kittys"] = kittys
	if err := setter(s, params, db, c); err != nil {
		return nil, err
	}

	if callbk != nil {
		if err := callbk(s, db); err != nil {
			return nil, err
		}
	}

	return s.raw, nil
}

//
func queryObj(s *Structs, search *SearchCondition, db *gorm.DB, c Context) (interface{}, error) {
	crud := newcrud(&config{
		strs:   s,
		search: search,
		db:     db,
		ctx:    c,
	})
	return crud.queryObj()
}
