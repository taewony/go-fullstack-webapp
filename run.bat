@echo off
mkdir go-chichat
cd go-chichat

echo . > main.go
echo . > go.mod
echo . > go.sum
echo . > README.md
echo . > .env

mkdir internal
mkdir internal\handler
mkdir internal\model
mkdir internal\repository
mkdir internal\template

mkdir public

echo 디렉토리 구조 생성이 완료되었습니다.