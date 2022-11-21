# ya-gophkeeper


Часть тулинга размещена в отдельной директории `tools` и вендорится с помощью `go mod`:

```sh
make bootstrap-deps
```

## Команды sql-migrate

### Добавление новой миграции

`sql-migrate new -env=local migration-name`

### Применение миграций

`sql-migrate up -env=local`