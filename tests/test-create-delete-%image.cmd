footloose config create --config %t-%image.footloose --name %t --key %t-%image-key --image quay.io/footloose/%image
footloose create --config %t-%image.footloose
%out docker ps --format {{.Names}}
footloose delete --config %t-%image.footloose
%out docker ps --format {{.Names}}
