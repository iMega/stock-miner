## [0.0.29]

-   add workhours to js-client

## [0.0.28]

-   add banner with version to js client

## [0.0.27]

### Fixed

-   Too many re-renders

### Added

-   Main settings: global mining

## [0.0.26]

-   test render problem

## [0.0.25]

-   fixed start/stop status button

## [0.0.24]

-   fixed start/stop status button

## [0.0.23]

-   fixed start/stop status button

## [0.0.22]

-   marketCredentials

## [0.0.21]

-   marketCredentials

## [0.0.20]

### Fixed

-   a.settings.marketCredentials[0] is undefined

## [0.0.19]

-   Works only on weekdays

## [0.0.18]

-   revert, remove defence handler for static files

## [0.0.17]

-   remove defence handler for static files

## [0.0.16]

-   test GITHUB_REF

## [0.0.15]

-   check values of range stock item

## [0.0.14]

-   test GITHUB_REF

## [0.0.13]

-   test GITHUB_REF

## [0.0.12]

-   test GITHUB_REF

## [0.0.11]

-   Добавлена проверка корридора цены с yahoo finance. Бот сможет делать
    покупку если текущая цена входит в диапазон.

## [0.0.10]

### Fixed

-   #9 Исправлена работа с журналом операций. В случае невозможности получить
    результат сделки, будет создано отложенное задание. Количество попыток
    ограничено 20.

## [0.0.9]

### Changed

-   set url for ws

## [0.0.8]

### Fixed

-   fix render ssr

## [0.0.7]

### Fixed

-   window to globalThis

## [0.0.6]

### Fixed

-   env GRAPHQL_HOST
-   env GRAPHQL_SCHEMA

## [0.0.5]

### Added

-   Add version to title
-   Add env GRAPHQL_HOST - graphql host, default is 127.0.0.1
-   Add env GRAPHQL_SCHEMA - graphql schema, default is http

## [0.0.4]

### Added

-   If build has error then CI to stop.
-   Disable reconnect websocket.

### Changed

-   Update Nodejs image to 16.1.0

### Fixed

-   Build frontend for SSR with StaticRouter
-   Dealings page doesn't work because not check input data.

## [0.0.3]

### Added

-   Page Profile, dealings and slot
-   Test automatically buy and sell
-   Migration process

## [0.0.2]

### Added

-   Build frontend

## [0.0.1]

### Added

-   Support daemon kit
