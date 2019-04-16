# Test that common utilities are present in the base images
footloose config create --config %testName.footloose --override --name %testName --key %testName-key --image quay.io/footloose/%image
footloose create --config %testName.footloose
footloose --config %testName.footloose ssh root@node0 hostname
footloose --config %testName.footloose ssh root@node0 ps
footloose --config %testName.footloose ssh root@node0 ifconfig
footloose --config %testName.footloose ssh root@node0 ip route
footloose --config %testName.footloose ssh root@node0 -- netstat -n -l
footloose --config %testName.footloose ssh root@node0 -- ping -V
footloose --config %testName.footloose ssh root@node0 -- curl --version
footloose --config %testName.footloose ssh root@node0 -- wget --version
footloose --config %testName.footloose ssh root@node0 -- vi --help
footloose --config %testName.footloose ssh root@node0 -- sudo true
footloose delete --config %testName.footloose
