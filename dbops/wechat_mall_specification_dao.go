package dbops

import (
	"strconv"
	"time"
	"wechat-mall-backend/model"
)

const specColumnList = `
id, name, description, unit, standard, is_del, create_time, update_time
`

func QuerySpecificationList(page, size int) (*[]model.WechatMallSpecificationDO, error) {
	sql := "SELECT " + specColumnList + " FROM wechat_mall_specification WHERE is_del = 0 LIMIT ?, ?"
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return nil, err
	}
	rows, err := stmt.Query((page-1)*size, size)
	if err != nil {
		return nil, err
	}
	var specList []model.WechatMallSpecificationDO
	for rows.Next() {
		spec := model.WechatMallSpecificationDO{}
		err := rows.Scan(&spec.Id, &spec.Name, &spec.Description, &spec.Unit, &spec.Standard, &spec.Del, &spec.CreateTime, &spec.UpdateTime)
		if err != nil {
			return nil, err
		}
		specList = append(specList, spec)
	}
	return &specList, nil
}

func CountSpecification() (int, error) {
	sql := "SELECT COUNT(*) FROM wechat_mall_specification WHERE is_del = 0"
	rows, err := dbConn.Query(sql)
	if err != nil {
		return 0, err
	}
	total := 0
	for rows.Next() {
		err := rows.Scan(&total)
		if err != nil {
			return 0, err
		}
	}
	return total, nil
}

func AddSpecification(spec *model.WechatMallSpecificationDO) error {
	sql := "INSERT INTO wechat_mall_specification ( " + specColumnList[4:] + " ) VALUES(?, ?, ?, ?, ?, ?, ?)"
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(spec.Name, spec.Description, spec.Unit, spec.Standard, 0, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func QuerySpecificationById(id int) (*model.WechatMallSpecificationDO, error) {
	sql := "SELECT " + specColumnList + " FROM wechat_mall_specification WHERE id = " + strconv.Itoa(id)
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	spec := model.WechatMallSpecificationDO{}
	for rows.Next() {
		err := rows.Scan(&spec.Id, &spec.Name, &spec.Description, &spec.Unit, &spec.Standard, &spec.Del, &spec.CreateTime, &spec.UpdateTime)
		if err != nil {
			return nil, err
		}
	}
	return &spec, nil
}

func QuerySpecificationByName(name string) (*model.WechatMallSpecificationDO, error) {
	sql := "SELECT " + specColumnList + " FROM wechat_mall_specification WHERE is_del = 0 AND name = '" + name + "'"
	rows, err := dbConn.Query(sql)
	if err != nil {
		return nil, err
	}
	spec := model.WechatMallSpecificationDO{}
	for rows.Next() {
		err := rows.Scan(&spec.Id, &spec.Name, &spec.Description, &spec.Unit, &spec.Standard, &spec.Del, &spec.CreateTime, &spec.UpdateTime)
		if err != nil {
			return nil, err
		}
	}
	return &spec, nil
}

func UpdateSpecificationById(spec *model.WechatMallSpecificationDO) error {
	sql := `
UPDATE wechat_mall_specification 
SET name = ?, description = ?, unit = ?, standard = ?, is_del = ?, update_time = ? 
WHERE id = ?
`
	stmt, err := dbConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(spec.Name, spec.Description, spec.Unit, spec.Standard, spec.Del, time.Now(), spec.Id)
	if err != nil {
		return err
	}
	return nil
}
