footloose config create --override --config %testName.footloose --name %testName --key %testName-key --image quay.io/footloose/ubuntu18.04
footloose create --config %testName.footloose
footloose delete --config %testName.footloose
%out footloose show --config %testName.footloose
%out footloose show -o json --config %testName.footloose
