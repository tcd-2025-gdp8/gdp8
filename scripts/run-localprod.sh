#! /bin/sh

docker compose -f ./scripts/compose-localprod.yml --project-directory . up --build
