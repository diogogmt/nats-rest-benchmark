var benchrest = require('bench-rest');
var serverIp = 'http://192.168.44.10:3000';

//var qs = '';
//for (var i = 0; i < 2000; i++) {
//  qs += 'qs' + i + '=12345&';
//}

var flow = {
  main: [
    //{ get: serverIp + '/nats/items?' + qs },
    { get: serverIp + '/nats/items'},
    { get: serverIp + '/nats/items/93416ae8-7225-4bb6-bc93-31dd30dc55a6' },
    { get: serverIp + '/nats/items/error' }
  ]
};

var runOptions = {
  limit: 2500,
  iterations: 50000
};

benchrest(flow, runOptions)
.on('error', function (err, ctxName) { console.error('Failed in %s with err: ', ctxName, err); })
.on('end', function (stats, errorCount) {
  console.log('error count: ', errorCount);
  console.log('stats', stats);
});
