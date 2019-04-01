footloose config create --config %testName.footloose --name %testName --key %testName-key --image quay.io/footloose/%image
footloose create --config %testName.footloose
%out docker ps --format {{.Names}}
footloose stop --config %testName.footloose
%out docker inspect %testName-node0 -f "{{.State.Running}}"
footloose start --config %testName.footloose
%out docker inspect %testName-node0 -f "{{.State.Running}}"
footloose delete --config %testName.footloose
