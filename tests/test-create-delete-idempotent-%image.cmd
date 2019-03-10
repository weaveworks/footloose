footloose config create --config %testName.footloose --name %testName --key %testName-key --image quay.io/footloose/%image
footloose create --config %testName.footloose
%out footloose create --config %testName.footloose echo success
footloose delete --config %testName.footloose
%out footloose delete --config %testName.footloose echo success
