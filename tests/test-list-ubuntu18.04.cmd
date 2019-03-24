footloose config create --config %testName.footloose --name %testName --key %testName-key --image quay.io/footloose/ubuntu18.04
footloose create --config %testName.footloose
footloose delete --config %testName.footloose
%out footloose list --config %testName.footloose
%out footloose list -f json --config %testName.footloose
