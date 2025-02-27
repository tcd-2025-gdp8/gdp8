#! /bin/sh

docker compose -f ./scripts/compose-prod.yml --project-directory . up --build
