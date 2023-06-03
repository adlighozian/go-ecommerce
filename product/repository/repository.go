package repository

import (
	"database/sql"
	"errors"
	"log"
	"product-go/helper/failerror"
	"product-go/helper/timeout"
	"product-go/model"
	"product-go/publisher"
	"time"
)

type repository struct {
	db   *sql.DB
	sent publisher.Publisher
}

func NewRepository(db *sql.DB, sent publisher.Publisher) Repositorier {
	return &repository{
		db:   db,
		sent: sent,
	}
}

func (repo *repository) GetProduct(req model.ProductSearch) ([]model.Product, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	querySearchProduct := `select p.id, p.store_id, p.category_id, p.size_id, p.color_id ,p.name, p.brand, p.subtitle, p.description ,p.unit_price ,p.status ,p.stock ,p.sku ,p.weight ,p.created_at ,p.updated_at from products p  join categories c on c.id = p.category_id join product_sizes s on s.id = p.size_id join product_colors co on co.id = p.color_id where c.name like '%' || $1 || '%' and p.name like '%' || $2 || '%' and p.brand like '%' || $3 || '%'`

	result, err := repo.db.QueryContext(ctx, querySearchProduct, req.Category, req.Name, req.Brand)
	failerror.FailError(err, "error query")

	var data = []model.Product{}

	for result.Next() {
		var temp model.ProductResult
		result.Scan(&temp.Id, &temp.StoreID, &temp.CategoryID, &temp.SizeID, &temp.ColorID, &temp.Name, &temp.Brand, &temp.Subtitle, &temp.Description, &temp.UnitPrice, &temp.Status, &temp.Stock, &temp.Sku, &temp.Weight, &temp.Created_at, &temp.Update_at)

		queryImage := `select i.name, i.image_url from image i join image_products ip on ip.image_id = i.id join products p on p.id  = ip.product_id  where ip.product_id  = $1`

		rows, err := repo.db.QueryContext(ctx, queryImage, temp.Id)
		failerror.FailError(err, "error query")

		var imageProduct = []model.ProductImage{}
		for rows.Next() {
			var images model.ProductImage
			err := rows.Scan(&images.Name, &images.ImageURL)
			failerror.FailError(err, "error scan")
			imageProduct = append(imageProduct, images)
		}

		data = append(data, model.Product{
			Id:           temp.Id,
			StoreID:      temp.StoreID,
			CategoryID:   temp.CategoryID,
			SizeID:       temp.SizeID,
			ColorID:      temp.ColorID,
			Name:         temp.Name,
			Brand:        temp.Brand,
			Subtitle:     temp.Subtitle,
			Description:  temp.Description,
			UnitPrice:    temp.UnitPrice,
			Status:       temp.Status,
			Stock:        temp.Stock,
			Sku:          temp.Sku,
			Weight:       temp.Weight,
			ProductImage: imageProduct,
			Created_at:   temp.Created_at,
			Update_at:    temp.Update_at,
		})
	}

	return data, nil
}

func (repo *repository) ShowProduct(id int) (model.Product, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	// var data = model.Product{}

	query := `select  id, store_id, category_id, size_id, color_id , name,brand, subtitle , description , unit_price , status , stock ,sku , weight , created_at , updated_at from products p where id = $1`

	result, err := repo.db.QueryContext(ctx, query, id)
	failerror.FailError(err, "error prepare")

	var temp = model.Product{}
	for result.Next() {
		result.Scan(&temp.Id, &temp.StoreID, &temp.CategoryID, &temp.SizeID, &temp.ColorID, &temp.Name, &temp.Brand, &temp.Subtitle, &temp.Description, &temp.UnitPrice, &temp.Status, &temp.Stock, &temp.Sku, &temp.Weight, &temp.Created_at, &temp.Update_at)
	}

	if temp.Id <= 0 {
		return temp, errors.New("product not found")
	}

	queryImage := `select i.name, i.image_url from image i join image_products ip on ip.image_id = i.id join products p on p.id  = ip.product_id  where ip.product_id  = $1`

	rows, err := repo.db.QueryContext(ctx, queryImage, temp.Id)
	failerror.FailError(err, "error query")

	var imageProduct = []model.ProductImage{}
	for rows.Next() {
		var images model.ProductImage
		err := rows.Scan(&images.Name, &images.ImageURL)
		failerror.FailError(err, "error scan")
		imageProduct = append(imageProduct, images)
	}

	data := model.Product{
		Id:           temp.Id,
		StoreID:      temp.StoreID,
		CategoryID:   temp.CategoryID,
		SizeID:       temp.SizeID,
		ColorID:      temp.ColorID,
		Name:         temp.Name,
		Brand:        temp.Brand,
		Subtitle:     temp.Subtitle,
		Description:  temp.Description,
		UnitPrice:    temp.UnitPrice,
		Status:       temp.Status,
		Stock:        temp.Stock,
		Sku:          temp.Sku,
		Weight:       temp.Weight,
		ProductImage: imageProduct,
		Created_at:   temp.Created_at,
		Update_at:    temp.Update_at,
	}

	return data, nil
}

func (repo *repository) CreateProduct(req []model.Product) ([]model.Product, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	err := repo.sent.Public(req, "create_product")
	if err != nil {
		return nil, errors.New("failed publisher")
	}

	time.Sleep(3 * time.Second)

	var result []model.Product

	for _, v := range req {

		queryImage := `select i.name, i.image_url from image i join image_products ip on ip.image_id = i.id join products p on p.id  = ip.product_id  where sku = $1`

		rows, err := repo.db.QueryContext(ctx, queryImage, v.Sku)
		failerror.FailError(err, "error query")

		var imageProduct = []model.ProductImage{}
		for rows.Next() {
			var images model.ProductImage
			err := rows.Scan(&images.Name, &images.ImageURL)
			failerror.FailError(err, "error scan")
			imageProduct = append(imageProduct, images)
		}

		queryProduct := `select  id, store_id, category_id, size_id, color_id , name, brand,subtitle , description , unit_price , status , stock ,sku , weight , created_at , updated_at from products p where sku = $1`

		rowsProduct, err := repo.db.QueryContext(ctx, queryProduct, v.Sku)
		failerror.FailError(err, "error prepare")

		var temp model.ProductResult
		for rowsProduct.Next() {
			rowsProduct.Scan(&temp.Id, &temp.StoreID, &temp.CategoryID, &temp.SizeID, &temp.ColorID, &temp.Name, &temp.Brand, &temp.Subtitle, &temp.Description, &temp.UnitPrice, &temp.Status, &temp.Stock, &temp.Sku, &temp.Weight, &temp.Created_at, &temp.Update_at)
		}
		if temp.Id == 0 {
			continue
		}

		result = append(result, model.Product{
			Id:           temp.Id,
			StoreID:      temp.StoreID,
			CategoryID:   temp.CategoryID,
			SizeID:       temp.SizeID,
			ColorID:      temp.ColorID,
			Name:         temp.Name,
			Brand:        temp.Brand,
			Subtitle:     temp.Subtitle,
			Description:  temp.Description,
			UnitPrice:    temp.UnitPrice,
			Status:       temp.Status,
			Stock:        temp.Stock,
			Sku:          temp.Sku,
			Weight:       temp.Weight,
			ProductImage: imageProduct,
			Created_at:   temp.Created_at,
			Update_at:    temp.Update_at,
		})
	}

	if result == nil {
		return []model.Product{}, errors.New("error create product")
	}

	return result, nil
}

func (repo *repository) UpdateProduct(req model.ProductUpd) (model.Product, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	err := repo.sent.Public(req, "update_product")
	if err != nil {
		return model.Product{}, errors.New("failed publisher")
	}

	log.Println(req)
	time.Sleep(1 * time.Second)

	query := `select id, store_id, category_id, size_id, color_id , name, brand,subtitle , description , unit_price , status , stock ,sku , weight , created_at , updated_at from products p where id = $1`

	result, err := repo.db.QueryContext(ctx, query, req.Id)
	failerror.FailError(err, "error prepare")

	var temp model.Product
	for result.Next() {
		result.Scan(&temp.Id, &temp.StoreID, &temp.CategoryID, &temp.SizeID, &temp.ColorID, &temp.Name, &temp.Brand, &temp.Subtitle, &temp.Description, &temp.UnitPrice, &temp.Status, &temp.Stock, &temp.Sku, &temp.Weight, &temp.Created_at, &temp.Update_at)
	}

	if temp.Id == 0 {
		return model.Product{}, errors.New("error create product")
	}

	return temp, nil
}

func (repo *repository) DeleteProduct(id int) (int, error) {
	ctx, cancel := timeout.NewCtxTimeout()
	defer cancel()

	var idCheck int
	queryCheck := `select id from products where id = $1`
	err := repo.db.QueryRowContext(ctx, queryCheck, id).Scan(&idCheck)
	failerror.FailError(err, "error exec")

	if idCheck == 0 {
		return 0, errors.New("product tidak ditemukan")
	}

	var arrImg []int
	queryImage := `select image_id from image_products where product_id = $1`
	rows, err := repo.db.QueryContext(ctx, queryImage, idCheck)
	failerror.FailError(err, "error query")

	for rows.Next() {
		var imageID int
		err := rows.Scan(&imageID)
		failerror.FailError(err, "error scan")
		arrImg = append(arrImg, imageID)
	}

	queryMapping := `delete from image_products where product_id = $1`
	_, err = repo.db.ExecContext(ctx, queryMapping, idCheck)
	failerror.FailError(err, "")

	query := `DELETE FROM products WHERE id = $1`
	_, err = repo.db.ExecContext(ctx, query, idCheck)
	failerror.FailError(err, "")

	for _, v := range arrImg {
		if v != 0 {
			queryDeleteIMG := `delete from image where id = $1`
			_, err := repo.db.ExecContext(ctx, queryDeleteIMG, v)
			failerror.FailError(err, "")
		}
		continue
	}

	return idCheck, nil
}
