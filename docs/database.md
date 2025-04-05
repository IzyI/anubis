##Установите goose(занимается миграцией) :

```bash 
   go install github.com/pressly/goose/v3/cmd/goose@latest
```

Добавляем в bashrc

```bash 
   echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
   source ~/.bashrc
```

### Управление миграциями с Goose

Goose — инструмент для управления миграциями базы данных в экосистеме Go.

#### Команды

1. **Создание миграции**

```bash 
  goose --dir migrations create user sql
```

- Создает новый SQL-файл миграции в директории `migrations`.
- `user` — имя миграции.

2. **Применение миграций**

```bash 
     goose --dir migrations up
```

- Применяет все новые миграции из директории `migrations`.

3. **Откат миграций**

```bash 
  goose --dir migrations down-to 0
```

- Откатывает все миграции, возвращая базу данных в начальное состояние.