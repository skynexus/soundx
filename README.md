
# SoundX

This is the soundx demo demonstrating a solution for the Sound recommender API.

Project structure:

```
.
├── Makefile   <-- deploy commands + more
├── README.md  <-- This file
├── api        <-- OpenAPI assets
├── build      <-- Build and deploy assets (not committed)
├── cmd        <-- Entry point for the applications
├── db/sql     <-- Database migration files
├── internal   <-- Private application and library code
└── pkg        <-- Library code available for use by external applications
```

The principle tools used by this project are:

- Go
- Postgresql
- OpenAPI
- Postman
- Docker

## Installation

Main requirements:

- [Go 1.21.0 or higher](https://go.dev/doc/install)
- [Docker Engine](https://docs.docker.com/engine/install/)

### Caveats

Requires macOS with Apple silicon because of the Postman CLI tool used for
the integration tests.

TODO: use the Postman Docker image instead.

## Testing

To run the integration tests:

```sh
make integration
```

This rule will:

1. Install the Postman CLI, if necesary
2. Build the SoundX docker image, if necessary
3. Stop any running SoundX container, if a new image was built
4. Start a Postgresql container, if not already running
5. Start the Soundx container, if not already running
6. Run the Sound recommender tests using the Postman CLI

Note that the database is discarded when the Postgresql container stops.

To run the unit tests:

```sh
make test
```

### Caveats

Unit tests are only available for the data access layer.

TODO: write more tests.

## Usage

To start the SoundX service in the background:

```sh
make docker-run
```

To also use the SoundX CLI tool, make sure to also build it:

```sh
make build
```

Then use the SoundX CLI tool to add sounds, for example:

```sh
$ build/soundx-cli sound add "Purple Rain" \
    --bpm 113 --duration 521 --genre rock,rb,gospel \
    --credit "Artist:Prince,Composer:Prince,Producer:Prince"
              Id: 6
           Title: Purple Rain
             BPM: 113
        Duration: 521
          Genres: [rock rb gospel]
          Credit: [{Prince Artist} {Prince Composer} {Prince Producer}]
       CreatedAt: 2023-11-28 14:46 +00:00
       UpdatedAt: 2023-11-28 14:46 +00:00
```

To add a playlist, you will need to use curl for now:

```sh
curl --location 'http://localhost:8080/playlists' \
--header 'Content-Type: application/json' \
--data '{
    "data":
    [
        {
            "title": "Blue",
            "sounds": ["6", "8"]
        }
    ]
}'
```

The CLI tool can be used to list songs and playlists. You can also use it to
get a list of recommended songs based on a playlist, for example:

```sh
$ build/soundx-cli sound recommend 2
  Id Title
   6 Purple Rain
   2 Let's Go Crazy
```

## Development

### OpenAPI specification

The OpenAPI specification is located in:

- api/openapi.yaml

Note that any change to the specification is detected by most rules. For
example, the `docker-run` command will rebuild the image, and restart the
container, if changes are detected in the `api/openapi.yaml` file.

To just update the generated OpenAPI assets, run:

```sh
make api
```

To browse the API specification, run:

```sh
make api-browser
```

This will start a [swagger-ui](https://swagger.io/tools/swagger-ui/) container
and open OpenAPI specification, in your web browser.

The original OpenAPI specification was ported from the Postman collection
and manually adjusted to improve the available schemas, and more.

Note that an additional endpoint was added to the specification:

```yaml
  /sounds/{id}:
    get:
      summary: Get sound
      description: Fetches a specific sound.
```

### Working with the database

To start the postgresql container:

```sh
make db
```

To stop the container, pass the `stop=1` option. Note that the database
is automatically started by the `docker-run` rule. Also note that the
database is discarded when the container stops.

To access the database console:

```sh
make psql
```

Note that this will require `psql` installed on your host.

To create a new migration, use the `migration` rule. For example:

```sh
make migration name=add_awesome_table
```

Note that this will automatically install the
[migration](https://github.com/golang-migrate/migrate) tool.

To run migrations, run:

```sh
make migration
```

Note that the SoundX service automatically runs all database migrations
(that have not been previously run) when started.

### Postman collection

To download a new version of the Postman collection, you will need to
export the collection ID and Postman API key. For example:

```sh
export collection=31378900-dc378d30-94aa-44c7-b1a3-b999e7d5615f
export key=PMAK-65639dd8a349d30030d3c82f-f9f5456035762deb7c4fa5495dbd1a895d
```

Then, if you update the integration tests in the Postman collection, you
can replace the new version like so:

```sh
rm api/postman.json
make api/postman.json
```
