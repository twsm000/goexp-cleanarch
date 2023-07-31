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

- Execute o comando abaixo para buildar e executar o servidor.

```bash
go build -o server ./cmd/ordersystem && ./server
```

### Execução via gRPC com Evans

#### Instalação do Evans utilizando Go

```bash
go install github.com/ktr0731/evans@latest
```

### Acessando os serviço de ordens via Evans

- Com o serviço rodando abra um novo terminal e acesse o Evans pelo comando:

```bash
evans -r repl
```

- Conecte-se ao package:

```bash
package pb
```

- Conecte-se ao serviço de ordens:

```bash
service OrderService
```

#### Criando uma nova ordem

```bash
call CreateOrder
```

#### Listando as ordens

```bash
call ListOrders
```
