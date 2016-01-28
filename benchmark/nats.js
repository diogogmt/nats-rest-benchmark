var benchrest = require('bench-rest');
var serverIp = 'http://192.168.44.10:3000';


// OR more powerfully define an array of REST operations with substitution 
// This does a unique PUT and then a GET for each iteration 
var flow = {
  main: [
    { get: serverIp + '/nats/items' },
    { get: serverIp + '/nats/items/93416ae8-7225-4bb6-bc93-31dd30dc55a6' },
    { get: serverIp + '/nats/items/error' },
  ]
};

var runOptions = {
  limit: 1024,
  iterations: 10240
};

benchrest(flow, runOptions)
.on('error', function (err, ctxName) { console.error('Failed in %s with err: ', ctxName, err); })
.on('end', function (stats, errorCount) {
  console.log('error count: ', errorCount);
  console.log('stats', stats);
});
