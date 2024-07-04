# BCAF REST API

### Endpoints

All endpoints give a JSON response with a `response` field.\
Authorization data which is given by the Discord oauth2 API is passed via an `authorization` header.\
`POST` requests require a JSON body.

| Endpoint | Methods | Required authorization | `POST` body fields | Description |
| - | - | - | - | - |
| `/v1/ping` | `GET` | | | Pings the API |
| `/v1/minecraft/name/:uuid` | `GET` | | | `GET`s the corresponding Minecraft username for a UUID |
| `/v1/accounts/lookup/:id` | `GET`, `POST` | Access token | `color`, `backgroundImageUrl`, `foregroundImageUrl` | `GET`s a users account data with the corresponding Discord user id or `POST`s changes to one's own account data given the changes made are valid
| `/v1/accounts/list` | `GET` | Access token | | `GET`s a list of every registered users account data
| `/v1/accounts/auth` | `GET` | Code | | `GET`s Discord authentication data as given by `https://discord.com/api/oauth2/token`
| `/v1/accounts/refresh` | `GET` | Refresh token | | `GET`s refreshed Discord authentication data |
| `/v1/smp/info` | `GET` | (Optional) Access token | | `GET`s a list of all currently running Minecraft SMPs (IPs are omitted if no access token is given) |
| `/v1/smp/worlds` | `GET` | Access token | | `GET`s a list of all previous Minecraft SMP world downloads |