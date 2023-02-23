# 🏊 SwimLogs 🌊

A web application for creating and sharing swim trainings in a swimming team.

## To fix

1. Migrate getting data from .vugu files to state package
2. migrate from zap to zerolog
3. Also get stacktrace for error in client logging
4. Add unique constraint to session, [like this](https://stackoverflow.com/a/14221810)
5. When adding new session/training mark saved collection of them as dirty and try to reload them, on failure just keep what is there
6. Change CORS! main.go

## TODO

1. Frontend -> Go Vugu / Hugo
2. Zitadel
3. Caddy
4. Docker-compose

## Known issues

1. couldn't implement constraint into Postgres for checking if training.date is
   actually the training.day day of the week. This was my attempt:
   `constraint training_date_day_check check (lower(to_char(date::date, 'Day')) = day::text)`

## Useful links

- Read and apply [Owasp Top 10](https://owasp.org/www-project-top-ten/)
