package service

import (
	"wechat-mall-backend/dbops"
	"wechat-mall-backend/defs"
	"wechat-mall-backend/model"
)

type ICartService interface {
	DoEditCart(userId, goodsId, skuId, num int)
	GetCartGoods(userId, page, size int) (*[]defs.PortalCartGoodsVO, int)
}

type cartService struct {
}

func NewCartService() ICartService {
	service := cartService{}
	return &service
}

func (s *cartService) DoEditCart(userId, goodsId, skuId, num int) {
	if num == 0 {
		return
	}
	cartDO, err := dbops.QueryCartByParams(userId, goodsId, skuId)
	if err != nil {
		panic(err)
	}
	if num > 0 {
		if cartDO.Id == 0 {
			userCartDO := model.WechatMallUserCartDO{}
			userCartDO.UserId = userId
			userCartDO.GoodsId = goodsId
			userCartDO.SkuId = skuId
			userCartDO.Num = num
			err = dbops.AddUserCart(&userCartDO)
		} else {
			cartDO.Num += num
			err = dbops.UpdateCartById(cartDO)
		}
	} else {
		if cartDO.Id == 0 {
			return
		}
		if cartDO.Num+num > 0 {
			cartDO.Num += num
			err = dbops.UpdateCartById(cartDO)
		} else {
			cartDO.Del = 1
			err = dbops.UpdateCartById(cartDO)
		}
	}
	panic(err)
}

func (s *cartService) GetCartGoods(userId, page, size int) (*[]defs.PortalCartGoodsVO, int) {
	cartList, err := dbops.QueryCartList(userId, page, size)
	if err != nil {
		panic(err)
	}
	total, err := dbops.CountCartGoods(userId)
	if err != nil {
		panic(err)
	}
	cartGoodsVOList := []defs.PortalCartGoodsVO{}
	for _, v := range *cartList {
		goodsDO, err := dbops.QueryGoodsById(v.GoodsId)
		if err != nil {
			panic(err)
		}
		skuDO, err := dbops.GetSKUById(v.SkuId)
		if err != nil {
			panic(err)
		}
		status := 0
		if goodsDO.Id == 0 || goodsDO.Online == 0 || skuDO.Id == 0 || skuDO.Online == 0 {
			status = 2
		} else {
			if skuDO.Stock < v.Num {
				status = 1
			}
		}
		cartGoodsVO := defs.PortalCartGoodsVO{}
		cartGoodsVO.GoodsId = v.Id
		cartGoodsVO.Title = goodsDO.Title
		cartGoodsVO.Price = goodsDO.Price
		cartGoodsVO.DiscountPrice = goodsDO.DiscountPrice
		cartGoodsVO.Picture = goodsDO.Picture
		cartGoodsVO.Tags = goodsDO.Tags
		cartGoodsVO.SkuId = v.SkuId
		cartGoodsVO.Specs = skuDO.Specs
		cartGoodsVO.Num = v.Num
		cartGoodsVO.Status = status
		cartGoodsVOList = append(cartGoodsVOList, cartGoodsVO)
	}
	return &cartGoodsVOList, total
}