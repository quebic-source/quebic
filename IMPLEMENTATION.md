# Quebic — Implementation & Architecture Analysis

> Generated: 2026-05-29
> Scope: Full repository analysis of the architecture, components, runtime flows,
> and implementation details of the Quebic FaaS framework.

---

## 1. What Quebic Is

**Quebic** is an open-source **Functions-as-a-Service (FaaS) / event-driven
microservices framework** written in **Go**, designed to run serverless functions
on **Kubernetes** (with a Docker-based deployment mode as well). Developers write
functions in **Python (2.7 / 3.6), Java 8, and NodeJS**; Quebic packages each into
a container, deploys it to the cluster, wires it to an **event bus (RabbitMQ /
AMQP)**, and exposes it through an **API Gateway**.

It is a real, working multi-binary system (not a scaffold). The codebase is
≈14k lines of Go across four executables plus shared libraries.

### The four binaries
| Binary | Module dir | Role |
|--------|-----------|------|
| **quebic-faas-mgr** | `quebic-faas-mgr/` | Control plane / manager. Owns the REST management API, persistence, function build & deployment, and lifecycle of platform components. |
| **quebic-faas-apigateway** | `quebic-faas-apigateway/` | Data plane ingress. Maps external HTTP routes to events and fires them onto the event bus (sync or async). |
| **quebic-faas-cli** | `quebic-faas-cli/` | `quebic` command-line client that talks to the manager's REST API. |
| **quebic-faas-mgr-dashboard** | `quebic-faas-mgr-dashboard/` | Serves the management web dashboard (static UI). |

Functions themselves run as separate containers built from the language runtime
base images in `docker-files/` and the manager's `function/` package.

---

## 2. High-Level Architecture

```
                ┌──────────────┐         ┌──────────────────────┐
   quebic CLI ──┤  REST (JWT)  ├────────►│   quebic-faas-mgr     │  CONTROL PLANE
  (or dashboard)└──────────────┘         │  (Manager)            │
                                         │                       │
                                         │  httphandler  ── auth │  JWT auth
                                         │  dao ── db (bolt)     │  persistence
                                         │  function build/image │  build images
                                         │  deployment (docker/  │  deploy
                                         │             kube)     │
                                         │  components mgr       │  bring up infra
                                         │  gc                   │  garbage collect
                                         └──────┬────────────────┘
                       brings up / manages      │
        ┌──────────────────────┬────────────────┼───────────────────┐
        ▼                      ▼                 ▼                   ▼
 ┌────────────┐       ┌────────────────┐  ┌──────────────┐   ┌──────────────┐
 │  eventbus  │       │   apigateway   │  │  mgr-dashboard│   │   eventbox   │
 │ (RabbitMQ) │◄──────┤ (HTTP ingress) │  │  (web UI)     │   │ (req tracker)│
 └─────┬──────┘       └───────┬────────┘  └──────────────┘   └──────────────┘
       │  publish/subscribe   │ external HTTP in
       │  (messenger pkg)     ▼
       │              client request → route → event → function
       ▼
 ┌────────────────────────────────────────────────────┐
 │  Function containers (Python 2.7/3.6, Java 8,       │  DATA PLANE
 │  NodeJS) — each subscribes to its events via the    │
 │  AMQP topic exchange, runs handler, replies via     │
 │  messenger                                          │
 └────────────────────────────────────────────────────┘
```

Diagrams shipped in the repo: `docs/digram-1..4.png`, `docs/quebic.png`,
`docs/quebic.xml` (draw.io source), and `docs/quebic-introduction.pptx`.

---

## 3. Repository Layout

```
.
├── auth/                       JWT auth, roles, token utilities (shared)
├── common/                     Cross-cutting constants & helpers (shared)
├── config/                     Base app config struct + load/save (shared)
├── messenger/                  RabbitMQ/AMQP event bus client (shared)
├── types/                      Domain types + DTOs (shared)
│
├── quebic-faas-mgr/            MANAGER binary (control plane)
│   ├── quebic-faas-mgr.go      main()
│   ├── cmd/                    cobra root + flags
│   ├── config/                 manager config (extends base config)
│   ├── app/                    App bootstrap: wires everything together
│   ├── db/                     BoltDB open (file .quebic-faas-mgr.db)
│   ├── dao/                    data-access: function, user, mgr-component
│   ├── httphandler/            REST endpoints (function/route/event/user/...)
│   ├── components/             brings up eventbus / apigateway / dashboard / eventbox
│   ├── deployment/             deployment abstraction
│   │   ├── docker_deployment/  Docker driver
│   │   └── kube_deployment/    Kubernetes driver (largest file, ~1180 LOC)
│   ├── function/               function build pipeline
│   │   ├── function_create/    create/save function spec
│   │   ├── function_image/     build container image for a function
│   │   ├── function_common/    shared consts + file utils
│   │   ├── function_util/      packaging / artifact handling
│   │   ├── function_life_controller/  container life (awake/idle)
│   │   └── function_runtime/   per-language runtime builders
│   │       ├── function_java_8_runtime/
│   │       ├── function_nodejs_runtime/
│   │       ├── function_python_2_7_runtime/
│   │       └── function_python_3_6_runtime/
│   ├── gc/                     garbage collector for stale resources
│   ├── logger/                 request-scoped logging
│   └── Dockerfile              manager image
│
├── quebic-faas-apigateway/     API GATEWAY binary (data plane ingress)
│   ├── quebic-faas-apigateway.go  main()
│   ├── config/                 gateway config
│   ├── db/                     gateway local store
│   ├── httphandler/            request routing + resource_setup (route table)
│   └── Dockerfile
│
├── quebic-faas-cli/            CLI binary (`quebic`)
│   ├── quebic-faas-cli.go      main()
│   ├── app/                    CLI bootstrap
│   ├── cmd/                    cobra commands (function, route, event, user,
│   │                           mgr, apigateway, ingress, eventbox,
│   │                           request_tracker, config, version, ...)
│   ├── service/                HTTP client services calling the manager API
│   ├── config/                 cli-config.yml (~/.quebic-faas/)
│   ├── common/                 yaml spec loading
│   ├── sample_yml_spec/        example function/route specs
│   └── build_script.sh
│
├── quebic-faas-mgr-dashboard/  DASHBOARD binary (gorilla/mux server for web UI)
│
├── docker-files/               base Dockerfiles for function runtimes
│   ├── quebic-faas-function-nodejs-docker-file
│   └── quebic-faas-function-java-docker-file
│
├── docs/                       diagrams, draw.io source, intro deck
├── vendor/                     vendored deps (govendor: vendor/vendor.json)
├── README.md                   user-facing documentation
└── LICENSE                     Apache License v2
```

**Import path / module root:** The project predates Go modules and uses **GOPATH +
govendor** (`vendor/vendor.json`). Internal packages are imported under the
`quebic-faas/...` root (e.g. `quebic-faas/common`, `quebic-faas/types`,
`quebic-faas/quebic-faas-mgr/app`). There is no `go.mod`.

---

## 4. Shared Libraries

### 4.1 `types/` — domain model
- **`types.go`** (~328 LOC): the core entities persisted and exchanged over REST:
  - `Function` — name, version, source artifact, handler, runtime, replicas,
    events list, route, env vars, container life (awake mode + idle timeout).
  - `Route` — request method, URL, async flag, success status, target event,
    request/header attribute mappings, headers-to-pass.
  - `Event` — event id and metadata.
  - `User` — username, password (hashed), role.
  - `ManagerComponent` — describes platform components (eventbus, apigateway,
    dashboard, eventbox) and their deployment state.
  - Request/response wrappers and enums (runtimes, awake modes, time units).
- **`dto.go`** (~90 LOC): data-transfer objects for API request/response shaping.

### 4.2 `messenger/` — the event bus client (RabbitMQ / AMQP)
The abstraction every function and the gateway use to talk over **RabbitMQ**, via
the `github.com/streadway/amqp` driver. Everything flows through one durable
**topic exchange** named `quebic-faas-eventbus`, where the **eventID is the AMQP
routing key**.
- **`messenger.go`** — connection/channel management, exchange declaration,
  `Messenger` type, queue naming (`<AppID>.<routingKey>`), and a
  `WaitInit(timeout)` that also pings the RabbitMQ management HTTP API.
- **`publish.go`** — `Publish(...)` (fire-and-forget with async callbacks) and
  `PublishBlocking(...)` (caller waits). Implements **request/reply**: each publish
  generates a `requestID`, subscribes a temporary reply queue on it, and times out
  if no reply arrives. Also wires the **request-tracker** event so progress is
  recorded.
- **`subscribe.go`** — declares a durable queue, binds it to the exchange on the
  eventID routing key, and consumes (auto-ack) in a goroutine.
- **`reply.go`** — `ReplySuccess` / `ReplyError`: publish a response back on the
  requester's `requestID` routing key.
- **`baseevent.go`** — `BaseEvent` envelope: payload bytes + a header map carrying
  `eventID`, `requestID`, `created`, `statuscode`, `error`, plus arbitrary
  request/header data. Provides typed getters and JSON payload (un)marshalling.

This is the backbone of Quebic's event-driven model: routes fire events, functions
subscribe to events, and replies travel back as AMQP request/reply.

### 4.3 `auth/` — authentication & authorization
- **`jwt.go`** — issue / parse / validate JWT tokens.
- **`roles.go`** — role model: **admin | developer | tester** with permission
  checks.
- **`util.go`**, **`const.go`**, **`error_message.go`** — helpers, claim keys,
  auth error messages.

### 4.4 `config/` — base configuration
- **`config.go`** — base config structs shared by all binaries: `ServerConfig`
  (host/port) and `EventBusConfig` (RabbitMQ AMQP host/port + management
  host/port/username/password).
- **`config_util.go`** — config dir/path helpers.
- **`config_const.go`** — config keys & defaults.

### 4.5 `common/` — cross-cutting constants & helpers
A large set of single-purpose constant files plus utilities, including:
`component_const.go` (component names/ids), `event_const.go` (well-known internal
event ids), `function_const.go`, `docker_const.go` / `docker_util.go` /
`docker_file.go` (image naming, Docker helpers), `resource_const.go`,
`ingress_const.go`, `server_const.go`, `http_const.go`, `header_keys_const.go`,
`env_keys_const.go`, `consumer_ids.go`, `log_const.go`, `time_const.go`,
`error_const.go`, and `common_util.go`.

---

## 5. The Manager (`quebic-faas-mgr/`) — Control Plane

The manager is the heart of the system. Flow on startup (`app/app.go`):

1. **Load config** (`config/`) — merges flags (`cmd/flags.go`, cobra) with the
   on-disk manager config.
2. **Open DB** (`db/db.go`) — opens the embedded **BoltDB** file
   `.quebic-faas-mgr.db`, used by the DAOs. The generic DAO (`dao/dao.go`) stores
   each entity type in its own bucket keyed by the entity's ID, via a small
   `Entity` interface (`GetID/SetID/GetReflectObject/SetModifiedAt`) all domain
   types implement.
3. **Seed admin user** (`setupAdminUser`) — writes the default admin user
   (`admin/admin`, configurable) through `user_dao`.
4. **Setup deployment driver** (`setupDeployment`) — instantiates the Kubernetes
   driver (`kube_deployment`); calls `Init()` (creates namespace + RBAC admin role)
   when not running in-cluster.
5. **Setup GC** (`gc/gc.go`) and **ingress** (`IngressCreateOrUpdate`).
6. **Bring up platform components** (`components/`):
   - `eventbus_component.go` — deploy the **RabbitMQ** event bus, then connect the
     manager's `messenger` to it (`WaitInit`, up to 10 min).
   - `apigateway_component.go` — deploy a **versioned** API gateway deployment +
     service, inject eventbus env vars + an access key, and broadcast the new
     version event (old gateway versions self-drain and request shutdown via GC).
   - `mgr_dashboard_component.go` — deploy the dashboard.
   - `eventbox_component.go` — deploy the **eventbox** (request tracker) and its
     backing **MongoDB** (`eventbox-db`); the eventbox stores request-tracker
     state for async request tracking.
   - The manager also subscribes internal listeners on the bus: apigateway-data
     fetch, function-data fetch, and shutdown-request.
7. **Restore functions** (`setupFunctions`) — re-deploys every persisted function
   on startup.
8. **Start the REST API** (`httphandler/`) with CORS, on `ServerConfig` host/port.

The GC (`gc/gc.go`) deletes a deployment once all its replicas have reported
completion (driven by shutdown-request events), supporting the awake/idle
lifecycle.

### 5.1 REST API (`httphandler/`)
The manager exposes management endpoints, JWT-protected:
- `function_handler.go` (~812 LOC, the biggest handler) — create / update /
  deploy(version) / scale / delete / list / inspect / test functions, and stream
  function logs.
- `route_handler` (via `event_handler.go` + route flows) — route CRUD.
- `event_handler.go` — event CRUD.
- `auth_handler.go` — login, change-password, auth-context, user CRUD.
- `resource_handler.go` — resource/replica operations.
- `request_tracker_handler.go` — async request status lookup by request-id.
- `eventbox_handler.go` — eventbox interactions.
- `mgr_component_handler.go` — manage platform components.
- `apigateway_data_serve_handler.go` — serves route/event data to the gateway so
  it can build its routing table.
- `httphandler.go` — router setup, middleware (auth), and server wiring.

### 5.2 Function build pipeline (`function/`)
When a function is created/deployed:
1. `function_create/` validates and stores the function spec.
2. `function_util/` (~431 LOC) packages the source artifact (`.tar` / `.jar` /
   single file) per the deployment spec.
3. `function_runtime/<lang>/` selects the language runtime, generating the right
   Dockerfile/entrypoint for **Java 8 / NodeJS / Python 2.7 / Python 3.6**.
4. `function_image/` builds the container image.
5. `deployment/` deploys it (Docker or Kubernetes).
6. `function_life_controller/` enforces container life policy:
   - `awake: instant | request` (start eagerly vs. on first request)
   - `idleState.timeout` + `timeunit` for idle shutdown.

### 5.3 Deployment abstraction (`deployment/`)
- `deployment.go` — common `Deployment` interface (create/update/scale/delete
  service, list, logs).
- `docker_deployment/docker_deployment.go` (~514 LOC) — Docker driver using the
  Docker SDK.
- `kube_deployment/kube_deployment.go` (~1183 LOC) — Kubernetes driver using
  client-go: deployments, services, and ingress for functions and platform
  components. This is the primary/production target per the README.

---

## 6. The API Gateway (`quebic-faas-apigateway/`) — Data Plane

- `quebic-faas-apigateway.go` — `main()`: load config, open local store, start
  HTTP server.
- `httphandler/resource_setup.go` — builds the in-memory **route table** by
  pulling route + event definitions from the manager
  (`apigateway_data_serve_handler`).
- `httphandler/httphandler.go` (~303 LOC) — the request path:
  1. Match incoming HTTP request to a configured `Route`.
  2. Apply `requestMapping` (HTTP form/path attrs → event attrs), `headerMapping`
     (HTTP headers → event attrs), and `headersToPass`; also stamps
     `requestPath` / `requestHTTPMethod` / `async` into the event headers.
  3. Fire the target event onto the bus via the `messenger`.
  4. **Sync route** (`async: false`): `PublishBlocking` with `requestTimeout`,
     return the function's result + status (`successResponseStatus`).
     **Async route** (`async: true`): `Publish` and immediately return a
     **request-id** (in body + `requestID` header); the client later polls
     `<api-gateway>/request-tracker/{request-id}`, which the gateway resolves by
     publishing a `request-tracker-data-fetch` event answered from the eventbox.

The gateway is versioned: it listens for a `new-version-apigateway` broadcast and,
when a newer version exists, drains in-flight requests (atomic usage counter) then
publishes a shutdown-request so the manager's GC removes the old deployment.

---

## 7. The CLI (`quebic-faas-cli/`)

`quebic` is a **cobra**-based CLI (`cmd/root.go`) that authenticates against the
manager and calls its REST API via the `service/` layer (`http_client.go` +
per-domain service files). Config lives at `~/.quebic-faas/cli-config.yml`.

Command groups (from `cmd/`):
- **manager**: `quebic manager start | connect | logs | status` (`mgr.go`,
  `mgr_components.go`)
- **user**: `create | update | login | change-password | auth-context` (`user.go`)
- **function**: `create | update | deploy | scale | delete | ls | inspect | test`
  and `function logs` (`function.go`, `function_logs.go`)
- **route**: `create | update | ls | inspect` (`route.go`)
- **event** (`event.go`), **eventbox** (`eventbox.go`),
  **apigateway** (`apigateway.go`), **ingress** (`ingress.go`),
  **request-tracker** (`request_tracker.go`), **config** (`config.go`),
  **version** (`version.go`).

Function/route specs are YAML; samples live in
`quebic-faas-cli/sample_yml_spec/` (Java, NodeJS, Python 2.7/3.6, caller/receiver,
route).

---

## 8. End-to-End Runtime Flows

### 8.1 Bootstrapping the platform
```
quebic manager start     # manager deploys eventbus(RabbitMQ), apigateway,
                         # dashboard, eventbox(+mongo) into K8s/Docker
quebic manager connect   # CLI stores manager endpoint in cli-config.yml
quebic user login -u admin -p admin
```

### 8.2 Deploying a function
```
quebic function create --file function_java.yml
```
→ CLI POSTs spec (multipart, incl. the source artifact) → manager validates →
packages artifact (`function_util`) → selects runtime (`function_runtime/java_8`)
→ builds image (`function_image`, Docker build) → deploys container
(`deployment/kube` or `deployment/docker`). At runtime the function container
subscribes to its declared events on the AMQP exchange (`messenger/subscribe`).
For `awake: request` functions the manager instead subscribes a thin listener that
deploys the container on the first matching event.

### 8.3 Invoking a function
```
HTTP → apigateway → match Route → map request→event → messenger.Publish(eventID)
     → AMQP topic exchange → function container handler runs → callback/reply
     → (sync) reply travels back on requestID routing key → gateway returns result
     → (async) gateway returned request-id earlier; result stored in eventbox,
       fetched via /request-tracker/{request-id}
```

### 8.4 Function-to-function (event-driven microservices)
Inside a handler, `context.messenger().publish("users.UserValidate", payload, ...)`
fires another event onto the bus, letting functions compose into event-driven
microservice chains without direct coupling.

---

## 9. Function Lifecycle

A function moves through these phases. Each is implemented in the manager's
`function/` packages plus `deployment/`.

```
 spec (YAML)                                                   delete
    │ create/update                                              ▲
    ▼                                                            │
┌─────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐   ┌──────────┐
│ VALIDATE│──►│ PACKAGE  │──►│  BUILD   │──►│  DEPLOY   │──►│ RUNNING  │
│ + STORE │   │ artifact │   │  image   │   │ container │   │ (serving)│
└─────────┘   └──────────┘   └──────────┘   └────┬─────┘   └────┬─────┘
                                                 │              │ idle / new ver
                                                 │ scale        ▼
                                                 │         ┌──────────┐
                                                 └────────►│  REAPED  │ (GC deletes)
                                                           └──────────┘
```

### 9.1 Create / build (one-time, on `function create|update`)
Driven by `httphandler/function_handler.go → saveFunctionDTO → function_util.FunctionCreate`:
1. **Validate** — runtime must be one of `java | nodejs | python_2.7 | python_3.6`
   (`common.RuntimeValidate`); container validations run before persist.
2. **Persist spec** — the `Function` entity is saved to BoltDB; versions are tracked
   (`Version` / `Versions`, initial version `0.1.0`).
3. **Package artifact** (`function_create.CreateFunction`) — creates a build dir,
   writes a per-runtime Dockerfile (`common.DockerFileContent_*`), and copies the
   uploaded source into the build context. Per runtime:
   - **Java** — expects a `.jar`; handler = fully-qualified class
     (`com.x.HelloFunction`); baked as `/app/function.jar`.
   - **NodeJS** — `.js` single file (auto-tarred to `function_handler.tar`) or a
     `.tar`/`.gz` package; handler = `file.exportedFn` (e.g. `app.helloHandler`).
   - **Python 2.7 / 3.6** — single `.py` or a package tar; handler = `index.handler`.
4. **Build image** (`function_image.FunctionImageBuild`) — Docker build off the base
   image `quebicdocker/quebic-faas-container-<runtime>:0.1.0`, tagged
   `quebic-faas-function-<name>`, with the function's `access_key` as a build-arg.
   `Publish: true` pushes to the configured registry. The build dir is then removed.

### 9.2 Deploy (on create with `start`, on `deploy`, or on manager restart)
Driven by `function_util.FunctionDeploy → functionDeploy`:
- Computes IDs: service `quebic-faas-function-<name>`, deployment
  `…-<version>` (`GetID` / `GetDeploymentID`).
- Builds the container env (`prepareEnvKeys`) — see §11.3.
- Creates the deployment + service via the active `Deployment` driver
  (`CreateOrUpdateDeployment` + `CreateService`). On Kubernetes a port 80 service
  is added.
- Waits for readiness, then **broadcasts a new-version event**
  (`quebic-faas-event.internal.new-version.function.<name>`) so older versions can
  drain and self-terminate.
- On manager startup, `app.setupFunctions()` re-deploys every persisted function so
  the cluster state is reconciled.

### 9.3 Awake modes (cold vs. warm) — `function.life.awake`
- **`default`** — the container is deployed immediately and stays up (warm),
  subscribing directly to its events.
- **`request`** — **lazy / scale-from-zero**. The manager does *not* run the
  container; instead it subscribes a lightweight listener to the function's events.
  On the first matching event it deploys the container, then proxies the original
  event through a per-function **awake event**
  (`quebic-faas-event.internal.function-awake.<name>`) and relays the function's
  reply back to the original caller. See §10.3 and §11.4.

### 9.4 Idle / reaping — `function.life.idleState`
- Each function carries an idle policy: `timeout` + `timeunit`
  (`seconds | minutes | hours`; minimums `FunctionLifeMinAgeMinutes = 2`). This is
  passed to the container as `functionAge = "<unit>:<timeout>"`.
- The **function container** counts down its own idle timer; when idle it requests
  shutdown via the `shutdown-request` event.
- The manager's **GC** (`gc/gc.go`) tallies completion reports per deployment and,
  once all replicas have reported, deletes the deployment — yielding scale-to-zero
  for `awake: request` functions.

### 9.5 Scale / update / delete
- **Scale** — `POST /functions/scale` validates then changes the replica count via
  the deployment driver (`scaleFunction`).
- **Upgrade/Downgrade** — `POST /functions/deploy {name, version}` re-points the
  active version to an existing build (must already exist in `Versions`).
- **Delete** — `DELETE /functions/{name}` removes the deployment
  (`function_util.StopFunction → DeleteDeployment`) and the stored spec.
- **Versioning model** — every build is retained; deploy switches the live version,
  and the new-version broadcast drains the previous one.

---

## 10. Function Triggers

A function never exposes HTTP itself. **Every invocation is an event** delivered
over the RabbitMQ topic exchange; the *trigger* is whatever publishes that event.

### 10.1 Trigger sources
| Trigger | Who publishes | Path |
|---------|---------------|------|
| **HTTP via API Gateway** | apigateway, from a matched `Route` | external client → gateway → event |
| **Function → function** | another function's handler, via `messenger.publish(eventID,…)` | in-cluster, no HTTP |
| **Manager internal events** | manager components | platform control (data-fetch, new-version, shutdown, request-tracker) |
| **Direct test** | `quebic function test` → manager | publishes the function's event with a test payload |

### 10.2 What "binds" a trigger to a function
A function's `events:` list **is** its trigger set. At deploy time the runtime
subscribes a durable AMQP queue per event id (routing key = event id) on the
shared exchange:
```yaml
function:
  name: hello
  events:                 # ← triggers
    - orders.OrderCreate
    - orders.OrderUpdate
```
Routes connect HTTP to one of these events:
```yaml
route:
  requestMethod: POST
  url: /orders
  event: orders.OrderCreate   # ← fired by the gateway
  async: true
```

### 10.3 Event id namespacing (`common/event_const.go`)
User events (e.g. `orders.OrderCreate`) ride under the user-defined prefix.
Internal/control events use reserved prefixes so they never collide:
- `quebic-faas-event.defined.function.*` — user/function events
- `quebic-faas-event.execution-ack.*` — execution acknowledgements
- `quebic-faas-event.internal.*` — control plane:
  `apigateway-data-fetch`, `function-data-fetch`, `request-tracker`,
  `request-tracker-data-fetch`, `new-version.*`, `function-awake.<name>`,
  `shutdown-request`.

### 10.4 Sync vs. async triggering (route `async` flag)
- **`async: false`** — gateway uses `PublishBlocking`, waits up to the route's
  `requestTimeout`, returns the handler's result with `successResponseStatus`.
- **`async: true`** — gateway `Publish`es and immediately returns a **request-id**;
  the caller later polls `/request-tracker/{request-id}` (resolved from the
  eventbox) for the result.

### 10.5 Cold-start trigger (awake: request)
For a `request`-awake function the *first* trigger event is intercepted by the
manager's listener, which deploys the container and then re-fires the event through
the function-awake channel; subsequent events hit the now-warm container directly
until it idles out and is reaped (§9.4).

---

## 11. How API Gateway, Event Bus, and Functions Connect

Everything is wired through one durable **RabbitMQ topic exchange**
`quebic-faas-eventbus`, where **eventID = AMQP routing key** and replies travel
back on a per-request routing key (`requestID`). No component calls a function over
HTTP directly.

```
                        RabbitMQ topic exchange "quebic-faas-eventbus"
  ┌───────────┐   publish(event)        ┌──────────────────────────────┐
  │ API       │ ───────────────────────►│  routingKey = eventID         │
  │ Gateway   │                         │                               │
  │           │◄─── reply(requestID) ───┤  routingKey = requestID       │
  └───────────┘                         └───────────┬──────────────────┘
                                                    │ bound queue
                                                    ▼
                                          ┌───────────────────┐
                                          │ Function container │
                                          │ subscribe(eventID) │
                                          │ handler→reply      │
                                          └───────────────────┘
       ▲ all three connect with the same EventBusConfig (host/port/creds)
       │ manager injects it as container env vars at deploy time
  ┌────┴─────┐
  │ Manager  │  also subscribes: data-fetch, request-tracker, shutdown
  └──────────┘
```

### 11.1 Connection bootstrap (the `messenger` package)
Each process builds a `Messenger{AppID, EventBusConfig}` and calls `Init()`
(or `WaitInit(timeout)`), which:
1. dials `amqp://<user>:<pass>@<host>[:<port>]/` (driver `streadway/amqp`),
2. opens a channel,
3. declares the durable topic exchange `quebic-faas-eventbus`.
Queues are named `<AppID>.<routingKey>` so each app gets its own queue per event.

### 11.2 The connection points
- **API Gateway → Bus.** On boot the gateway pulls its route table from the manager
  (publishes `…internal.apigateway-data-fetch`, manager replies with
  `ApigatewayData`). For each `Route` it registers an HTTP handler; on a request it
  maps HTTP→event and publishes to `route.Event`. (`quebic-faas-apigateway/httphandler/httphandler.go`)
- **Bus → Function.** At deploy the container subscribes a durable queue bound to
  each of its `events:` (routing keys) and runs the handler per message
  (`messenger/subscribe.go`).
- **Function → Bus (reply).** The handler's `callback`/`success` publishes back on
  the request's `requestID` routing key (`messenger/reply.go`), which the original
  publisher is subscribed to — implementing request/reply.
- **Function → Function.** A handler may `publish` any other event id, chaining
  functions without HTTP.
- **Manager ↔ Bus.** The manager subscribes control listeners
  (apigateway-data-fetch, function-data-fetch, request-tracker,
  request-tracker-data-fetch, shutdown-request) and publishes new-version events.

### 11.3 How a function learns to connect (env vars injected at deploy)
`function_util.prepareEnvKeys` bakes the bus coordinates + event ids into the
container environment, so the runtime can connect and subscribe on startup:
- Bus: `rabbitmq_host`, `rabbitmq_port`, `rabbitmq_exchange`,
  `rabbitmq_management_username/password`
  (in Docker mode host = the `quebic-faas-eventbus` service name, port `0`).
- Identity: `appID`, `deploymentID`, `version`, `access_key`.
- Events: `events` (CSV of trigger event ids), plus the reserved event ids
  (`eventConst_eventLog`, `…eventFunctionAwake`, `…eventDataFetch`,
  `…eventNewVersion`, `…eventShutDownRequest`).
- Handler: `artifactLocation` (`HandlerFile`), `functionPath` (`HandlerPath`).
- Lifecycle: `functionAge` (`<unit>:<timeout>`).
- Tracking: `EVENTBOX_URI` (eventbox HTTP endpoint).

### 11.4 Request lifecycle across all three (sync route example)
```
client ─HTTP─► API Gateway
                 │ map request→event, stamp headers (requestPath, async, …)
                 │ PublishBlocking(route.Event)            [routingKey=eventID]
                 ▼
            RabbitMQ exchange ──► function queue ──► handler runs
                                                       │ callback(result)
                                                       ▼
                                  ReplySuccess(requestID)   [routingKey=requestID]
                 ┌─────────────────────────────────────────┘
                 ▼
            API Gateway receives reply ─HTTP─► client (status = successResponseStatus)
            (in parallel: request-tracker events recorded in eventbox)
```
For an **async** route the gateway returns the `requestID` immediately and the
result is retrieved later from the eventbox via `/request-tracker/{requestID}`.

### 11.5 Network / ports (defaults)
- Manager REST: `:1028` · API Gateway: `:3000` · Dashboard: `:80`
- RabbitMQ: AMQP `:5672`, management `:15672` (default creds `guest/guest`)
- Eventbox: `:8080` · Eventbox Mongo: `:27017`
- Eventbus image: `rabbitmq:3.7-management-alpine`

---

## 12. RabbitMQ Usage & Messaging Patterns

RabbitMQ is the single backbone of the platform — there is no other inter-service
transport. It is **deployed by the manager itself** as the `quebic-faas-eventbus`
component (image `rabbitmq:3.7-management-alpine`, AMQP `:5672`, management
`:15672`, default creds `guest/guest`), and every other process (manager, gateway,
each function container, eventbox) connects to it. All of the integration lives in
the shared **`messenger/`** package via the `github.com/streadway/amqp` driver.

### 12.1 Topology: one durable topic exchange
- A single **topic exchange** named **`quebic-faas-eventbus`** (`messenger.Exchange`)
  is declared by every process on connect (`messenger.go → Init`):
  `ExchangeDeclare(name, "topic", durable=true, autoDelete=false, internal=false, noWait=false)`.
- **eventID = routing key.** Publishing to an event id routes to whichever queues
  bound that key; subscribing binds a queue to that key.
- **Per-app queues.** Each process queue is named `<AppID>.<routingKey>`
  (`getQueueName`) — so the same event can fan out to multiple distinct consumers
  (manager, gateway, function), each with its own queue.
- **Connection.** `amqp.Dial("amqp://<user>:<pass>@<host>[:<port>]/")`; one
  `*amqp.Connection` + one `*amqp.Channel` per `Messenger`. `WaitInit(timeout)`
  retries the dial **and** HTTP-pings the management API (`pingToRabbitmqMgr`,
  `http://<mgmtHost>:<mgmtPort>`) before declaring ready — used by the manager to
  wait (up to 10 min) for the freshly-deployed broker.

#### 12.1.1 Exchange → Queue → Bind → Consumer/Handler

The whole subscribe path is: a publisher sends to the **exchange** with
`routingKey = eventID`; a **queue** named `<appID>.<eventID>` is asserted and
**bound** to that exchange on the same routing key; a **consumer** drains the queue
and invokes the **handler**. Because the queue name is derived purely from the
event id (`getQueueName(routingKey) = appID + "." + routingKey`), the queue is the
unit shared per event — not per handler.

```
                        Publisher
                  publish(eventID, payload)
                  routingKey = eventID
                            │
                            ▼
        ┌───────────────────────────────────────────┐
        │             RabbitMQ EXCHANGE              │
        │      messengerConfig.rabbitmq.exchange     │
        │           "quebic-faas-eventbus"           │
        │              (topic, durable)              │
        └───────────────────────────────────────────┘
             │                              │
   routingKey = "EVENT_A"        routingKey = "EVENT_B"
   (QueueBind)                   (QueueBind)
             │                              │
             ▼                              ▼
   ┌───────────────────┐         ┌───────────────────┐
   │       QUEUE        │         │       QUEUE        │
   │ appID + ".EVENT_A" │         │ appID + ".EVENT_B" │
   │ durable, autoDelete│         │ durable, autoDelete│
   └───────────────────┘         └───────────────────┘
             │                              │
        consume(noAck,                 consume(noAck,
        exclusive:false)               exclusive:false)
             │                    ┌─────────┴─────────┐
             ▼                    ▼                   ▼
   ┌───────────────────┐  ┌──────────────┐   ┌──────────────┐
   │   CONSUMER +      │  │  CONSUMER +  │   │  CONSUMER +  │
   │   Handler_1       │  │  Handler_2   │   │  Handler_3   │
   │ (listens EVENT_A) │  │(listen EVENT_B)│ │(listen EVENT_B)│
   └───────────────────┘  └──────────────┘   └──────────────┘

   ── SINGLE handler ──    ──── MULTIPLE handlers, one event ────
   one unique event        same queue (appID.EVENT_B), competing
   → its own queue         consumers → round-robin delivery
```

**Single handler, unique event** (left): `subscribe("EVENT_A", handler1)` asserts
queue `appID.EVENT_A`, binds it on routing key `EVENT_A`, and attaches one
consumer. One event → one queue → one handler.

**Multiple handlers, one event** (right): `subscribe("EVENT_B", handler2)` and
`subscribe("EVENT_B", handler3)` both compute the **same** queue name
`appID.EVENT_B` (queue name is a function of the routing key only). `assertQueue`
is idempotent, so both attach as **competing consumers on the same queue**.

> ⚠️ **Delivery semantics:** because the two handlers share one queue (not a
> fan-out queue per handler), RabbitMQ delivers each `EVENT_B` message to **only
> one** of them (round-robin / competing-consumers), **not** to both. With
> `noAck: true` there is no redelivery either. To make every handler receive every
> message (true pub/sub fan-out within one app), each handler would need its own
> distinct queue (e.g. fold a handler/consumer id into `getQueueName`) bound to the
> topic exchange on the same routing key.

#### 12.1.2 Fan-out across functions + competing consumers within replicas

This is the case Quebic actually relies on for multi-function event delivery. The
queue name is `getQueueName(routingKey) = appID + "." + routingKey`, and `appID`
is injected per **function** (the function name/id), **not** per replica
(`config/messenger_config.js → process.env['appID']`). So replicas of the same
function share one queue, while different functions get different queues.

**Setup:**
- `function_a` → `appID = "function_a"`, 2 replicas
- `function_b` → `appID = "function_b"`, 2 replicas
- Both subscribe to the **same** routing key `EVENT_X`

**Queues the code produces:**

| Process | appID | Queue name |
|---------|-------|-----------|
| function_a replica 1 | `function_a` | `function_a.EVENT_X` |
| function_a replica 2 | `function_a` | `function_a.EVENT_X` ← same |
| function_b replica 1 | `function_b` | `function_b.EVENT_X` |
| function_b replica 2 | `function_b` | `function_b.EVENT_X` ← same |

→ **two distinct queues**, both bound to the exchange on `EVENT_X`; within each
queue the 2 replicas attach as competing consumers (`exclusive: false`).

```
                         Publisher
                  publish("EVENT_X", payload)
                  routingKey = EVENT_X
                            │
                            ▼
        ┌───────────────────────────────────────────┐
        │             RabbitMQ EXCHANGE              │
        │           "quebic-faas-eventbus"           │
        │                  (topic)                   │
        └───────────────────────────────────────────┘
            │  bind EVENT_X            │  bind EVENT_X
            ▼                          ▼
  ┌──────────────────────┐  ┌──────────────────────┐
  │       QUEUE          │  │       QUEUE          │
  │  function_a.EVENT_X  │  │  function_b.EVENT_X  │
  └──────────────────────┘  └──────────────────────┘
       │           │              │           │
   consume     consume        consume     consume
       ▼           ▼              ▼           ▼
  ┌─────────┐ ┌─────────┐   ┌─────────┐ ┌─────────┐
  │ fn_a r1 │ │ fn_a r2 │   │ fn_b r1 │ │ fn_b r2 │
  └─────────┘ └─────────┘   └─────────┘ └─────────┘
   └──── one of these ───┘    └──── one of these ───┘
        handles it                 handles it

  ── FAN-OUT across functions (different queues) ──
  ── COMPETING CONSUMERS within a function (replicas share one queue) ──
```

**Result (the intended behavior):**
1. **Fan-out across functions** — `function_a` and `function_b` have *different*
   `appID`s → *different* queues, so RabbitMQ copies each `EVENT_X` message into
   **both** queues. Both functions receive the event. ✅
2. **One instance per function** — the 2 replicas of a function share the *same*
   queue, so they are competing consumers and the broker delivers each message to
   **only one** replica (round-robin). Exactly one `function_a` replica and one
   `function_b` replica process it. ✅

This is the standard **"fan-out by topic + competing consumers per group"** pattern,
where **the queue name is the consumer group**: a different queue name = an
independent copy of the message (fan-out); the same queue name = load-balanced
within the group. The ⚠️ caveat in §12.1.1 only bites when consumers share the
*same* queue (same `appID` + routing key); distinct functions avoid it by having
distinct `appID`s.

### 12.2 Subscribe side (`subscribe.go`)
For each subscription:
```
QueueDeclare(name="<AppID>.<eventID>", durable=true, autoDelete=true, exclusive=false)
QueueBind(queue, routingKey=eventID, exchange="quebic-faas-eventbus")
Consume(queue, consumerID, autoAck=true, exclusive=false)
→ goroutine: for msg := range msgs { handler(BaseEvent{body, headers}) }
```
Notable choices:
- **Auto-ack** (`autoAck=true`) — fire-and-forget delivery; no manual `Ack`/`Nack`,
  so there is no redelivery/at-least-once guarantee at the AMQP layer.
- **autoDelete queues** — a consumer's queue disappears when it goes away, keeping
  the broker clean as functions/gateways scale and version-rotate.
- **Unique consumer tags** — `consumerID = <name>-<UUID>` (`common/consumer_ids.go`)
  so cancelling one consumer never disturbs another.

### 12.3 Publish side & the Request/Reply (RPC) pattern (`publish.go`)
This is the core pattern. Quebic implements **RPC-over-AMQP** by hand (it does not
use AMQP's `ReplyTo`/`CorrelationId` fields — it uses a dedicated reply *routing
key* instead):
1. Each publish mints a **`requestID` (UUID)** stored in the message headers
   (`baseevent.go`), alongside `eventID`, `created`, `statuscode`, `error`.
2. The publisher **subscribes a temporary queue bound to `requestID`** as the reply
   channel, then `channel.Publish(exchange, routingKey=eventID, body, headers)` with
   `ContentType: application/json`.
3. The consumer handles the event and **replies by publishing to the `requestID`
   routing key** (`reply.go → ReplySuccess/ReplyError`), which the publisher's reply
   queue receives.
4. On reply (or timeout) the publisher **releases the reply queue** (`ReleseQueue` =
   `channel.Cancel` + `QueueDelete`), so reply queues are ephemeral and
   self-cleaning.

Two front-end variants over the same `publish` core:
- **`Publish(...)`** — non-blocking: the wait runs in a goroutine; success/error
  handlers fire asynchronously. Used for async routes and fire-and-forget
  broadcasts (new-version, shutdown-request).
- **`PublishBlocking(...)`** — blocks the caller until the reply arrives or the
  timeout elapses. Used by the API gateway for **sync routes** and by the manager
  for `function test` and internal data-fetch calls.
- **Timeout** — `select { case <-waitForResponse: … ; case <-time.After(timeout): release + errHandler }`;
  default `2h` when none supplied; the gateway passes the route's `requestTimeout`.

### 12.4 Patterns in use (summary)
| Pattern | Where | How it maps to AMQP |
|---------|-------|---------------------|
| **Request/Reply (RPC)** | gateway↔function (sync), manager internal data-fetch, `function test` | publish on `eventID`, reply on temp queue bound to `requestID` |
| **Publish/Subscribe (fan-out by topic)** | events delivered to functions; one event → many subscribers | topic exchange + per-app queues on the same routing key |
| **Fire-and-forget broadcast** | `new-version.*`, `shutdown-request`, `request-tracker` logging | `Publish` with no/ignored reply handler |
| **Async + deferred result** | async routes | `Publish` returns `requestID`; result later fetched from eventbox |
| **Competing consumers** | a function with `replicas > 1` | replicas share `<AppID>.<eventID>`; broker round-robins deliveries |
| **Event-carried state transfer** | function→function chains | handler `publish`es a new event id with its payload |

### 12.5 The message envelope (`BaseEvent`)
Every message is a `BaseEvent` = **JSON payload bytes + an AMQP header map**. The
headers carry the protocol metadata (`eventID`, `requestID`, `created`,
`statuscode`, `error`) plus user data (mapped request attributes, passed HTTP
headers, `requestPath`, `requestHTTPMethod`, `async`). Helpers
(`GetPayloadAsObject`, `ParsePayloadAsObject`, typed header getters) wrap
(de)serialization so handlers work with objects, not raw frames.

### 12.6 Reserved routing keys (control plane on the same bus)
The control plane multiplexes onto the same exchange using reserved event ids
(`common/event_const.go`), so platform traffic and user traffic share one broker
without collision:
- `…internal.apigateway-data-fetch` / `…function-data-fetch` — RPC to the manager.
- `…internal.request-tracker` — fire-and-forget logging into the eventbox.
- `…internal.request-tracker-data-fetch` — RPC to read async results.
- `…internal.new-version.<component|function>.*` — broadcast for rolling versions.
- `…internal.function-awake.<name>` — cold-start proxy channel.
- `…internal.shutdown-request` — drained component asks the GC to delete it.

### 12.7 Operational characteristics & trade-offs
- **At-most-once delivery** — auto-ack + autoDelete queues favor simplicity and
  cleanliness over durability; an in-flight message is lost if a consumer crashes
  mid-handle, and replies rely on liveness within the timeout window.
- **No dead-letter / retry** configured — failures surface as RPC timeouts or error
  replies, not broker-level redelivery.
- **Broker is a hard dependency & single component** — `replicas: 1` for the
  eventbus; if it is down, the manager blocks on `WaitInit` and the gateway records
  an `eventbus-connect-failed` health status.
- **Credentials/coordinates are injected** into every function container as env vars
  (`rabbitmq_host/port/exchange/management_username/password`) at deploy time
  (see §11.3); in Docker mode the host is the service name `quebic-faas-eventbus`.

---

## 13. Data Storage Choices (Why NoSQL / Key-Value, not Relational)

Quebic uses two non-relational stores and **no SQL database**:

| Store | Used by | Holds | Driver |
|-------|---------|-------|--------|
| **BoltDB** (`coreos/bbolt`) | manager + apigateway | platform state: `Function`, `Resource`(route), `Event`, `User`, `ManagerComponent`, `RequestTracker` | embedded Go lib |
| **MongoDB** | eventbox component | request-tracker / execution results for async invocations | `eventbox` service |

### 13.1 What the access pattern actually is
The decision follows directly from how the code touches data. The entire manager
DAO (`dao/dao.go` + the per-entity DAOs) only ever does **whole-object
get/put keyed by one natural ID**:
- `Add` / `Update` / `Save` → marshal the struct to JSON, `bucket.Put(id, json)`.
- `GetByID` → `bucket.Get(id)`, unmarshal.
- `GetAll` → `bucket.ForEach` (full scan of a small bucket).
- Each entity type maps to its own **bucket** (`Function`, `User`, …); the key is
  the entity's single ID (`function.Name`, `user.Username`, `resource.URL:method`).

There are **no joins, no multi-field queries, no foreign keys, and no
cross-entity transactions** anywhere. Even "filtered" reads (e.g. find a route by
name) are done by scanning a bucket in Go, not by a query — see
`ResourceHandler` iterating `GetAll` and breaking on a match. The data is
naturally **document-shaped** (a `Function` embeds its route, env vars, life policy,
versions, and last log as nested structs/arrays in one JSON blob).

### 13.2 Why BoltDB (key-value) for platform state
- **Access pattern is pure key-value.** Every read/write is by primary ID over a
  self-contained JSON document — exactly what a K/V store is best at. A relational
  schema would add table/column modeling and an ORM/marshalling layer for zero
  query benefit.
- **Embedded = zero operational footprint.** Bolt is an in-process Go library
  writing a single file (`.quebic-faas-mgr.db`, `.quebic-faas-apigateway.db`). The
  manager and gateway ship as **`FROM scratch`** static binaries with no sidecar DB
  container, no connection pool, no DBA, no migrations — the platform can bootstrap
  itself before any database service exists.
- **ACID where it's needed.** Bolt gives serializable ACID transactions
  (`db.Update` / `db.View`) and an `Entity` abstraction, which is enough for the
  small, low-write control-plane state. A full RDBMS's concurrency/scaling machinery
  is unnecessary here.
- **Schema-flexible documents.** Nested, evolving structs (function spec + route +
  env + life + versions) serialize as one JSON value; adding a field is a struct
  change, not a migration. Relational normalization would fragment one logical
  function across several tables for no gain.
- **Small, bounded dataset.** Entities number in the tens/hundreds, so full-bucket
  `ForEach` scans are cheap — there's no need for the secondary indexes/query
  planner a relational engine provides.

### 13.3 Why MongoDB for the eventbox (request tracker)
- **Append-heavy, document-per-request data.** A `RequestTracker` accumulates a
  growing `logs` array and a free-form `response.message interface{}` (arbitrary
  function output). A document store maps to this `{requestID → {logs[], response}}`
  shape directly; a relational model would need a `logs` child table and a way to
  store schemaless payloads.
- **Variable / schemaless payloads.** Function responses have no fixed schema, so a
  schemaless document fits better than typed SQL columns.
- **Independently scalable from control-plane state.** The eventbox is a **separate
  deployed component** (its own service + `mongo` container, §5.6) because async
  request tracking is higher-volume and on the request path, unlike the rarely-
  written Bolt state. Keeping it in a standalone server (not embedded) lets it scale
  and be queried by request id without coupling to the manager process.

### 13.4 Net rationale
A FaaS control plane is **a small set of configuration documents read/written by
primary key**, plus **a stream of per-request result documents**. Neither workload
needs relational joins, normalization, or a query language — so Quebic picks the
stores whose strengths match: an **embedded K/V engine** for self-contained,
operationally-free platform state, and a **document database** for schemaless,
append-style request tracking. Choosing an RDBMS would add operational weight and
modelling overhead with no functional payoff for these access patterns.

---

## 14. Technology Stack

| Concern | Technology |
|---------|-----------|
| Language | Go (GOPATH + govendor, no go.mod) |
| Function runtimes | Java 8, NodeJS, Python 2.7, Python 3.6 |
| Orchestration | Kubernetes (`k8s.io/client-go`) — primary; Docker SDK (swarm services) — alt mode |
| Event bus / messaging | RabbitMQ / AMQP (`streadway/amqp`), one topic exchange `quebic-faas-eventbus`, request/reply |
| Request-tracker store | MongoDB (`eventbox-db`, behind the eventbox) |
| Persistence (manager) | Embedded BoltDB (`coreos/bbolt`), file `.quebic-faas-mgr.db` |
| AuthN/Z | JWT (`dgrijalva/jwt-go`) + role model (admin/developer/tester) |
| HTTP router | `gorilla/mux` (+ `gorilla/handlers` CORS, `gorilla/context`) |
| Management API | REST over HTTP |
| CLI | `spf13/cobra` (+ `olekukonko/tablewriter`, `fatih/color`) |
| Function spec | YAML deployment/route specs |
| Java runtime SDK | `quebic-faas-runtime-java` (external lib) |
| License | Apache License v2 |

---

## 15. Key Design Characteristics

- **Clear control-plane / data-plane split:** the manager builds & deploys;
  the gateway routes live traffic. They communicate over REST + a shared event bus.
- **Event-first model:** HTTP routes are thin adapters that translate requests into
  AMQP events; functions never expose HTTP directly. This makes function-to-function
  composition and async invocation first-class.
- **Pluggable deployment backend:** a single `Deployment` interface with Docker and
  Kubernetes implementations.
- **Pluggable language runtimes:** each runtime is an isolated package producing the
  appropriate container build, so adding a language is additive.
- **Self-bootstrapping:** the manager deploys its own supporting components
  (eventbus, gateway, dashboard, eventbox) into the cluster.
- **Lifecycle-aware functions:** per-function awake mode and idle timeout, reaped by
  a background GC — a scale-to-zero style capability.

---

## 16. Build & Run

- Each binary is built from its own `main` package
  (`quebic-faas-mgr/quebic-faas-mgr.go`, `quebic-faas-apigateway/...go`,
  `quebic-faas-cli/quebic-faas-cli.go`, `quebic-faas-mgr-dashboard/...go`).
- `quebic-faas-cli/build_script.sh` builds the CLI; manager and gateway ship
  `Dockerfile`s for containerized deployment.
- Function runtime base images are defined in `docker-files/`.
- There is **no top-level `main.go`, `Makefile`, or `go.mod`** — builds rely on
  GOPATH layout under `quebic-faas/` and the vendored deps in `vendor/`.

---

## 17. Notes / Caveats

- The repo uses **govendor** (`vendor/vendor.json`), not Go modules; to build today
  you would place it under `$GOPATH/src/quebic-faas/` or migrate to modules.
- The dashboard binary serves a UI; the front-end asset bundle itself is delivered
  as a packaged component rather than checked-in source here.
- Documentation primarily lives in the root `README.md` (per-runtime programming
  models, spec formats, and CLI reference) plus the diagrams/deck in `docs/`.
```

---

## 18. Future Improvements

### 18.1 Awake via KEDA (Kubernetes Event Driven Autoscaling), triggered by RabbitMQ

**Today (§9.3–§9.4).** For `awake: request` functions Quebic hand-rolls its own
scale-from-zero: the manager subscribes a lightweight listener to the function's
events, deploys the container on the first matching event, proxies that event
through a `function-awake.<name>` channel, and later the container idles out and
publishes a `shutdown-request` that the manager's **GC** acts on to delete the
deployment. This works but couples cold-start and reaping to bespoke control-plane
code (the listener, the awake-proxy channel, and the GC completion tallying).

**Proposed.** Delegate scale-from-zero (and scale-to-zero) to **KEDA** using its
**RabbitMQ scaler**. KEDA watches the function's queue depth on the
`quebic-faas-eventbus` broker and drives the Kubernetes Deployment's replica count
directly — `0 → N` when messages are waiting, back to `0` after they drain. Quebic
keeps owning build/deploy and the event model; KEDA owns the autoscaling loop.

**What this replaces / simplifies:**
- The manager's per-function **awake listener** and the `function-awake.<name>`
  proxy channel — no more intercepting the first event to boot the container.
- The **`shutdown-request` + GC completion-tally** path for idle reaping — KEDA
  scales the deployment to zero on an idle cooldown instead.
- Idle policy (`idleState.timeout`) maps onto KEDA's `cooldownPeriod`;
  `awake: request` maps onto `minReplicaCount: 0`.

**What Quebic would generate at deploy time** (in `kube_deployment`): alongside the
Deployment + Service, emit a `ScaledObject` (and, if per-message concurrency is
desired, a `ScaledJob`) pointing at the function's queue.

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: quebic-faas-function-<name>
spec:
  scaleTargetRef:
    name: quebic-faas-function-<name>-<version>   # the Deployment
  minReplicaCount: 0            # ← awake: request  (scale-from-zero)
  maxReplicaCount: <replicas>
  cooldownPeriod: <idleState.timeout in seconds>   # ← idle reaping
  triggers:
    - type: rabbitmq
      metadata:
        protocol: amqp
        queueName: <appID>.<eventID>               # §12.1: getQueueName
        mode: QueueLength
        value: "1"
      authenticationRef:
        name: quebic-faas-rabbitmq-auth            # host/creds from EventBusConfig
```

**Trade-offs to weigh:** adds KEDA + its metrics adapter as a cluster dependency;
`autoDelete: true` queues (§12.2) must become durable/non-autoDelete so a queue
exists for KEDA to poll while replicas are at zero; and the current auto-ack
"at-most-once" semantics may want revisiting so queue depth accurately reflects
outstanding work.

#### ASCII diagram — KEDA-driven awake on RabbitMQ

```
                        external client / another function
                                     │  HTTP  (or function→function publish)
                                     ▼
                            ┌──────────────────┐
                            │   API Gateway     │  map request → event
                            └────────┬─────────┘
                                     │ publish(eventID)         [routingKey = eventID]
                                     ▼
        ┌──────────────────────────────────────────────────────────────┐
        │       RabbitMQ topic exchange  "quebic-faas-eventbus"          │
        └───────────────────────────────┬──────────────────────────────┘
                                         │ bound (durable, NOT autoDelete)
                                         ▼
                            ┌──────────────────────────┐
                            │   QUEUE  <appID>.<eventID> │   ← messages pile up
                            │   depth = N (>0)           │      while replicas = 0
                            └────────────┬──────────────┘
                                         │
             poll queue depth            │ (QueueLength via AMQP / mgmt API)
        ┌────────────────────────────────┤
        │                                ▼
   ┌─────────┐   metrics   ┌──────────────────────────┐
   │  KEDA    │◄───────────┤  RabbitMQ scaler / trigger │
   │ operator │            └──────────────────────────┘
   │          │
   │  drives  │  patch replicas:  0 → N   (depth>0)      ┌───────────────────────┐
   │  replica │───────────────────────────────────────► │ Kubernetes Deployment  │
   │  count   │  patch replicas:  N → 0   (idle cooldown)│ quebic-faas-function-* │
   └─────────┘                                           └───────────┬───────────┘
        ▲                                                            │ pods start
        │ cooldownPeriod  (= idleState.timeout)                      ▼
        │ elapses with depth == 0                          ┌───────────────────┐
        └──────────────────────────────────────────────── │ Function container │
                                                           │ subscribe(eventID) │
                                                           │  handler → reply   │
                                                           └───────────────────┘

   ── SCALE-FROM-ZERO ──   queue depth > 0  →  KEDA scales Deployment 0 → N
   ── SCALE-TO-ZERO   ──   depth == 0 for cooldownPeriod  →  KEDA scales N → 0
   (replaces: manager awake-listener, function-awake proxy, shutdown-request + GC)
```
