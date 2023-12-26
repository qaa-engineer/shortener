# Сокращатель ссылок на Go

#### Обновление шаблона

```
git remote add -m main template https://github.com/yandex-praktikum/go-musthave-shortener-tpl.git
```

#### Обновление автотестов

```
git fetch template && git checkout template/main .github
```
затем добавьте полученые изменения в свой репозиторий.

Для успешного запуска автотестов вам необходимо давать вашим веткам названия вида `iter<number>`, где `<number>` -
порядковый номер итерации.

Например в ветке с названием `iter4` запустятся автотесты для итераций с первой по четвертую.

При мерже ветки с итерацией в основную ветку (`main`) будут запускаться все автотесты.


#### Полезное

      go mod init github.com/qaa-engineer/shortener
      go build -o shortener *.go