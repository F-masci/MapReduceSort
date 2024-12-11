:: Master
start /B go run master.go --idx 0
start /B go run master.go --idx 1

:: Mapper
start /B go run worker.go --address localhost --proto tcp --port 45980 --map
start /B go run worker.go --address localhost --proto tcp --port 45981 --map

:: Reducer
start /B go run worker.go --address localhost --proto tcp --port 45989 --reduce
start /B go run worker.go --address localhost --proto tcp --port 45990 --reduce

:: Client
:: start /B go run client.go --client fmasci --master-idx 0