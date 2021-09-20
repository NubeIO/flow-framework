package database

func (d *GormDatabase) CreateModel(model interface{}) error {
	query := d.DB.Create(model)
	return query.Error
}

func (d *GormDatabase) GetModel(id int, model *interface{}) error {
	query := d.DB.Where("id = ? ", id).First(&model)
	return query.Error
}

func (d *GormDatabase) GetModelList(modelList interface{}) error {
	query := d.DB.Find(modelList)
	return query.Error
}

func (d *GormDatabase) DeleteModel(id int, model interface{}) (bool, error) {
	query := d.DB.Where("id = ? ", id).Delete(model)
	return query.RowsAffected > 0, query.Error
}

func (d *GormDatabase) UpdateModel(id int, model *interface{}) error {
	query := d.DB.Model(&model).Where("id = ?", id).Updates(model)
	return query.Error
}

func (d *GormDatabase) DropBlockList(model *interface{}) (bool, error) {
	query := d.DB.Where("1 = 1").Delete(model)
	return query.RowsAffected > 0, query.Error
}
