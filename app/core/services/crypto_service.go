//Package services
/*
 * @Time    : 2021年05月06日 09:19:41
 * @Author  : user
 * @Project : SEGI
 * @File    : crypto-services.go
 * @Software: GoLand
 * @Describe:
 */
package services

import (
	"context"
	"fmt"
	"scrm-sync-data/app/config"
	"scrm-sync-data/app/core/services/crypto-services"
	"time"
)

type CryptoService struct {
	context context.Context
	cancel  context.CancelFunc
	config  *config.Config
	logger  *LoggerService
	AES     *crypto.AES
	DES     *crypto.DES
	RSA     *crypto.RSA
	MD5     *crypto.MD5
	Sha1    *crypto.SHA1
	Base64  *crypto.Base64
	State   bool
}

func (cryptoService *CryptoService) Init(parentContext context.Context, config *config.Config, logger *LoggerService) {
	cryptoService.config = config
	cryptoService.logger = logger
	cryptoService.State = false
	cryptoService.context, cryptoService.cancel = context.WithCancel(parentContext)
	cryptoService.Start()
}

func (cryptoService *CryptoService) Start() {
	fmt.Println("start crypto service...", time.Now())
	cryptoService.logger.Debug(fmt.Sprintf("%s,%v", "start crypto service...", time.Now()))
	var err error
	cryptoService.AES = new(crypto.AES)
	cryptoService.DES = new(crypto.DES)
	cryptoService.RSA = new(crypto.RSA)
	cryptoService.RSA, err = cryptoService.RSA.Init()
	if err != nil {
		cryptoService.Stop()
		return
	}
	cryptoService.MD5 = new(crypto.MD5)
	cryptoService.Sha1 = new(crypto.SHA1)
	cryptoService.Base64 = new(crypto.Base64)
	cryptoService.Base64 = cryptoService.Base64.Init()
	cryptoService.State = true
}

func (cryptoService *CryptoService) Stop() {
	fmt.Println("stop crypto service...", time.Now())
	cryptoService.logger.Debug(fmt.Sprintf("%s,%v", "stop crypto service...", time.Now()))
	cryptoService.cancel()
	cryptoService.State = false
}
