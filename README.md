# SSH Machine Controller

A machine controller for the [Cluster API](http://sigs.k8s.io/cluster-api) reference implementation.


## Testing

### Machine Actuator

The machine actuator tests will run against a mock SSH server, but until then they run against an actual host, and require valid SSH credentials. Generate them before running tests:

    $ ./machine/testdata/generate-sshcredentials-secret.sh username /path/to/private-ssh-key > ./machine/testdata/sshcredentials-secret.yaml