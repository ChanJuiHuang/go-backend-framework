package booter

import (
	"flag"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/ChanJuiHuang/go-backend-framework/pkg/booter/config"
	"github.com/spf13/viper"
)

type Config struct {
	RootDir        string
	ConfigFileName string
	Debug          bool
	Testing        bool
}

func NewConfig(rootDir string, configFileName string, debug bool, testing bool) *Config {
	return &Config{
		RootDir:        rootDir,
		ConfigFileName: configFileName,
		Debug:          debug,
		Testing:        testing,
	}
}

func NewConfigWithCommand() *Config {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var rootDir string
	var configFileName string
	var debug bool
	var testing bool
	flag.StringVar(&rootDir, "rootDir", wd, "root directory which the executable file in")
	flag.StringVar(&configFileName, "configFileName", "config.yml", "config file name")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.BoolVar(&testing, "testing", false, "does run in testing mode")
	flag.Parse()

	return NewConfig(rootDir, configFileName, debug, testing)
}

func NewDefaultConfig() *Config {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return NewConfig(wd, "config.yml", false, false)
}

func NewProductionConfig() *Config {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return NewConfig(wd, "config.production.yml", false, false)
}

func NewTestingConfig() *Config {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime caller cannot get file information")
	}
	wd := path.Join(path.Dir(file), "../..")

	return NewConfig(wd, "config.testing.yml", false, true)
}

type NewConfigFunc func() *Config
type LoadEnvFunc func()

type Registrar interface {
	Boot()
	Register()
}

type RegistrarCenter struct {
	registrars []Registrar
}

func NewRegistrarCenter(registrars []Registrar) *RegistrarCenter {
	return &RegistrarCenter{
		registrars: registrars,
	}
}

func (r *RegistrarCenter) GetRegistrars() []Registrar {
	return r.registrars
}

func (r *RegistrarCenter) Execute() {
	for _, registrar := range r.registrars {
		registrar.Boot()
		registrar.Register()
	}
}

type RegisterExecutor interface {
	BeforeExecute()
	Execute()
	AfterExecute()
}

func bootConfigRegistry(booterConfig *Config) {
	config.Registry.Set("booter", booterConfig)
	byteYaml, err := os.ReadFile(path.Join(booterConfig.RootDir, booterConfig.ConfigFileName))
	if err != nil {
		panic(err)
	}
	stringYaml := os.ExpandEnv(string(byteYaml))

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(strings.NewReader(stringYaml)); err != nil {
		panic(err)
	}

	config.Registry.SetViper(v)
}

func Boot(
	loadEnvFunc LoadEnvFunc,
	newConfigFunc NewConfigFunc,
	registrarCenter RegisterExecutor,
) {
	loadEnvFunc()
	bootConfigRegistry(newConfigFunc())
	registrarCenter.BeforeExecute()
	registrarCenter.Execute()
	registrarCenter.AfterExecute()
}
