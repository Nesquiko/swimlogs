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

1. training distance split by group
2. add/accept buttons to upper right and cancel to upper left?
3. parachute, other equipment
4. set intensity slider
5. set style, make only usable when really necessary
6. dragable sets in preview
7. multi select sets for delete/duplicate (and move?)
8. set groups
9. graphs, or some other preview of swam distance
10. rework history
11. Not found page in frontend
12. One method through which requests go for better monitoring?
13. Tests?
14. PWA

## TODOS backend

1. [read this](https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/)
2. [generic handler functions](https://www.willem.dev/articles/generic-http-handlers/)
3. uptimerobot endpoint
4. docker scan

## Cleanup

1. DB types in lowercase so I can I don't have to lower them manually in frontend
   when accessing translations
