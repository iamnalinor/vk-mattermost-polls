# vk-mattermost-polls

Решение тестового задания от компании VK на стажировку по направлению Backend.

Это бот для мессенджера Mattermost, который позволяет создавать и проводить опросы. В качестве БД используется Tarantool.

## Запуск

1. Скопируйте `.env.dist` в `.env`.
2. Укажите URL сервера Mattermost и токен бота (на [developers.mattermost.com](https://developers.mattermost.com/integrate/reference/bot-accounts/) описано, как его получить).
Обратите внимание: если Mattermost запущен локально (на том же хосте, где находится Docker), то в URL вместо `localhost` нужно указать `host.docker.internal`.
3. Запустите команду `docker compose up -d`. Если Mattermost запущен локально, то нужно запустить `docker compose -f docker-compose.yml -f docker-compose.local.yml up -d`.

Бот будет работать в личных сообщениях и в любом канале, куда его добавят.

## Команды

- `!newpoll` - создать опрос
- `!poll <id>` - просмотр опроса
- `!vote <poll id> <option number>` - проголосовать
- `!closepoll <id>` - закрыть опрос, в нём нельзя будет проголосовать
- `!deletepoll <id>` - удалить опрос
