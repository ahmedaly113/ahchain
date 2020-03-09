#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/ahchaind/${BINARY:-ahchaind}
ID=${ID:-0}
LOG=${LOG:-ahchaind.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'ahchaind' E.g.: -e BINARY=ahchaind_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export ahchainDHOME="/ahchaind/node${ID}/ahchaind"

if [ -d "`dirname ${ahchainDHOME}/${LOG}`" ]; then
  "$BINARY" --home "$ahchainDHOME" "$@" | tee "${ahchainDHOME}/${LOG}"
else
  "$BINARY" --home "$ahchainDHOME" "$@"
fi

chmod 777 -R /ahchaind

