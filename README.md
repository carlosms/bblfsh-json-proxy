# bblfsh-json-proxy

A personal project to proxy calls to a gRPC [Babelfish server](https://github.com/bblfsh/bblfshd) using a JSON REST API.

Based on [gitbase-web](https://github.com/src-d/gitbase-web) and @dennwc's work on https://github.com/bblfsh/bblfshd/pull/199.

You can find its [Docker image in Docker Hub](https://hub.docker.com/r/carlosms/bblfsh-json-proxy/).

## Usage

```
Usage:
  bblfsh-json-proxy [OPTIONS] serve [serve-OPTIONS]

starts serving the application

Help Options:
  -h, --help                                     Show this help message

[serve command options]
          --host=                                IP address to bind the HTTP server (default: 0.0.0.0) [$BBLFSH_JSON_HOST]
          --port=                                Port to bind the HTTP server (default: 8095) [$BBLFSH_JSON_PORT]
          --bblfsh=                              Address where bblfsh server is listening (default: 127.0.0.1:9432) [$BBLFSH_JSON_SERVER_URL]

```

## API

### /version

```bash
$ curl localhost:8095/version
{"build":"","version":"v2.9.2"}
```

## /languages

```bash
$ curl localhost:8095/languages
```
<details>
<summary>Response</summary>

```json
[
   {
      "name" : "Python",
      "version" : "v2.4.0",
      "status" : "beta",
      "language" : "python",
      "features" : [
         "ast",
         "uast",
         "roles"
      ]
   },
   {
      "name" : "Java",
      "language" : "java",
      "features" : [
         "ast",
         "uast",
         "roles"
      ],
      "version" : "v2.3.0",
      "status" : "beta"
   },
   {
      "version" : "v2.3.0",
      "status" : "beta",
      "language" : "javascript",
      "features" : [
         "ast",
         "uast",
         "roles"
      ],
      "name" : "JavaScript"
   },
   {
      "version" : "v2.1.0",
      "status" : "beta",
      "language" : "bash",
      "features" : [
         "ast",
         "uast",
         "roles"
      ],
      "name" : "Bash"
   },
   {
      "status" : "beta",
      "version" : "v2.4.0",
      "language" : "ruby",
      "features" : [
         "ast",
         "uast",
         "roles"
      ],
      "name" : "Ruby"
   },
   {
      "version" : "v2.3.0",
      "status" : "beta",
      "language" : "go",
      "features" : [
         "ast",
         "uast",
         "roles"
      ],
      "name" : "Go"
   },
   {
      "version" : "v2.4.0",
      "status" : "beta",
      "language" : "php",
      "features" : [
         "ast",
         "uast",
         "roles"
      ],
      "name" : "PHP"
   }
]
```

</details>

## /parse

You can set the following options in the `post` json body:

- `language`:
- `filename`:
- `content`:
- `filter`: TODO: not implemented yet
- `mode`: One of `native`, `annotated`, `semantic`. Default is `semantic`.

```bash
$ curl -X POST \
  http://localhost:8095/parse \
  -d '{
    "language": "javascript",
    "content": "console.log(\"hello world\");",
    "filename": "hello.js"
  }'
```

<details>
<summary>Response</summary>

```json
{
   "language" : "javascript",
   "uast" : {
      "comments" : [],
      "program" : {
         "@type" : "Program",
         "body" : [
            {
               "@type" : "ExpressionStatement",
               "expression" : {
                  "callee" : {
                     "computed" : false,
                     "property" : {
                        "Name" : "log",
                        "@pos" : {
                           "start" : {
                              "col" : 9,
                              "@type" : "uast:Position",
                              "line" : 1,
                              "offset" : 8
                           },
                           "end" : {
                              "col" : 12,
                              "@type" : "uast:Position",
                              "line" : 1,
                              "offset" : 11
                           },
                           "@type" : "uast:Positions"
                        },
                        "@type" : "uast:Identifier"
                     },
                     "@type" : "MemberExpression",
                     "object" : {
                        "@pos" : {
                           "start" : {
                              "line" : 1,
                              "offset" : 0,
                              "@type" : "uast:Position",
                              "col" : 1
                           },
                           "end" : {
                              "@type" : "uast:Position",
                              "offset" : 7,
                              "line" : 1,
                              "col" : 8
                           },
                           "@type" : "uast:Positions"
                        },
                        "Name" : "console",
                        "@type" : "uast:Identifier"
                     },
                     "@role" : [
                        "Qualified",
                        "Expression",
                        "Identifier",
                        "Call",
                        "Callee"
                     ],
                     "@pos" : {
                        "start" : {
                           "offset" : 0,
                           "line" : 1,
                           "@type" : "uast:Position",
                           "col" : 1
                        },
                        "@type" : "uast:Positions",
                        "end" : {
                           "col" : 12,
                           "@type" : "uast:Position",
                           "line" : 1,
                           "offset" : 11
                        }
                     }
                  },
                  "@pos" : {
                     "start" : {
                        "line" : 1,
                        "offset" : 0,
                        "@type" : "uast:Position",
                        "col" : 1
                     },
                     "end" : {
                        "col" : 27,
                        "@type" : "uast:Position",
                        "line" : 1,
                        "offset" : 26
                     },
                     "@type" : "uast:Positions"
                  },
                  "@role" : [
                     "Expression",
                     "Call"
                  ],
                  "@type" : "CallExpression",
                  "arguments" : [
                     {
                        "@role" : [
                           "Call",
                           "Argument"
                        ],
                        "@pos" : {
                           "end" : {
                              "@type" : "uast:Position",
                              "line" : 1,
                              "offset" : 25,
                              "col" : 26
                           },
                           "@type" : "uast:Positions",
                           "start" : {
                              "col" : 13,
                              "@type" : "uast:Position",
                              "offset" : 12,
                              "line" : 1
                           }
                        },
                        "Value" : "hello world",
                        "Format" : "",
                        "@type" : "uast:String"
                     }
                  ]
               },
               "@role" : [
                  "Statement"
               ],
               "@pos" : {
                  "@type" : "uast:Positions",
                  "end" : {
                     "line" : 1,
                     "offset" : 27,
                     "@type" : "uast:Position",
                     "col" : 28
                  },
                  "start" : {
                     "col" : 1,
                     "line" : 1,
                     "offset" : 0,
                     "@type" : "uast:Position"
                  }
               }
            }
         ],
         "@pos" : {
            "start" : {
               "col" : 1,
               "@type" : "uast:Position",
               "offset" : 0,
               "line" : 1
            },
            "@type" : "uast:Positions",
            "end" : {
               "@type" : "uast:Position",
               "line" : 1,
               "offset" : 27,
               "col" : 28
            }
         },
         "directives" : [],
         "sourceType" : "module",
         "@role" : [
            "Module"
         ]
      },
      "@type" : "File",
      "@role" : [
         "File"
      ],
      "@pos" : {
         "start" : {
            "@type" : "uast:Position",
            "offset" : 0,
            "line" : 1,
            "col" : 1
         },
         "@type" : "uast:Positions",
         "end" : {
            "col" : 28,
            "@type" : "uast:Position",
            "offset" : 27,
            "line" : 1
         }
      }
   }
}
```

</details>