# Star Wars
## Run App
>```
>docker-compose up
>```


## Allow Methods

### GetAll
>```
>curl --location -g --request GET 'http://localhost:8080/planets?name='
>```



### GetOne
>```
>curl --location -g --request GET 'http://localhost:8080/planets/{ID}'
>```



### Create
>```
>curl --location --request POST 'http://localhost:8080/planets/' \
>  --header 'Content-Type: application/json' \
>  --data-raw '{
>    "Name": "Planet",
>    "Climate": "Climate",
>    "Terrain": "Terrain"
>}'
>```

### Delete
>```
>curl --location -g --request DELETE 'http://localhost:8080/planets/{ID}'
>```

