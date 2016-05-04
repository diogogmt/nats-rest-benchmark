var express = require('express');
var app = express();
var request = require('request');
var uuid = require('node-uuid');

var natsServers = process.env.NATS_SERVERS && process.env.NATS_SERVERS.split(',') || ['nats://server.nats:4222'];
var serviceUrl = process.env.SERVICE_URL || 'http://service01:3000';
var port = process.env.HTTP_PORT || 8000;
console.log('natsServers: ', natsServers);
var config = { 
  servers: natsServers,
  json: true
};
var nats = require('nats').connect(config);

var pendingRequests = {};

var reply = 'proxy';

nats.subscribe(reply, 'proxy', function(msg) {
  var options = msg.options;
  pendingRequests[options.requestUUID].res.json(msg);
  delete pendingRequests[options.requestUUID];
});

app.get('/nats/pub/sub', function (req, res) {
  // console.log('NATS /pub/sub');
  var requestUUID = uuid.v4();
  var options = {
    requestUUID: requestUUID,
    reply: reply
  };
  pendingRequests[requestUUID] = {
    req: req,
    res: res
  };
  nats.publish('pub.sub', {
    options: options,
    query: req.query
  });
});

app.get('/nats/req/rep', function (req, res) {
  // console.log('NATS /req/rep');
  nats.request('req.rep', { data: true }, { max: 1 }, function(data) {
    res.status(200).json(data);
  });
});


app.get('/rest', function (req, res) {
  // console.log('HTTP /rest');
  request(serviceUrl + '/rest', function (error, response, body) {
    return res.status(200).json(body);
  })
});

app.listen(port, function () {
  console.log('Service listening on port ', port);
});