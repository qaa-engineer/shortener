# Сокращатель ссылок на Go

http://localhost:8080/

#### Обновление шаблона

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

#### Обновление автотестов

```
git fetch template && git checkout template/main .github
```

#### Полезное

        go mod init github.com/qaa-engineer/shortener
        go build -o shortener *.go
