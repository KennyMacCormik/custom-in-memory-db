Это простая in-memry key-value база данных

Архитектура выглядит таким образом

![](/simple-design.jpg)

Синтаксис запросов в базу:

>[!IMPORTANT]
>`query = set_command | get_command | del_command`
>
>`set_command = "SET" argument argument`
>
>`get_command = "GET" argument`
>
>`del_command = "DEL" argument`
>
>`argument    = punctuation | letter | digit { punctuation | letter | digit }`
>
>`punctuation = "*" | "/" | "_" | ...`
>
>`letter      = "a" | ... | "z" | "A" | ... | "Z"`
>
>`digit       = "0" | ... | "9"`

Всего реализовано несколько сущностей:

1. TCP server
- `Server struct`
  
   Принимает подключения по tcp от cli, для обработки подключения сервер использует функцию с сигнатурой `func(r io.Reader, lg *slog.Logger) (string, error)`, которую принимает в качестве параметра в метод `Listen()`. Обработка каждого клиента запускается в отдельной горутине.
- `connMeter struct`
  
   Контролирует, чтобы единомоментно сервер не принимал больше чем `NET_MAX_CONN` подключений.
2. Database
- `Database struct`

  Абстрагирует всю логику работы базы данных, оставляя только один метод `HandleRequest(r io.Reader, lg *slog.Logger) (string, error)`, который возвращшает либо результат запроса к базе, либо ошибку, если таковая возникла в процессе обработки запроса.
3. Parser
- `Read(r io.Reader, lg *slog.Logger) (Command, error)`

  Единственная публичная функция в пакете Parser отвечает за то, чтобы проанализировать переданные данные и вернуть ошшибку, если данные не содержат корерктной команды, либо вернуть `Command struct`, которая будет обработана далее.
- `Command struct`

  Содержит команду, которую будет выполнять база, а также её аргументы.
4. Compute
- `Compute interface`

  Данный интерфейцс призван принимать корректную команду от функции `Read()` из пакета `parser` и отправлять её на исполнение через метод `Exec(cmd parser.Command, lg *slog.Logger) (string, error)`.
- `Comp struct`

  Реализация интерфейса `Compute`, которая требует инициализации через метод `New()`, который в качестве параметра принимает интерфейс `Storage`. Реаализует метод `Exec()` в котром вызывает медоты интерфейса `Storage` в соответствии с переданной командой.
5. Storage
- `Storage interface`

  Данный интерфейс отражает действия с данными в базе данных. Содержит методы на каждую домустимую команду, а также методы `Close() error` для корреткного завершения работы и метод `Recover(conf cmd.Config, lg *slog.Logger)` error для восстановления данных, если конкретная реализация это позволяет.
- `MapStorage struct`

  Реализация интерфейса `Storage`, которая использует `map[string][string]` для хранения данных и `sync.Mutex` для организации потокобезопасного доступа к `map`. Методы `Close()` и `Recover()` являются заглушками, так как данная реализация не предполагает возможности восстановления данных и не требует специального завершения работы.
6. Wal
- `Wal struct`

  Реализует интерфейс `Storage`. Для корректного завершения работы требуется вызвать метод `Close()`. Реализует васстановление данных через метод `Recover()` (wal логи). Инициализируется методом `New(conf cmd.Config, st storage.Storage) error`. Является wrapper'ом для настоящей имплементации интерфейса `Storage`, так что принимает данный интерфейс, как один из памаетров для инициализации.

  Данный интерфейс призван реализовать технологию write-ahead logging, а также дополнительную логику. Он принимает команду, которая предназначается для интерфейса `Storage`. Если принята команда на чтение, то она отправляется сразу в `Storage`. Если же была принята мутирующая команда, то вызывается метод `waitForWal(inData []byte)` который отправляет данные от команды в `barrier` и ждёт сигнала о том, что можно продолжать исполнение (канал `NotifyDone: make(chan struct{})`).
- `barrier struct`

  Представляет собой горутину. В её задачи входит принимать мутирующие запросы в базу данных и держать их неисполненными до тех пор, пока мы не накопим `WAL_BATCH_SIZE` запросов, либо не пройдёт `WAL_BATCH_TIMEOUT` единиц измерения времени. После достижения одно из двух условий производит запись запросов на диск с использонием `writer` и и синализирует всем ждущим горутинам, что можно продолжать исполнение запроса. Для корректной работы требует инициализации с использвоанием `writer`. После инициализации возвращает канал, куда необходимо писать все запросы (`chan Input`), которые мы хотим записать в wal. Для корретного завершения работы требует выхова метода `Close() error`.
  - `Input struct`

    Данная структура представляет собой данные, полученные от запроса (`Data []byte`). Содержат команду для `Storage`, а также канал для уведомления горутины о том, что можно продолжать исполение (`NotifyDone chan struct{}`).
- `writer struct`

  Записывает данные в wal. Ротирует сегменты wal при достижении ими размера в `WAL_SEG_SIZE`. Имена файлов сегментов начинаются с 1 и представляют собой натуральные числа. Для корретного завершения работы требует выхова метода `Close() error`.
7. Config
- Config struct
  Читает перемнные окружения и инициализаует себя корректными параметрами конфигурации для запуска базы данных. Детальная документация каждого параметра содержится в файле `cmd.go`, а домустимые занчения содержатся в `cmd_test.go`.