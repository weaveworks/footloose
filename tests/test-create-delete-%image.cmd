footloose config create --config %testName.footloose --name %testName --key %testName-key --image quay.io/footloose/%image
footloose create --config %testName.footloose
%out docker ps --format {{.Names}}
footloose delete --config %testName.footloose
%out docker ps --format {{.Names}}
