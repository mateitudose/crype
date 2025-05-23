package service

import (
	"crype/api/gen_jet/crype_db/public/model"
	"crype/server/utils"
	"database/sql"
	"fmt"

	"github.com/google/uuid"

	// Dot import so that the Go code resembles native SQL
	. "crype/api/gen_jet/crype_db/public/table"

	"github.com/go-jet/jet/v2/postgres"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// SavePaymentAddress persists a payment address to the database
func (r *OrderRepository) SavePaymentAddress(address, privateKey string) error {
	stmt := PaymentAddresses.INSERT(PaymentAddresses.Address, PaymentAddresses.PrivateKey).VALUES(address, privateKey)
	_, err := stmt.Exec(r.db)
	if err != nil {
		return fmt.Errorf("failed to save payment address: %w", err)
	}
	return nil
}

// SaveOrder persists an order to the database
func (r *OrderRepository) SaveOrder(id uuid.UUID, amount float64, currency, paymentAddress string, createdAt, expiresAt interface{}) error {
	stmt := Orders.INSERT(Orders.ID, Orders.Amount, Orders.Currency, Orders.PaymentAddress, Orders.CreatedAt, Orders.OrderExpiration).
		VALUES(id, amount, currency, paymentAddress, createdAt, expiresAt)
	_, err := stmt.Exec(r.db)
	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}
	return nil
}

// GetOrderByID retrieves an order from the database by ID
func (r *OrderRepository) GetOrderByID(id uuid.UUID) (*model.Orders, error) {
	stmt := Orders.SELECT(Orders.AllColumns).WHERE(Orders.ID.EQ(postgres.UUID(id)))
	var orders []model.Orders
	err := stmt.Query(r.db, &orders)
	if err != nil {
		return nil, fmt.Errorf("database query error: %w", err)
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("order not found")
	}
	return &orders[0], nil
}

// UpdateOrderStatus updates an order's status in the database
func (r *OrderRepository) UpdateOrderStatus(id uuid.UUID, status model.OrderStatus, txHash *string) error {
	// Convert model.OrderStatus to postgres.StringExpression so that we can use it in the SQL query
	statusExpr, err := utils.ModelToEnumOrderStatus(status)
	if err != nil {
		return fmt.Errorf("failed to convert model.OrderStatus to enum: %w", err)
	}
	// TODO: Is there a better way to do this? Using a second SET overrides the first one...
	var stmt postgres.Statement
	if txHash == nil {
		stmt = Orders.UPDATE(Orders.Status).WHERE(Orders.ID.EQ(postgres.UUID(id))).SET(Orders.Status.SET(statusExpr))
	} else {
		stmt = Orders.UPDATE(Orders.Status, Orders.TxHash).WHERE(Orders.ID.EQ(postgres.UUID(id))).SET(Orders.Status.SET(statusExpr), Orders.TxHash.SET(postgres.String(*txHash)))
	}
	_, err = stmt.Exec(r.db)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	return nil
}
