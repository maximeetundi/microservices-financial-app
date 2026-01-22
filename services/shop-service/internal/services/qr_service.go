package services

import (
	"bytes"
	"fmt"

	qrcode "github.com/skip2/go-qrcode"
)

type QRService struct {
	baseURL string
}

func NewQRService(baseURL string) *QRService {
	return &QRService{baseURL: baseURL}
}

// GenerateShopQR generates a QR code for a shop
func (s *QRService) GenerateShopQR(shopSlug string) ([]byte, error) {
	url := fmt.Sprintf("%s/shops/%s", s.baseURL, shopSlug)
	return s.generateQR(url)
}

// GenerateProductQR generates a QR code for a product
func (s *QRService) GenerateProductQR(shopSlug, productSlug string) ([]byte, error) {
	url := fmt.Sprintf("%s/shops/%s/product/%s", s.baseURL, shopSlug, productSlug)
	return s.generateQR(url)
}

// GenerateCategoryQR generates a QR code for a category
func (s *QRService) GenerateCategoryQR(shopSlug, categorySlug string) ([]byte, error) {
	url := fmt.Sprintf("%s/shops/%s/category/%s", s.baseURL, shopSlug, categorySlug)
	return s.generateQR(url)
}

func (s *QRService) generateQR(content string) ([]byte, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}
	
	var buf bytes.Buffer
	buf.Write(png)
	return buf.Bytes(), nil
}
