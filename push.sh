#!/bin/sh

commit_msg="$1"

if [ -z "$commit_msg" ]
then
    printf "Commit message not provided\n"
    exit
fi

printf "Testing ado help function: "
go_test=$(go test)
fail_count=$(echo $go_test | grep "FAIL" | wc -l)

if [ $fail_count -ne 0 ]
then
    printf "FAIL\n"
    printf "$go_test\n"
    exit
fi
printf "DONE\n"

printf "Building binaries: "
go build -o ./bin/ado
GOOS=windows GOARCH=amd64 go build -o ./bin/ado.exe
printf "DONE\n"

printf "pushing to git: "
git add .
git commit -m "$commit_msg"
git push
printf "DONE\n"
