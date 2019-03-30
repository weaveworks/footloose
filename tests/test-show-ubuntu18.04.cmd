footloose config create --config %testName.footloose --name %testName --key %testName-key --image quay.io/footloose/ubuntu18.04
footloose create --config %testName.footloose
footloose delete --config %testName.footloose
%out footloose show --config %testName.footloose
%out footloose show -f json --config %testName.footloose
