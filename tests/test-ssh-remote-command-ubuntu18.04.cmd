footloose config create --config %t.footloose --name %t --key %t-key --image quay.io/footloose/ubuntu18.04
footloose create --config %t.footloose
%out footloose --config %t.footloose ssh root@node0 hostname
footloose delete --config %t.footloose
