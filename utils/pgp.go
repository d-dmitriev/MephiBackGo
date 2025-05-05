package utils

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	_ "golang.org/x/crypto/ripemd160"
)

var (
	pgpPrivateKey, pgpPublicKey string
)

const (
	privateKeyFile = "pgp_private.key"
	publicKeyFile  = "pgp_public.key"
)

func init() {
	if fileExists(privateKeyFile) && fileExists(publicKeyFile) {
		// Читаем из файла
		priv, _ := os.ReadFile(privateKeyFile)
		pub, _ := os.ReadFile(publicKeyFile)
		pgpPrivateKey = string(priv)
		pgpPublicKey = string(pub)
	} else {
		// Генерируем новые и сохраняем
		entity, err := generatePGPKey()
		if err != nil {
			panic("failed to generate PGP keys: " + err.Error())
		}
		pgpPublicKey, pgpPrivateKey = entityToArmored(entity)

		err = os.WriteFile(privateKeyFile, []byte(pgpPrivateKey), 0600)
		if err != nil {
			panic("failed to save private key: " + err.Error())
		}

		err = os.WriteFile(publicKeyFile, []byte(pgpPublicKey), 0644)
		if err != nil {
			panic("failed to save public key: " + err.Error())
		}
	}
}

// generatePGPKey — генерирует пару ключей PGP
func generatePGPKey() (*openpgp.Entity, error) {
	cfg := &packet.Config{
		RSABits: 2048,
	}
	// Создаем сущность с указанием размера ключа через RSABitsOption
	entity, err := openpgp.NewEntity(
		"Bank API",
		"",
		"",
		cfg,
	)
	if err != nil {
		return nil, err
	}

	return entity, nil
}

// entityToArmored — преобразует ключи в формат ASCII-armored
func entityToArmored(entity *openpgp.Entity) (pub, priv string) {
	var pubBuf, privBuf bytes.Buffer

	// Публичный ключ
	w, _ := armor.Encode(&pubBuf, openpgp.PublicKeyType, nil)
	entity.Serialize(w)
	w.Close()

	// Приватный ключ
	w2, _ := armor.Encode(&privBuf, openpgp.PrivateKeyType, nil)
	entity.SerializePrivate(w2, nil)
	w2.Close()

	return pubBuf.String(), privBuf.String()
}

// PGPEncrypt — шифрует данные с использованием публичного ключа
func PGPEncrypt(data string) (string, error) {
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(pgpPublicKey))
	if err != nil {
		return "", err
	}

	var encBuf bytes.Buffer
	encodeWriter, err := armor.Encode(&encBuf, "PGP MESSAGE", nil)
	if err != nil {
		return "", err
	}

	encryptWriter, err := openpgp.Encrypt(encodeWriter, entityList, nil, nil, nil)
	if err != nil {
		return "", err
	}

	_, err = io.WriteString(encryptWriter, data)
	if err != nil {
		return "", err
	}

	if err := encryptWriter.Close(); err != nil {
		return "", err
	}

	if err := encodeWriter.Close(); err != nil {
		return "", err
	}

	return encBuf.String(), nil
}

// PGPDecrypt — расшифровывает данные с использованием приватного ключа
func PGPDecrypt(data string) (string, error) {
	entityList, err := openpgp.ReadArmoredKeyRing(strings.NewReader(pgpPrivateKey))
	if err != nil {
		return "", err
	}

	entity := entityList[0]
	if entity.PrivateKey == nil {
		return "", errors.New("private key not available")
	}

	dataBuffer := bytes.NewBuffer([]byte(data))
	armorBlock, err := armor.Decode(dataBuffer)

	md, err := openpgp.ReadMessage(armorBlock.Body, entityList, nil, nil)
	if err != nil {
		return "", err
	}

	plaintext, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Вспомогательная функция
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
