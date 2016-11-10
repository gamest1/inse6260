export MGO_HOSTS=localhost:27017
export MGO_DATABASE=inse6260
export MGO_WATCHDB=local
export MGO_WATCHCOLLECTION=oplog.rs
export MGO_USERNAME=guest
export MGO_PASSWORD=welcome
export BUOY_DATABASE=inse6260
export USERS_DATABASE=inse6260
export USERS_COLLECTION=users
export REQUESTS_DATABASE=inse6260
export REQUESTS_COLLECTION=service_requests

cd $GOPATH/src/github.com/goinggo/beego-mgo
go clean -i
go build

./beego-mgo
