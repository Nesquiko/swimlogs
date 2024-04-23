# Swimlogs

## Overall TODOS

1. Fix CI/CD, I want to:
   - be able to deploy frontend without backend
   - run frontend tests if frontend changed
   - run backend tests if backend changed
   - only deploy backend on tagged commits
   - run tests on pull request
   - run tests on merge to main
2. Update versions of everything

## TODOS frontend

1. add/accept buttons to upper right and cancel to upper left?
2. duplicate training
3. floating action button with popup menu for adding set, superset, pyramid
4. parachute, other equipment
5. set intensity slider
6. set style, make only usable when really necessary
7. dragable sets in preview
8. multi select sets for delete/duplicate (and move?)
9. set groups
10. graphs, or some other preview of swam distance
11. rework history
12. Not found page in frontend
13. One method through which requests go for better monitoring?
14. Tests?
15. PWA

## TODOS backend

1. [Makefile](https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects)
2. [read this](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/)
3. [generic handler functions](https://www.willem.dev/articles/generic-http-handlers/)
4. uptimerobot endpoint
5. docker scan

## Cleanup

1. DB types in lowercase so I can I don't have to lower them manually in frontend
   when accessing translations
