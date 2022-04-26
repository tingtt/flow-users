# flow-users

## Usage

### With `docker-compose`

#### Variables `.env`

| Name                    | Description                 | Required           |
| ----------------------- | --------------------------- | ------------------ |
| `PORT`                  | Published port              |                    |
| `MYSQL_DATABASE`        | MySQL database name         |                    |
| `MYSQL_USER`            | MySQL user name             | :heavy_check_mark: |
| `MYSQL_PASSWORD`        | MySQL password              | :heavy_check_mark: |
| `MYSQL_ROOT_PASSWORD`   | MySQL root user password    |                    |
| `LOG_LEVEL`             | API log level               |                    |
| `GZIP_LEVEL`            | API Gzip level              |                    |
| `MYSQL_HOST`            | MySQL host                  | :heavy_check_mark: |
| `MYSQL_PORT`            | MySQL port                  | :heavy_check_mark: |
| `JWT_ISSUER`            | JWT issuer                  | :heavy_check_mark: |
| `JWT_SECRET`            | JWT secret                  | :heavy_check_mark: |
| `GITHUB_CLIENT_ID`      | GitHub OAuth client id      |                    |
| `GITHUB_CLIENT_SECRET`  | GitHub OAuth client secret  |                    |
| `GOOGLE_CLIENT_ID`      | Google OAuth client id      |                    |
| `GOOGLE_CLIENT_SECRET`  | Google OAuth client secret  |                    |
| `TWITTER_CLIENT_ID`     | Twitter OAuth client id     |                    |
| `TWITTER_CLIENT_SECRET` | Twitter OAuth client secret |                    |

```bash
$ docker-compose up
```