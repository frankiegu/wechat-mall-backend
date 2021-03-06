package service

import (
	"wechat-mall-backend/dbops"
	"wechat-mall-backend/model"
)

type ICategoryService interface {
	GetCategoryList(pid, page, size int) (*[]model.WechatMallCategoryDO, int)
	GetCategoryById(id int) *model.WechatMallCategoryDO
	GetCategoryByName(name string) *model.WechatMallCategoryDO
	AddCategory(category *model.WechatMallCategoryDO)
	UpdateCategory(category *model.WechatMallCategoryDO)
}

type categoryService struct {
}

func NewCategoryService() ICategoryService {
	service := &categoryService{}
	return service
}

func (cs *categoryService) GetCategoryList(pid, page, size int) (*[]model.WechatMallCategoryDO, int) {
	cateList, err := dbops.QueryCategoryList(pid, page, size)
	if err != nil {
		panic(err)
	}
	total, err := dbops.CountCategory(pid)
	if err != nil {
		panic(err)
	}
	return cateList, total
}

func (cs *categoryService) GetCategoryById(id int) *model.WechatMallCategoryDO {
	category, err := dbops.QueryCategoryById(id)
	if err != nil {
		panic(err)
	}
	return category
}

func (cs *categoryService) GetCategoryByName(name string) *model.WechatMallCategoryDO {
	category, err := dbops.QueryCategoryByName(name)
	if err != nil {
		panic(err)
	}
	return category
}

func (cs *categoryService) AddCategory(category *model.WechatMallCategoryDO) {
	err := dbops.InsertCategory(category)
	if err != nil {
		panic(err)
	}
}

func (cs *categoryService) UpdateCategory(category *model.WechatMallCategoryDO) {
	err := dbops.UpdateCategoryById(category)
	if err != nil {
		panic(err)
	}
	syncSubCategoryAndGoodsOnline(category.ParentId, category.Id, category.Online)
}

// 同步其子分类和商品的上下架状态
func syncSubCategoryAndGoodsOnline(parentId, categoryId, online int) {
	if parentId == 0 {
		err := dbops.UpdateSubCategoryOnline(categoryId, online)
		if err != nil {
			panic(err)
		}
		ids, err := dbops.QuerySubCategoryByParentId(categoryId)
		if err != nil {
			panic(err)
		}
		for _, v := range *ids {
			err := dbops.UpdateCategoryGoodsOnlineStatus(v, online)
			if err != nil {
				panic(err)
			}
		}
	} else {
		err := dbops.UpdateCategoryGoodsOnlineStatus(categoryId, online)
		if err != nil {
			panic(err)
		}
	}
}
