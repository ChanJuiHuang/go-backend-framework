package provider

import (
	"github.com/ChanJuiHuang/go-backend-framework/pkg/authentication"
	"github.com/casbin/casbin/v2"
	"github.com/go-playground/form/v4"
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type registry struct {
	logger        *zap.Logger
	consoleLogger *zap.Logger
	fileLogger    *zap.Logger
	db            *gorm.DB
	redis         *redis.Client
	authenticator *authentication.Authenticator
	casbin        *casbin.SyncedCachedEnforcer
	validator     *validator.Validate
	formDecoder   *form.Decoder
	modifier      *mold.Transformer
}

var Registry *registry

func init() {
	Registry = NewRegistry()
}

func NewRegistry() *registry {
	return &registry{}
}

func (r *registry) Register(
	logger *zap.Logger,
	consoleLogger *zap.Logger,
	fileLogger *zap.Logger,
	db *gorm.DB,
	redis *redis.Client,
	authenticator *authentication.Authenticator,
	casbin *casbin.SyncedCachedEnforcer,
) {
	r.logger = logger
	r.consoleLogger = consoleLogger
	r.db = db
	r.redis = redis
	r.authenticator = authenticator
	r.casbin = casbin
	r.validator = validator.New()
	r.formDecoder = form.NewDecoder()
	r.modifier = modifiers.New()
}

func (r *registry) Logger() *zap.Logger {
	return r.logger
}

func (r *registry) ConsoleLogger() *zap.Logger {
	return r.consoleLogger
}

func (r *registry) FileLogger() *zap.Logger {
	return r.fileLogger
}

func (r *registry) DB() *gorm.DB {
	return r.db
}

func (r *registry) Redis() *redis.Client {
	return r.redis
}

func (r *registry) Authenticator() *authentication.Authenticator {
	return r.authenticator
}

func (r *registry) Casbin() *casbin.SyncedCachedEnforcer {
	return r.casbin
}

func (r *registry) Validator() *validator.Validate {
	return r.validator
}

func (r *registry) FormDecoder() *form.Decoder {
	return r.formDecoder
}

func (r *registry) Modifier() *mold.Transformer {
	return r.modifier
}
