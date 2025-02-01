# Reppy

Thanks for taking a look at Reppy. Below are some helpful terminal commands, a guide for importing & testing my APIs in Postman with the provided collections, and PGAdmin login/navigation notes for those who prefer to view database data in table format.

## Terminal Commands

make setup     # Download images and start fresh build
make start     # Start the app (using start-app.sh)
make stop      # Stop containers and remove volumes
make logs      # View backend logs
make sqlc      # Run SQLc code generation
make prefetch  # Just download the Docker images
make clean     # Stop containers and clean up Docker system

## API Testing with Postman

To test the APIs, import the Postman collections provided:

1. Open Postman
2. Click "Import" in the top left
3. Navigate to `reppy/backend/postman-collections`
4. Select all .json files in this directory

The collections contain pre-configured requests for all available API endpoints. Make sure the Docker container is running before testing the APIs. Run these in the order they're provided; the majority of routes depend on data from previous requests in the collection and are protected by auth middleware.


## PGAdmin Login Credentials/Navigation Guide
If you prefer to see data visualized in table format in the browser, head to <a href="http://localhost:8083/browser/" target="_blank">http://localhost:8083/browser/</a> after starting up the docker container. 

Log in with the first set of credentials:
- `dev@test.com`
- `123lng@#N5las`

Then in the left sidebar, hit the first dropdown (Servers (1)) and log in with the second set of credentials:
- `user01`
- `user01239nTGN35pio!$`

From there, follow absolute path:

Servers (1) > Reppy DB > Databases (2) > reppydb > Schemas (1) > public > Tables (8)

