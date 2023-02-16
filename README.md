# Digiutilsapi
>
> A simple RESTFUL API that provides endpoints to manage football leagues.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)

## Prerequisites

- `Go 1.19.3` or higher.
- [Swaggo](https://github.com/swaggo/swag).
- [Migrate](https://github.com/golang-migrate/migrate).

## Installation

- Clone this repository.
- Rename [.env.example](.env.example) file to `.env` and change the `variable value` with your desired value.
- Next, run below command to migrate the database.

```shell
migrate -database "mysql://root:{YOUR_MYSQL_PASSWORD}@tcp(localhost:{YOUR_PORT})/{YOUR_DB_NAME}" -path db/migrations up 
```

> **_NOTE:_**  Replace `{YOUR_MYSQL_PASSWORD}, {YOUR_PORT}, {YOUR_DB_NAME}` with the corresponding value in your `.env` file.

- Then run this command to synchronize all the dependencies.

```shell
go mod tidy
```

- After that run the command below to generate the required files for swaggo.

```shell
swag init
```

- Finally, run the project.

```shell
go run .
```

## Usage

Open your browser and go to `http://localhost:3000/swagger/index.html`.
> **_NOTE:_**  Replace the 3000 with your defined `PORT` value in your `.env` file.

## Meta

Hanif Naufal – [@Digisata](https://twitter.com/Digisata) – [hnaufal123@gmail.com](mailto:hnaufal123@gmail.com)

Distributed under the MIT license. See [LICENSE](LICENSE.md) for more information.

## Contributing

1. Fork this repository.
2. Create your own branch (`git checkout -b fooBar`).
3. Commit your changes (`git commit -am 'Add some fooBar'`).
4. Push to the branch (`git push origin fooBar`).
5. Create your awesome Pull Request.
