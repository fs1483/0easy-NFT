package orders

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"strings"
	"time"

	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createOrderRequest struct {
	Maker        string `json:"maker" binding:"required"`
	NFTAddress   string `json:"nftAddress" binding:"required"`
	TokenID      string `json:"tokenId" binding:"required"`
	PaymentToken string `json:"paymentToken" binding:"required"`
	Price        string `json:"price" binding:"required"`
	Expiry       int64  `json:"expiry" binding:"required"`
	Nonce        string `json:"nonce" binding:"required"`
	Side         string `json:"side" binding:"required"`
	Signature    string `json:"signature" binding:"required"`
}

type cancelOrderRequest struct {
	Maker     string `json:"maker" binding:"required"`
	Nonce     string `json:"nonce" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type orderResponse struct {
	ID           uint      `json:"id"`
	Maker        string    `json:"maker"`
	NFTAddress   string    `json:"nftAddress"`
	TokenID      string    `json:"tokenId"`
	PaymentToken string    `json:"paymentToken"`
	Price        string    `json:"price"`
	Expiry       time.Time `json:"expiry"`
	Nonce        string    `json:"nonce"`
	Side         string    `json:"side"`
	Status       string    `json:"status"`
	Signature    string    `json:"signature"`
	Hash         string    `json:"hash"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// RegisterRoutes 注册订单相关的HTTP路由端点
func RegisterRoutes(rg *gin.RouterGroup, svc *Service) {
	rg.POST("/orders", svc.createOrder)
	rg.GET("/orders", svc.listOrders)
	rg.POST("/orders/:id/cancel", svc.cancelOrder)
}

func (s *Service) createOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	order, err := s.processCreateOrder(c.Request.Context(), &req)
	if err != nil {
		status := http.StatusInternalServerError
		msg := err.Error()
		if errors.Is(err, ErrInvalidOrderPayload) {
			status = http.StatusBadRequest
		} else if errors.Is(err, ErrSignatureMismatch) {
			status = http.StatusUnauthorized
		} else if errors.Is(err, gorm.ErrDuplicatedKey) {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusCreated, toOrderResponse(order))
}

func (s *Service) listOrders(c *gin.Context) {
	side := c.Query("side")
	collection := c.Query("collection")
	status := c.Query("status") // 新增：支持按状态筛选

	var orders []Order
	var err error

	// 根据 status 参数决定查询逻辑
	if status == "" || status == "active" {
		// 默认或明确要求 active：使用原有逻辑
		orders, err = s.repository.ListActive(c.Request.Context(), side, collection)
	} else {
		// 查询特定状态的订单（filled, cancelled）
		orders, err = s.repository.ListByStatus(c.Request.Context(), status, side, collection)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}

	resp := make([]orderResponse, 0, len(orders))
	for _, ord := range orders {
		resp = append(resp, toOrderResponse(&ord))
	}
	c.JSON(http.StatusOK, gin.H{"orders": resp})
}

func (s *Service) cancelOrder(c *gin.Context) {
	idParam := c.Param("id")
	cancelReq := cancelOrderRequest{}
	if err := c.ShouldBindJSON(&cancelReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	orderID, ok := new(big.Int).SetString(idParam, 10)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	if err := s.processCancelOrder(c.Request.Context(), orderID.Uint64(), &cancelReq); err != nil {
		status := http.StatusInternalServerError
		msg := "failed to cancel order"
		switch {
		case errors.Is(err, ErrInvalidOrderPayload):
			status = http.StatusBadRequest
			msg = err.Error()
		case errors.Is(err, ErrSignatureMismatch):
			status = http.StatusUnauthorized
			msg = err.Error()
		case errors.Is(err, gorm.ErrRecordNotFound):
			status = http.StatusNotFound
			msg = "order not found"
		}
		c.JSON(status, gin.H{"error": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order cancelled"})
}

var (
	ErrInvalidOrderPayload  = errors.New("invalid order payload")
	ErrSignatureMismatch    = errors.New("order signature does not match maker")
	errInvalidSide          = errors.New("side must be either ask or bid")
	errInvalidAddressFormat = errors.New("invalid address format")
)

func (s *Service) processCreateOrder(ctx context.Context, req *createOrderRequest) (*Order, error) {
	makerAddr, err := parseAddress(req.Maker)
	if err != nil {
		return nil, ErrInvalidOrderPayload
	}
	nftAddr, err := parseAddress(req.NFTAddress)
	if err != nil {
		return nil, ErrInvalidOrderPayload
	}
	paymentTokenAddr, err := parseAddress(req.PaymentToken)
	if err != nil {
		return nil, ErrInvalidOrderPayload
	}

	tokenID, ok := new(big.Int).SetString(req.TokenID, 10)
	if !ok {
		return nil, ErrInvalidOrderPayload
	}
	price, ok := new(big.Int).SetString(req.Price, 10)
	if !ok {
		return nil, ErrInvalidOrderPayload
	}
	nonce, ok := new(big.Int).SetString(req.Nonce, 10)
	if !ok {
		return nil, ErrInvalidOrderPayload
	}

	if price.Sign() <= 0 {
		return nil, ErrInvalidOrderPayload
	}

	expiry := time.Unix(req.Expiry, 0)
	if expiry.Before(time.Now()) {
		return nil, ErrInvalidOrderPayload
	}

	sideValue, err := parseSide(req.Side)
	if err != nil {
		return nil, ErrInvalidOrderPayload
	}

	sigBytes, err := decodeSignature(req.Signature)
	if err != nil {
		return nil, ErrInvalidOrderPayload
	}

	message := apitypes.TypedDataMessage{
		"maker":        strings.ToLower(makerAddr.Hex()),
		"nft":          strings.ToLower(nftAddr.Hex()),
		"tokenId":      tokenID,
		"paymentToken": strings.ToLower(paymentTokenAddr.Hex()),
		"price":        price,
		"expiry":       big.NewInt(req.Expiry),
		"nonce":        nonce,
		"side":         big.NewInt(int64(sideValue)),
	}

	td := s.typedData
	td.Message = message
	digest, _, err := apitypes.TypedDataAndHash(td)
	if err != nil {
		return nil, err
	}

	if !verifySignature(digest[:], sigBytes, makerAddr) {
		return nil, ErrSignatureMismatch
	}

	order := &Order{
		Maker:        strings.ToLower(makerAddr.Hex()),
		NFTAddress:   strings.ToLower(nftAddr.Hex()),
		TokenID:      req.TokenID,
		PaymentToken: strings.ToLower(paymentTokenAddr.Hex()),
		Price:        req.Price,
		Expiry:       expiry,
		Nonce:        req.Nonce,
		Side:         strings.ToLower(req.Side),
		Status:       OrderStatusActive,
		Signature:    strings.ToLower(req.Signature),
		Hash:         hexutil.Encode(digest),
	}

	if err := s.repository.Create(ctx, order); err != nil {
		return nil, err
	}

	if err := s.cacheOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *Service) processCancelOrder(ctx context.Context, id uint64, payload *cancelOrderRequest) error {
	ord, err := s.repository.FindByID(ctx, uint(id))
	if err != nil {
		return err
	}

	if ord.Status != OrderStatusActive {
		return ErrInvalidOrderPayload
	}

	makerAddr, err := parseAddress(payload.Maker)
	if err != nil {
		return ErrInvalidOrderPayload
	}

	if !strings.EqualFold(ord.Maker, makerAddr.Hex()) {
		return ErrSignatureMismatch
	}

	nonce, ok := new(big.Int).SetString(payload.Nonce, 10)
	if !ok {
		return ErrInvalidOrderPayload
	}

	sigBytes, err := decodeSignature(payload.Signature)
	if err != nil {
		return ErrInvalidOrderPayload
	}

	message := apitypes.TypedDataMessage{
		"maker": strings.ToLower(makerAddr.Hex()),
		"nonce": nonce,
	}

	// 【企业级修复】：创建新的 TypedData 副本，避免修改共享状态导致并发竞态条件
	// 问题原因：多个并发请求共享同一个 s.typedData，修改 PrimaryType 会相互干扰
	// 解决方案：深拷贝 TypedData 对象，每个请求使用独立的实例
	td := apitypes.TypedData{
		Types:       s.typedData.Types,  // Types 是只读的，可以共享
		PrimaryType: "Cancel",           // 为取消请求设置正确的类型
		Domain:      s.typedData.Domain, // Domain 也是只读的
		Message:     message,            // 每个请求的消息都是独立的
	}

	digest, _, err := apitypes.TypedDataAndHash(td)
	if err != nil {
		return err
	}

	if !verifySignature(digest, sigBytes, makerAddr) {
		return ErrSignatureMismatch
	}

	if err := s.repository.UpdateStatus(ctx, ord.ID, OrderStatusCancelled); err != nil {
		return err
	}

	redisKey := "orders:active:" + ord.Side
	if err := s.redisClient.HDel(ctx, redisKey, ord.Hash).Err(); err != nil {
		return err
	}

	cancelledPayload, err := json.Marshal(struct {
		OrderID uint      `json:"orderId"`
		Maker   string    `json:"maker"`
		Nonce   string    `json:"nonce"`
		Hash    string    `json:"hash"`
		Time    time.Time `json:"time"`
	}{OrderID: ord.ID, Maker: ord.Maker, Nonce: ord.Nonce, Hash: ord.Hash, Time: time.Now()})
	if err == nil {
		_ = s.redisClient.Publish(ctx, "orders:cancelled", cancelledPayload).Err()
	}

	return nil
}

func parseAddress(addr string) (common.Address, error) {
	if !common.IsHexAddress(addr) {
		return common.Address{}, errInvalidAddressFormat
	}
	return common.HexToAddress(addr), nil
}

func parseSide(side string) (uint8, error) {
	switch strings.ToLower(side) {
	case "ask":
		return 0, nil
	case "bid":
		return 1, nil
	default:
		return 0, errInvalidSide
	}
}

func decodeSignature(sig string) ([]byte, error) {
	if !strings.HasPrefix(sig, "0x") {
		sig = "0x" + sig
	}
	bytes, err := hexutil.Decode(sig)
	if err != nil {
		return nil, err
	}
	if len(bytes) != 65 {
		return nil, ErrInvalidOrderPayload
	}
	return bytes, nil
}

func verifySignature(digest []byte, sig []byte, expected common.Address) bool {
	sigCopy := make([]byte, len(sig))
	copy(sigCopy, sig)

	if sigCopy[64] >= 27 {
		sigCopy[64] -= 27
	}

	if sigCopy[64] > 1 {
		return false
	}

	pub, err := crypto.SigToPub(digest, sigCopy)
	if err != nil {
		return false
	}
	recovered := crypto.PubkeyToAddress(*pub)
	return strings.EqualFold(recovered.Hex(), expected.Hex())
}

func toOrderResponse(ord *Order) orderResponse {
	return orderResponse{
		ID:           ord.ID,
		Maker:        ord.Maker,
		NFTAddress:   ord.NFTAddress,
		TokenID:      ord.TokenID,
		PaymentToken: ord.PaymentToken,
		Price:        ord.Price,
		Expiry:       ord.Expiry,
		Nonce:        ord.Nonce,
		Side:         ord.Side,
		Status:       string(ord.Status),
		Signature:    ord.Signature,
		Hash:         ord.Hash,
		CreatedAt:    ord.CreatedAt,
		UpdatedAt:    ord.UpdatedAt,
	}
}
