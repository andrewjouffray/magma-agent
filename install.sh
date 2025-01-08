# /bin/sh

# verify if the user is root
if [ "$(id -u)" != "0" ]; then
    echo "This script must be run as root, magma is intended to track files accross the entire system" 1>&2
    exit 1
fi

# check for the ./build/magma binary file
if [ ! -f ./build/magma ]; then
    echo "The magma binary file is missing, please build the project first" 1>&2
    exit 1
fi

# copy the magma binary file to /usr/bin
cp ./build/magma /usr/bin/magma

# check if the copy was successful
if [ ! -f /usr/bin/magma ]; then
    echo "Failed to copy the magma binary file to /usr/bin" 1>&2
    exit 1
fi

# set the permissions for the magma binary file
chmod 755 /usr/bin/magma

# check if the permissions were set successfully
if [ "$(stat -c %a /usr/bin/magma)" != "755" ]; then
    echo "Failed to set the permissions for the magma binary file" 1>&2
    exit 1
fi


# run magma init
magma init

# check if ssh is installed and enabled
if ! command -v ssh &> /dev/null; then
    echo "The ssh command is not installed, please install it if you want remote access capabilities to be enabled" 1>&2
    echo "The magma binary file was successfully installed, but remote access capabilities are disabled"

    else
        echo "The magma binary file was successfully installed"
fi

exit 0

