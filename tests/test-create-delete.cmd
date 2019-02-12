footloose config create --config %t.footloose --name %t --key %t-key
footloose create --config %t.footloose
%out docker ps --format {{.Names}}
footloose delete --config %t.footloose
%out docker ps --format {{.Names}}
