## [0.0.35]

### Fixed

-   Fixed updating of stock items settings.

## [0.0.34]

### Added

-   Добавлена проверка предпоследней цены для предотвращения воздействия
    скачкообразного объявления цены на премаркет
-   Минимальная цена для покупки бумаги

### Fixed

-   Рассчет профита в статистике слотов

## [0.0.33]

### Added

-   Добавлена медиана стэка для более сдержаных покупок. Обычно в стэке хранится
    5 последних цен на основе которых высчитывается SMA frame и определяется
    трэнд. Среднея цена сравнивается со средней после обновления стэка, тренд
    вверх, если последнее число больше.
    Разница цены по медиане от последней вошедшей в стэк цены должна быть
    отрицательным числом, тоже самое с разницей от средней цены к последней,
    таким образом будет получен более сильный крен трэнда вниз.
-   Добавлена функция расчета вхождения времени в заданный интервал,
    требуется сброс всех стэков раз в сутки, чтобы данные вчерашнего дня
    не влияли на принятие решения сейчас. Связано больше с торговлей
    на премаркете и переходных состояниях биржи.
-   В режиме разработки условия выходных дней игнорируется.

## [0.0.32]

### Fixed

-   #11 added a check for zero value of the payment field

### Added

-   detailed data to logs
-   read config for logger

## [0.0.31]

-   add hash to filename js-client
-   static with relative link

## [0.0.30]

-   add workhours to js-client

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
