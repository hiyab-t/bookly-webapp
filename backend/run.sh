#!/bin/bash

SKIP_PROMPTS=false

# checks if any arguments have been passed

if [ "$*" == "" ]; then
        go run ./cmd/api -help
        echo 'The above arguments exist if you would like to set your own settings '
        echo '(Recommendation: leave the API server port and cors-trusted-origins as default)'
        echo ""
        echo 'would you like to set your own setting?(y/N)'

        read setArgQ

        # if user is going by default settings it continues to excute the rest of commands
        if [[ "$setArgQ" != "y" ]]; then
            echo ""
        else
            echo 'Please run "./run.sh" with the settings argument you would like to set'
            exit 0
        fi
else [ "$*" == "-f" ]
    SKIP_PROMPTS=true
fi

if [ "$SKIP_PROMPTS" != "true" ]; then
    echo 'Are you sure you want to remove all running containers and volumes? (y/N)'
    echo 'If N, please manually remove any postgres running container to avoid errors.'

    read dckerDown

    if [[ "$dckerDown" != "y" ]]; then
        echo ""
    else
        echo '....stopping all running containers....'
        docker compose down -v
    fi
else
    docker compose up -d
fi

echo '....creating .env file to hold your secrets....'
touch .env

value=${input#*-f}
echo ARGS=$value > .env

if [ "$SKIP_PROMPTS" != "true" ]; then
    echo 'Enter postgresDB name:'
    read POSTGRES_DB
    echo POSTGRES_DB=$POSTGRES_DB >> .env

    echo 'Enter postgresDB username:'
    read POSTGRES_USER
    echo POSTGRES_USER=$POSTGRES_USER >> .env

    echo 'Enter postgresDB password:'
    read POSTGRES_PASSWORD
    echo POSTGRES_PASSWORD=$POSTGRES_PASSWORD >> .env

    echo DATABASE_URL="'postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:5432/$POSTGRES_DB?sslmode=disable'" >> .env
else
    # adding the env vars for clarity
    echo -e "POSTGRES_DB=bookly_db \nPOSTGRES_USER=webuser \nPOSTGRES_PASSWORD='444'" >> .env
    echo DATABASE_URL='postgres://webuser:444@localhost:5432/bookly_db?sslmode=disable' >> .env
fi

# to start browser 
echo ORIGIN=http://127.0.0.1:5500 >> .env

echo 'all your secrets are saved'

echo 'Running up migrations... '

make db/migrations/up 

xdg-open http://127.0.0.1:5500/ui/index.html || open http://127.0.0.1:5500/ui/index.html || start http://127.0.0.1:5500/ui/index.html

make run/api