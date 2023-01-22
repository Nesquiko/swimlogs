# 🏊 SwimLogs 🌊

A web application for creating and sharing swim trainings in a swimming team.

## To fix

1. In client, create a layer in between a views and api, which caches results
   for later user, for example on landing page, caches training details
2. Change CORS! main.go

## TODO

1. Frontend -> Go Vugu / Hugo
2. Zitadel
3. Caddy
4. Docker-compose

## Known issues

1. there is no locking mechanism on tables, so if two users update same row
   they will overwrite each other
2. couldn't implement constraint into Postgres for checking if training.date is
   actually the training.day day of the week. This was my attempt:
   `constraint training_date_day_check check (lower(to_char(date::date, 'Day')) = day::text)`

## Useful links

- Read and apply [Owasp Top 10](https://owasp.org/www-project-top-ten/)
