# Final_Skillbox
## Задача разработанного сетевого сервиса (предыстория)
Вы пришли работать разработчиком в компанию занимающуюся провайдингом современных средств коммуникации.

Компания предоставляет инструменты и API для автоматизации работы систем SMS, MMS, Голосовых звонков и Email. География клиентов распространяется на 160 стран и компания быстро растёт. Требуется всё больше ресурсов со стороны службы поддержки и было принято решение снизить количество заявок с помощью создания страниц информирования клиентов о текущем состоянии систем. С помощью этих страниц компания планирует снизить количество однотипных вопросов и высвободить время агентов службы поддержки для решения более сложных задач.

В числе прочего были составлены страницы с ответами на часто задаваемые вопросы, уведомления о неполадках и истории инцидентов, чтобы клиенты могли самостоятельно проверять действующие системы на работоспособность. Поскольку компания работает на широкую аудиторию и распространена по всему миру, практически невозможно уследить за всеми изменениями вручную. 

Поэтому каждое подразделение компании самостоятельно контролирует работу поставщиков услуг в автоматизированном режиме храня эти данные. Ваша задача — финализировать проект, объединив эти данные и разработав небольшой сетевой сервис, который будет принимать запросы по сети и возвращать данные о состоянии систем компании. Эти данные будут выводиться на web страницу сайта компании под названием Statuspage и содержать в себе географию и статусы сервисов. Так клиенты смогут проверить свой регион на наличие ошибок прежде чем обращаться в службу поддержки.

## Запуск 
1) В папке final/simulator выполняем команду:
```bash
go run main.go
```
После запуска в директории simuator будут сгенерированы файлы с тестовыми данными, именно их мы и будем забирать нашим сервисом для дальнейшей обработки

```
simulator/sms.data
simulator/voice.data
simulator/email.data
simulator/billing.data
```

Другая часть данных будет доступна нам по http запросу на порт 8383

```
http://127.0.0.1:8383/mms
http://127.0.0.1:8383/support"
http://127.0.0.1:8383/accendent"
```

2) В папке final/aggregator/cmd выполянем команду:
```bash
go run main.go
```
Происходит запуск разработанного сетевого сервиса (агрегатор), который собирает всю информацию. 

3) В папке final/web (открываем index.html)
Для демонстрации конечного результата можно открыть HTML страницу из директории web

```
index.html
```

Для работы с нашим сервисом агрегатором создана директория web, страница index.html делает запрос на http://127.0.0.1:8282/api

По этому адресу сервис отдаёт всю собранную информацию.
## Краткое описание директории aggregator
Основные директории:

1) dataset - настройка map для alpha2;

2) handlers -  7 обработчиков и реализация хранилища;

3) web - обработка http запросов к нашему агрегатору.

Основные файлы:
1) cmd - месторасположение main.go;

2) aggregator.go - все структуры в общем виде;

3) server.go - запуск и выключение сервера.

## Заметки по разработанному сетевому сервису
1) Язык реализации – Go 1.21, сетевой фреймворк - Gin Web Framework;
2) Идет обработка файлов  .data прямо из папки, где их симулятор создал. Поэтому, ничего перемещать не нужно. Достаточно просто запустить main.go;
3) Хранятся данные в структурах;
4) Обновление данных происходит раз в 30 секунд (обход всех .data и генерация запросов, их обработка и сохранение полученной информации в структуры). После обновления данных, сразу же формируется JSON, который всегда готов к отправке;
5) Сервер всегда слушает и готов отправить актуальный в текущие 30 секунд JSON;
6) Обработка всех .data и запросов идет в своих горутинах;
7) Логирование (Logrus);
8) Graceful shutdown.
