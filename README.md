# FitSync

Thanks for taking a look at FitSync. Below are some helpful terminal commands, PGAdmin login/navigation notes, and a basic guide for importing & testing my APIs in Postman with the stored collections provided.

## Terminal Commands

*Note: all of these commands are called from the "backend" directory (not root!).*

### For All Users

##### Start Docker container
```bash
./../start-app.sh
```

##### Stop Docker container
```bash
docker-compose down -v
```

##### View real-time Docker container logs
```bash
docker logs -f $(docker ps | grep fitsync-backend | awk '{print $1}')
```

### For contributors (just me, lol)

##### Compile & generate SQLc code
```bash
docker exec -it $(docker ps | grep fitsync-backend | awk '{print $1}') sqlc generate
```


## PGAdmin login credentials/navigation guide
If you prefer to see data visualized in table format in the browser, head to <a href="http://localhost:8083/browser/" target="_blank">http://localhost:8083/browser/</a> after starting up the docker container. 

Log in with the first set of credentials:
- `dev@test.com`
- `123lng@#N5las`

Then in the left sidebar, hit the first dropdown (Servers (1)) and log in with the second set of credentials:
- `user01`
- `user01239nTGN35pio!$`

From there, follow absolute path:
Servers (1) > FitSync DB > Databases (2) > fitsyncdb > Schemas (1) > public > Tables (8)


## API Testing with Postman

To test the APIs, import the Postman collections provided:

1. Open Postman
2. Click "Import" in the top left
3. Navigate to `fitsync/backend/postman-collections`
4. Select all .json files in this directory

The collections contain pre-configured requests for all available API endpoints. Make sure the Docker container is running before testing the APIs. Run these in the order they're provided; the majority of routes depend on data from previous requests in the collection and are protected by auth middleware.