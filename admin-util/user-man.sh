#!/bin/bash

DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

VIK_DIR=$( cd "${DIR}/../" && pwd )


# CMD_APP="docker exec vikunja_api /app/vikunja/vikunja" # for docker
CMD_APP=${VIK_DIR}/vikunja

# Function to list users
list_users() {
    echo "User list:"
    ${CMD_APP} user list
}

# Function to create a user
create_user() {
    # Prompt the user for the username
    read -p "Enter the username: " username

    # Prompt the user for the password (without showing it on the screen)
    read -sp "Enter the password: " password
    echo

    # Prompt the user for the email address
    read -p "Enter the email address: " email

    # Execute the docker command with the values provided by the user
    ${CMD_APP} user create -e "$email" -p "$password" -u "$username"
}

# Function to delete a user
delete_user() {
    # List users to display the IDs
    list_users

    # Prompt the user for the ID of the user to delete
    read -p "Enter the ID of the user you want to delete: " userID

    # Execute the docker command to delete the user with the provided ID
    ${CMD_APP} user delete "$userID" -c -n
}

# Main function
main() {
    while true; do
        echo "What action would you like to perform?"
        echo "1. List users"
        echo "2. Create user"
        echo "3. Delete user"
        echo "4. Exit"

        read -p "Enter the number of the action you want to perform: " option

        case $option in
            1)
                list_users
                ;;
            2)
                create_user
                ;;
            3)
                delete_user
                ;;
            4)
                echo "Exiting the program."
                exit 0
                ;;
            *)
                echo "Invalid option. Please select a valid option."
                ;;
        esac
    done
}

# Call to the main function
main