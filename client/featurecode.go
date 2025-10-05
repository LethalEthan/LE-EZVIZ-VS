package client

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"

	"go.uber.org/zap"
)

// Load featurecode from file, to note it is a random md5 hash and doesn't have much importance
// it is also known as hardwarecode, its name varies within API usage but is the same
func (LEZ *LE_EZVIZ_Client) LoadFeatureCode(featurecodepath string) string {
	f, err := os.Open(featurecodepath)
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create("featurecode")
			if err != nil {
				log.Error("Error creating featurecode file, creating in memory", zap.Error(err))
			} else {
				hash := GenerateFeatureCode()
				f.Write([]byte(hex.EncodeToString(hash[:])))
				f.Sync()
				f.Close()
				LEZ.FeatureCode = string(hex.EncodeToString(hash[:]))
				LEZ.Headers["featureCode"] = []string{LEZ.FeatureCode}
				return hex.EncodeToString(hash[:])
			}
			hash := GenerateFeatureCode()
			LEZ.FeatureCode = string(hex.EncodeToString(hash[:]))
			LEZ.Headers["featureCode"] = []string{LEZ.FeatureCode}
			return hex.EncodeToString(hash[:])
		}
	}
	featurebytes, err := io.ReadAll(f)
	if err != nil {
		log.Error("Error reading featurecode from file, creating in memory", zap.Error(err))
		hash := GenerateFeatureCode()
		LEZ.FeatureCode = string(hex.EncodeToString(hash[:]))
		LEZ.Headers["featureCode"] = []string{LEZ.FeatureCode}
		return hex.EncodeToString(hash[:])
	}
	if len(featurebytes) != 32 {
		log.Error("featurecode byte length is incorrect, creating in memory", zap.Error(err))
		hash := GenerateFeatureCode()
		LEZ.FeatureCode = string(hex.EncodeToString(hash[:]))
		LEZ.Headers["featureCode"] = []string{LEZ.FeatureCode}
		return hex.EncodeToString(hash[:])
	}
	LEZ.FeatureCode = string(featurebytes)
	LEZ.Headers["featureCode"] = []string{LEZ.FeatureCode}
	return LEZ.FeatureCode
}

func GetMd5(text string) string {
	h := md5.Sum([]byte(text))
	return hex.EncodeToString(h[:])
}

func GenerateFeatureCode() []byte {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)
	hash := md5.Sum(randomBytes)
	return hash[:]
}

func (LEZ *LE_EZVIZ_Client) SetFeatureCode(code string) {
	if len(code) != 32 {
		log.Info("featurecode must be 32 bytes")
		LEZ.FeatureCode = string(GenerateFeatureCode())
	}
	LEZ.FeatureCode = code
	LEZ.Headers["featureCode"] = []string{code}
}
