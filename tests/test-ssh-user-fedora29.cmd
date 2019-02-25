footloose config create --config %t.footloose --name %t --key %t-key --image quay.io/footloose/fedora29
footloose create --config %t.footloose
%out footloose --config %t.footloose ssh root@node0 whoami
footloose delete --config %t.footloose
