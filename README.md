# GoExpert - Clean Arch

## Rodando o projeto

- Instale a versão correta do [migrate](https://github.com/golang-migrate/migrate) utilizando uma [release build](https://github.com/golang-migrate/migrate/releases)

- Execute os comandos abaixo:

```bash:
go mod tidy
```

```bash:
docker-compose up -d
```

```bash
migrate -path=internal/infra/database/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" -verbose up
```

- Com esses comandos você irá inicializar as dependências do projeto, executar o docker para rodar o mysql e rabbitmq e inicializar as migrations para que o banco esteja com as tabelas corretas para serem utilizadas pelo projeto.
