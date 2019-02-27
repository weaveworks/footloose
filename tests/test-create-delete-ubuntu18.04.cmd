footloose config create --config %t.footloose --name %t --key %t-key --image quay.io/footloose/ubuntu18.04
footloose create --config %t.footloose
%out docker ps --format {{.Names}}
footloose delete --config %t.footloose
%out docker ps --format {{.Names}}
