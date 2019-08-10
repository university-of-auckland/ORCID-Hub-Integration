URL=${1:-https://415mdw939a.execute-api.ap-southeast-2.amazonaws.com/prod/v1/enqueue}

for id in 484378182 477579437 208013283 987654321 8524255 350622514 4306445 ; do
  curl -v -H "Content-Type: application/json" -d "{\"subject\":${id}}" $URL
done




	

