footloose config create --config %testName.footloose --name %testName --key %testName-key --image quay.io/footloose/%image
%out footloose config get --config %testName.footloose "machines[0].spec"
