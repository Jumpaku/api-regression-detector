# api-regression-detector

The api-regression-detector is a collection of command-line tools intended to be used for API regression testing. It provides functionalities to initialize tables in a database, perform API call (sending requests and receiving responses), dump modified tables in the database, and compare expected results in the JSON format.

Table of contents:
1. [Install](#install)
2. [Usage](#usage)
     1. [db-init](#db-init)
     1. [db-dump](#db-dump)
     1. [compare](#compare)
     1. [call-http](#call-http)
     1. [call-grpc](#call-grpc)

## Install

```sh
go install github.com/Jumpaku/api-regression-detector/cmd/call-http@latest
go install github.com/Jumpaku/api-regression-detector/cmd/call-grpc@latest
go install github.com/Jumpaku/api-regression-detector/cmd/compare@latest
go install github.com/Jumpaku/api-regression-detector/cmd/db-init@latest
go install github.com/Jumpaku/api-regression-detector/cmd/db-dump@latest
```

## Usage

### db-init

The `db-init` initializes tables in the database using JSON data.

```sh
db-init <database-driver> <connection-string>
db-init -h | --help
```

* `<database-driver>`: Supported database driver name, which can be one of mysql, spanner, sqlite3, or postgres.
* `<connection-string>`: Connection string corresponding to the database driver.
* `-h` `--help`: Show help.

#### Input

To initialize tables in the database, `db-init` expects a JSON data to be provided via stdin. The JSON data should be represented as the following type `DBInitInput`:

```ts
type DBInitInput = InitTable[];
type InitTable = {
        /** name of table to be initialized */
        name: string,
        /** rows to be inserted */
        rows: Row[] 
}
type Row = { [columnName: string]: ColumnValue };
type ColumnValue = boolean | string | number | null
```

Example JSON data:

```json
[
    {
        "name": "example_table",
        "rows": [
            {
                "c0": "abc",
                "c1": 123,
                "c2": true,
                "c3": "2022-12-25T00:45:17Z",
                "id": 1
            },
            {
                "c0": "",
                "c1": 0,
                "c2": false,
                "c3": "2022-12-24T15:54:17Z",
                "id": 2
            }
        ]
    },
    {
        "name": "child_example_table_1",
        "rows": [
            {
                "id": 1,
                "example_table_id": 1
            }
        ]
    },
    {
        "name": "child_example_table_2",
        "rows": [
            {
                "id": 2,
                "example_table_id": 2
            }
        ]
    }
]
```

#### Output

Nothing.

### db-dump

The `db-dump` outputs data within tables in JSON format.

```sh
db-dump <database-driver> <connection-string>
db-dump -h | --help
```

* `<database-driver>`: Supported database driver name, which can be one of mysql, spanner, sqlite3, or postgres.
* `<connection-string>`: Connection string corresponding to the database driver.
* `-h` `--help`: Show help.

#### Input

To output tables in the database, `db-dump` expects a JSON data to be provided via stdin. The JSON data should be represented as the following type `DBDumpInput`:

```ts
/** Array of table names to be dumped */
type DBDumpInput = string[];
```

Example JSON data:

```json
[
    "example_table",
    "child_example_table_1",
    "child_example_table_2"
]
```

#### Output

The `db-dump` tool outputs the data of rows in the specified tables to stdout as JSON data. The JSON data is represented as the following type `DBDumpOutput`:

```ts
type DBDumpOutput = { [tableName: string]: Row[] };
type Row = { [columnName: string]: ColumnValue };
type ColumnValue = boolean | string | number | null
```

Example:
```json
{
    "example_table": [
        {
            "c0": "abc",
            "c1": 123,
            "c2": true,
            "c3": "2022-12-25T00:45:17Z",
            "id": 1
        },
        {
            "c0": "",
            "c1": 0,
            "c2": false,
            "c3": "2022-12-24T15:54:17Z",
            "id": 2
        }
    ],
    "child_example_table_1": [
        {
            "example_table_id": 1,
            "id": 1
        }
    ],
    "child_example_table_2": [
        {
            "example_table_id": 2,
            "id": 2
        }
    ]
}
```

### compare

`compare` compares two JSON files.

```sh
compare [--show-diff] [--no-superset] <expected-json> <actual-json>
compare -h | --help
```

* `<expected-json>`: JSON file path of the expected data.
* `<actual-json>`: JSON file path of the actual data.
* `--show-diff`: Show the difference (default: false).
* `--no-superset`: Disallow superset match (default: false).
* `-h` `--help`: Show help.

#### Input

No specific input required.

#### Output

`compare` outputs the comparison result of the specified two JSON files to stdout. The result is represented as the following type `CompareOutput`:

```ts
type CompareOutput = 
        | "FullMatch" /** two JSON data match exactly */
        | "SupersetMatch" /** second JSON data is an extension of first JSON data */
        | "NoMatch" /** two JSON data are incompatible exactly */;
```

With `--show-diff` option, `compare` outputs difference of specified two JSON data after the result as follows:

```
NoMatch
 |[
 |    {
 |        "add": {
 |            ...skipped 2 object properties...,
+|            "c": 3,
+|            "d": 4
 |        }
 |    },
 |    {
 |        "remove": {
-|            "a": 1
 |        }
 |    },
 |    {
 |        "change": {
~|            "a": 1 => 2
 |        }
 |    }
 |]
```

expected.json
```json
[
    {
        "add": {
            "a": 1,
            "b": 2
        }
    },
    {
        "remove": {
            "a": 1
        }
    },
    {
        "change": {
            "a": 1
        }
    }
]
```

actual.json
```json
[
    {
        "add": {
            "d": 4,
            "a": 1,
            "b": 2,
            "c": 3
        }
    },
    {
        "remove": {}
    },
    {
        "change": {
            "a": 2
        }
    }
]
```

### call-http

`call-http` calls an HTTP API by sending a JSON request and receiving a JSON response.

```
call-http <endpoint-url> <http-method>
call-http -h | --help
```

* `<endpoint-url>`: The URL of the HTTP endpoint, which may have path parameters enclosed in '[' and ']'.
* `<http-method>`: One of GET, HEAD, POST, PUT, DELETE, CONNECT, OPTIONS, TRACE, or PATCH.
* `-h` `--help`: Show help.

#### Input

`call-http` sends a API request provided via stdin in JSON format.

#### Output

`call-http` receives the API response and outputs the response body to stdout in JSON format.

### call-grpc

`call-grpc` calls a GRPC API by sending a JSON request and receiving a JSON response.

```sh
call-grpc <grpc-endpoint> <grpc-full-method>
call-grpc -h | --help
```

* `<grpc-endpoint>`: Host and port joined by ':'.
* `<grpc-full-method>`: Full method in the form 'package.name.ServiceName/MethodName'.
* `-h` `--help`: Show help screen.

#### Input

`call-grpc` sends a API request provided via stdin in JSON format.

#### Output

`call-grpc` receives the API response and outputs the response body to stdout in JSON format.
