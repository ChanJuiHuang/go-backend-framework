#!/bin/bash

printf "Enter your project name:"
read PROJECT_NAME
PROJECT_NAME=${PROJECT_NAME:=go-backend}
mkdir $PROJECT_NAME

printf "Enter the framework version what you want(EX: 1.0.0):"
read VERSION

curl -L https://github.com/ChanJuiHuang/go-backend-framework/archive/refs/tags/v$VERSION.tar.gz | tar -zxv -C $PROJECT_NAME --strip-components=1
cd $PROJECT_NAME
cp .env.example .env
cp casbin_model.conf.example casbin_model.conf
cp .air.toml.exmaple .air.toml

OLD_MODULE_NAME=github.com/ChanJuiHuang/go-backend-framework
printf "Enter your go module name (\e[33mex: github.com/author/my-project\e[0m):"
read NEW_MODULE_NAME

go mod edit -module $NEW_MODULE_NAME
case $OSTYPE in
darwin*)  find . -type f -name "*.go" -exec sed -i "" "s,$OLD_MODULE_NAME,$NEW_MODULE_NAME,g" {} \; ;; 
linux*)   find . -type f -name "*.go" -exec sed -i "s,$OLD_MODULE_NAME,$NEW_MODULE_NAME,g" {} \; ;;
esac

printf "Do you want to enable database? \e[33m(y or n, default to n)\e[0m"
read IS_ENABLED
IS_ENABLED=${IS_ENABLED:=n}
if [ $IS_ENABLED == "y" ] ;then
  case $OSTYPE in
    darwin*) sed -i "" "s,DB_ENABLED.*,DB_ENABLED=true,g" .env ;; 
    linux*)  sed -i "s,DB_ENABLED.*,DB_ENABLED=true,g" .env ;;
  esac
else
  case $OSTYPE in
    darwin*) sed -i "" "s,DB_ENABLED.*,DB_ENABLED=false,g" .env ;; 
    linux*)  sed -i "s,DB_ENABLED.*,DB_ENABLED=false,g" .env ;;
  esac
fi

printf "Do you want to enable redis? \e[33m(y or n, default to n)\e[0m"
read IS_ENABLED
IS_ENABLED=${IS_ENABLED:=n}
if [ $IS_ENABLED == "y" ] ;then
  case $OSTYPE in
    darwin*) sed -i "" "s,REDIS_ENABLED.*,REDIS_ENABLED=true,g" .env ;; 
    linux*)  sed -i "s,REDIS_ENABLED.*,REDIS_ENABLED=true,g" .env ;;
  esac
else
  case $OSTYPE in
    darwin*) sed -i "" "s,REDIS_ENABLED.*,REDIS_ENABLED=false,g" .env ;; 
    linux*)  sed -i "s,REDIS_ENABLED.*,REDIS_ENABLED=false,g" .env ;;
  esac
fi

printf "Do you want to enable casbin? \e[33m(y or n, default to n)\e[0m"
read IS_ENABLED
IS_ENABLED=${IS_ENABLED:=n}
if [ $IS_ENABLED == "y" ] ;then
  case $OSTYPE in
    darwin*) sed -i "" "s,CASBIN_ENABLED.*,CASBIN_ENABLED=true,g" .env ;; 
    linux*)  sed -i "s,CASBIN_ENABLED.*,CASBIN_ENABLED=true,g" .env ;;
  esac
else
  case $OSTYPE in
    darwin*) sed -i "" "s,CASBIN_ENABLED.*,CASBIN_ENABLED=false,g" .env ;; 
    linux*)  sed -i "s,CASBIN_ENABLED.*,CASBIN_ENABLED=false,g" .env ;;
  esac
fi
