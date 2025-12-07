package product

import (
	"context"
	"time"

	"github.com/QosmuratSamat0/catalog/internal/domain/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	catalog "github.com/QosmuratSamat0/catalog/proto/gen/go/catalog"
)

const (
	emptyValue = 0
)

type Product interface {
	CreateProduct(ctx context.Context, name, description, category string, price float64) (id int64, err error)
	GetProductById(ctx context.Context, id int64) (product models.Product, err error)
	SearchProducts(ctx context.Context, category string, priceMin, priceMax float64, sortBy string) (products []models.Product, err error)
}

type serverAPI struct {
	catalog.UnimplementedProductServicesServer
	product Product
}

func Register(gRPC *grpc.Server, product Product) {
	catalog.RegisterProductServicesServer(gRPC, &serverAPI{product: product})
}

func (s *serverAPI) CreateProduct(
	ctx context.Context,
	req *catalog.CreateProductRequest,
) (*catalog.CreateProductResponse, error) {

	if err := validateCreateProductRequest(req); err != nil {
		return nil, err
	}

	id, err := s.product.CreateProduct(ctx, req.GetName(), req.GetDescription(), req.GetCategory(), req.GetPrice())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &catalog.CreateProductResponse{Id: id}, nil

}

func (s *serverAPI) GetProduct(
	ctx context.Context,
	req *catalog.GetProductRequest,
) (*catalog.GetProductResponse, error) {
	if err := validateGetProductRequest(req); err != nil {
		return nil, err
	}

	p, err := s.product.GetProductById(ctx, req.GetId())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &catalog.GetProductResponse{
		Product: toProtoProduct(p),
	}, nil
}

func (s *serverAPI) SearchProducts(
	ctx context.Context,
	req *catalog.SearchProductsRequest,
) (*catalog.SearchProductsResponse, error) {
	if err := validateSearchProductRequest(req); err != nil {
		return nil, err
	}

	products, err := s.product.SearchProducts(ctx, req.GetCategory(), req.GetPriceMin(), req.GetPriceMax(), req.GetSortBy())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "Internal error")
	}
	var protoProducts []*catalog.Product
	for _, product := range products {
		protoProducts = append(protoProducts, toProtoProduct(product))
	}

	return &catalog.SearchProductsResponse{Product: protoProducts}, nil
}

func toProtoProduct(m models.Product) *catalog.Product {
	return &catalog.Product{
		Id:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Category:    m.Category,
		Price:       m.Price,
		CreatedAt:   m.CreatedAt.Format(time.RFC3339),
	}
}

func validateSearchProductRequest(req *catalog.SearchProductsRequest) error {

	if req == nil {
		return status.Error(codes.InvalidArgument, "empty product request")
	}
	if req.GetCategory() == "" {
		return status.Error(codes.InvalidArgument, "category is required")
	}
	if req.GetPriceMax() == emptyValue {
		return status.Error(codes.InvalidArgument, "price is 0")
	}
	if req.GetPriceMin() < emptyValue {
		return status.Error(codes.InvalidArgument, "price is less than 0")
	}
	if req.GetSortBy() == "" {
		return status.Error(codes.InvalidArgument, "sort by is required")
	}

	return nil
}

func validateGetProductRequest(req *catalog.GetProductRequest) error {

	if req == nil {
		return status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	if req.GetId() < 0 {
		return status.Error(codes.InvalidArgument, "invalid product id")
	}
	return nil
}

func validateCreateProductRequest(req *catalog.CreateProductRequest) error {
	if req == nil {
		return status.Error(codes.InvalidArgument, "request cannot be empty")
	}
	if req.GetName() == "" {
		return status.Error(codes.InvalidArgument, "name is required")
	}
	if req.GetDescription() == "" {
		return status.Error(codes.InvalidArgument, "description is required")
	}
	if req.GetCategory() == "" {
		return status.Error(codes.InvalidArgument, "category is required")
	}
	if req.GetPrice() < 0 {
		return status.Error(codes.InvalidArgument, "price less than 0")
	}

	return nil
}
