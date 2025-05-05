package services

import (
	"bank-api/repositories"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"bank-api/models"
	"bank-api/utils"
)

type CardService struct {
	cardRepo repositories.CardRepository
}

func NewCardService(cardRepo repositories.CardRepository) *CardService {
	return &CardService{
		cardRepo: cardRepo,
	}
}

// GenerateCardNumber — генерирует валидный номер карты по алгоритму Луна
func (s *CardService) GenerateCardNumber(bin string, length int) (string, error) {
	if len(bin)+1 > length {
		return "", errors.New("invalid BIN or card length")
	}

	// Генерируем промежуточные цифры (всё, кроме последней)
	randomLength := length - len(bin) - 1 // одна цифра — контрольная
	randomDigits := make([]byte, randomLength)
	for i := range randomDigits {
		randomDigits[i] = byte(rand.Intn(10) + '0')
	}

	card := bin + string(randomDigits) + "0"

	// Применяем алгоритм Луна для подсчёта последней цифры
	digits := make([]int, len(card))
	for i := 0; i < len(card); i++ {
		digit, _ := strconv.Atoi(string(card[i]))
		digits[i] = digit
	}

	sum := 0
	double := false

	for i := len(digits) - 1; i >= 0; i-- {
		digit := digits[i]

		if double {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		double = !double
	}

	mod := sum % 10
	checkDigit := 0
	if mod != 0 {
		checkDigit = 10 - mod
	}

	card = card[:len(card)-1] + strconv.Itoa(checkDigit)

	return card, nil
}

// IssueCard — выпускает новую карту для пользователя
func (s *CardService) IssueCard(userID uint, accountID uint) (*models.Card, error) {
	// Генерируем номер карты (пример BIN: шесть цифр бина банка 453275 + две цифры банковская программа 01)
	number, err := s.GenerateCardNumber("45327501", 16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate card number: %v", err)
	}

	// Шифруем номер карты PGP
	encryptedNumber, err := utils.PGPEncrypt(number)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt card number: %v", err)
	}

	// Хешируем CVV
	cvv := fmt.Sprintf("%03d", rand.Intn(999))
	encryptedCVV, err := utils.PGPEncrypt(cvv)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt CVV: %v", err)
	}

	// Генерируем дату истечения (обычно 3 года)
	expireDate := time.Now().AddDate(3, 0, 0).Format("01/06") // MM/YY
	encryptedExpireDate, err := utils.PGPEncrypt(expireDate)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt expireDate: %v", err)
	}

	card := &models.Card{
		UserID:     userID,
		AccountID:  accountID,
		Number:     encryptedNumber,
		CVV:        encryptedCVV,
		ExpireDate: encryptedExpireDate,
		IssuedAt:   time.Now().Format(time.RFC3339),
		Status:     "active",
	}

	s.cardRepo.IssueCard(card)

	// Возвращаем не зашифрованный номер и CVV только что созданной карты
	card.ExpireDate = expireDate
	card.Number = number
	card.CVV = cvv
	return card, nil
}

// GetCards — получает все карты пользователя
func (s *CardService) GetCards(userID string) ([]models.Card, error) {
	userIDUint, err := utils.ParseUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %v", err)
	}

	cards, err := s.cardRepo.GetCards(userIDUint)
	if err != nil {
		return nil, fmt.Errorf("failed to get cards: %v", err)
	}

	// Расшифровываем номера карт перед возвратом клиенту
	for i := range cards {
		if decryptedNumber, err := utils.PGPDecrypt(cards[i].Number); err == nil {
			cards[i].Number = decryptedNumber
		} else {
			utils.Logger.Warnf("Failed to decrypt card number for card ID %d: %v", cards[i].ID, err)
			cards[i].Number = "****************"
		}
		if decryptedExpireDate, err := utils.PGPDecrypt(cards[i].ExpireDate); err == nil {
			cards[i].ExpireDate = decryptedExpireDate
		} else {
			utils.Logger.Warnf("Failed to decrypt expire date for card ID %d: %v", cards[i].ID, err)
			cards[i].ExpireDate = "**/**"
		}
		if decryptedCVV, err := utils.PGPDecrypt(cards[i].CVV); err == nil {
			cards[i].CVV = decryptedCVV
		} else {
			utils.Logger.Warnf("Failed to decrypt cvv for card ID %d: %v", cards[i].ID, err)
			cards[i].CVV = "***"
		}
	}

	return cards, nil
}
