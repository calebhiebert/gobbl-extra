for D in *; do
    if [ -d "${D}" ]; then
        cd "${D}"
        echo "UPDATING ${D}"
        go get -u
        go mod tidy
        cd ..
    fi
done