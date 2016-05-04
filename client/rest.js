var benchRest = require('bench-rest');
var proxyUrl = process.env.PROXY_URL || 'http://proxy01:8000';

var flow = {
  main: [
    { get: proxyUrl + '/rest'}
  ]
};

var runOptions = {
  limit: 10000,
  iterations: 50000
};

benchRest(flow, runOptions)
.on('error', function (err, ctxName) { console.error('Failed in %s with err: ', ctxName, err); })
.on('end', function (stats, errorCount) {
  console.log('error count: ', errorCount);
  console.log('stats', stats);
});
