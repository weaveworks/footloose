footloose config --config %t.footloose --name %t --key %t-key
footloose create --config %t.footloose
%out footloose --config %t.footloose ssh node0 hostname
footloose delete --config %t.footloose
