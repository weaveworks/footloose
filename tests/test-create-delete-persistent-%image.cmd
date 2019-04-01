footloose config create --config %testName.footloose --name %testName --key %testName-key --image quay.io/footloose/%image
footloose create --config %testName.footloose
%out docker ps --format {{.Names}}
%out docker inspect %testName-node0 -f "{{.HostConfig.AutoRemove}}"
footloose delete --config %testName.footloose
%out docker ps --format {{.Names}}
