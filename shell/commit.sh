
git pull

sleep 1


rm shell/chain/chain 

rm shell/api/api 

rm shell/midapi/midapi 

rm shell/midware/midware 

git add --all .
git commit -m"x"
git push origin master-evm:master-evm

# git config credential.helper store
