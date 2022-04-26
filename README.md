# flow-users

## Usage

### With `docker-compose`

#### Variables `.env`

| Name                    | Description                 | Default   | Required           |
| ----------------------- | --------------------------- | --------- | ------------------ |
| `PORT`                  | Published port              | 1323      |                    |
| `MYSQL_DATABASE`        | MySQL database name         | flow-user |                    |
| `MYSQL_USER`            | MySQL user name             | flow-user |                    |
| `MYSQL_PASSWORD`        | MySQL password              |           | :heavy_check_mark: |
| `MYSQL_ROOT_PASSWORD`   | MySQL root user password    |           |                    |
| `LOG_LEVEL`             | API log level               | 2         |                    |
| `GZIP_LEVEL`            | API Gzip level              | 6         |                    |
| `MYSQL_HOST`            | MySQL host                  | db        |                    |
| `MYSQL_PORT`            | MySQL port                  | 3306      |                    |
| `JWT_ISSUER`            | JWT issuer                  | flow-user |                    |
| `JWT_SECRET`            | JWT secret                  |           | :heavy_check_mark: |
| `GITHUB_CLIENT_ID`      | GitHub OAuth client id      |           |                    |
| `GITHUB_CLIENT_SECRET`  | GitHub OAuth client secret  |           |                    |
| `GOOGLE_CLIENT_ID`      | Google OAuth client id      |           |                    |
| `GOOGLE_CLIENT_SECRET`  | Google OAuth client secret  |           |                    |
| `TWITTER_CLIENT_ID`     | Twitter OAuth client id     |           |                    |
| `TWITTER_CLIENT_SECRET` | Twitter OAuth client secret |           |                    |

```bash
$ docker-compose up
```