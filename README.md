# fitsync

### Helpful terminal commands - these are called from the backend directory (not root!)

Docker container up:
`./../start-app.sh`

Docker container down:
`docker-compose down -v`

Real-time Docker container logs:
`docker logs -f $(docker ps | grep fitsync-backend | awk '{print $1}')`

Compile & generate SQLc code:
`docker exec -it $(docker ps | grep fitsync-backend | awk '{print $1}') sqlc generate`

PGAdmin login page credentials:
- dev@test.com
- 123lng@#N5las

PGAdmin view table credentials:
- user01
- user01239nTGN35pio!$