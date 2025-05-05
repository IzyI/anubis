# Настройки
```yml

list_services:
  chanki:  # Домен  chanki из (auth.chanki.com)
    service: chanki # отображение в бд 
    role:  # роли пишутся буквами в бд
      owner: "O"
      admin: "A"
      user: "U"
      guest: "G"
    auth: # разрешенные авторизации для сервиса
      - 'phone'
      - 'email'
    webhook:  
      token_revoked: '_-' #куда отправлять инфу что токен дропнулся
  calcer:
    service: chankidssd
  vector:
    service: chankidsds
name_app: # удалить
  anubis
short_jwt:  # короткий ли jwt 
  true
short_jwt_value:
  eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9  # первая часть jwt 


```