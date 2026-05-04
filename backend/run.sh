#!/bin/bash


input="$*"

value=${input#*-cors-trusted-origins}
echo $value

OtherArgs=${value#*-}
cut 
echo $value


# checks if any arguments have been passed
while [[ 1 = 1 ]]; do
    if [ "$input" != "" ]; then
        break
    elif [ "$input" = "" ]; then
        go run ./cmd/api -help
        echo 'The above arguments exist if you would like to set your own settings '
        echo '(Recommendation: leave the API server port and cors-trusted-origins as default)'
        edcho ""
        echo 'would you like to set your own setting?(y/N)'

        read setArgQ

        # if user is going by default settings it continues to excute the rest of commands
        if [[ "$setArgQ" != "y" ]]; then
            break
        else
            echo 'Please run "./run.sh" with the settings argument you would like to set'
            exit 0
        fi
    fi
done

echo 'Are you sure you want to remove all running containers and volumes? (y/N)'
echo 'If N, please manually remove any postgres running container to avoid errors.'

read dckerDown

if [[ "$dckerDown" != "y" ]]; then
    echo ""
else
    echo '....stopping all running containers....'
    docker compose down -v
fi


echo '....creating .env file to hold your secrets....'
touch .env

# credentials for db connection
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

# all arguments to be passed to go command in makefile rule 
echo ARGS=$input > .env

# to start browser 

DEFAULT_ORIGIN=http://127.0.0.1:5500
echo ORIGIN=$DEFAULT_ORIGIN >> .env

echo 'all your secrets are saved'

echo 'Running up migrations... '

make db/migrations/up 

if [[ $input =~ "-cors-trusted-origins" ]]; then
    value=${input#*-cors-trusted-origins}
    OtherArgs=${value#*-}
    echo $value | grep --exclude "-$OtherArgs"
    xdg-open $value/ui/index.html || open $DEFAULT_ORIGIN/ui/index.html || start $DEFAULT_ORIGIN/ui/index.html
else 
    xdg-open $DEFAULT_ORIGIN/ui/index.html || open $DEFAULT_ORIGIN/ui/index.html || start $DEFAULT_ORIGIN/ui/index.html
fi


make run/api




