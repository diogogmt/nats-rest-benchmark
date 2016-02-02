var express = require('express');
var app = express();
var request = require('request');
var uuid = require('node-uuid');

// Cluster
//var servers = [ 'nats://nats01.cluster.local:4222', 'nats://nats02.cluster.local:4223', 'nats://nats03.cluster.local:4224' ];
var servers = [ 'nats://nats.server.local:4222' ];
var config = { 
  servers: servers,
  json: true,
};
var nats = require('nats').connect(config);

var pendingRequests = {};

var reply = 'proxy';

nats.subscribe(reply, function(msg) {
  var options = msg.options;
  pendingRequests[options.requestUUID].res.json(msg);
  delete pendingRequests[options.requestUUID];
});

app.get('/nats/items', function (req, res) {
  var requestUUID = uuid.v4();
  var options = {
    requestUUID: requestUUID,
    reply: reply
  };
  var payload = {
    options: options,
    query: req.query
  };
  pendingRequests[requestUUID] = {
    req: req,
    res: res
  };
  nats.publish('listItems', payload);
});

app.get('/nats/items/:id', function (req, res) {
  var itemId = req.params.id;
  var requestUUID = uuid.v4();
  var options = {
    requestUUID: requestUUID,
    reply: reply
  };
  var payload = {
    id: itemId,
    options: options
  };
  pendingRequests[requestUUID] = {
    req: req,
    res: res
  };
  nats.publish('itemDetails', payload);
});

app.get('/rest/items', function (req, res) {
  request('http://go.client.local:4000/items?' + req._parsedUrl.query, function (error, response, body) {
    try {
      body = JSON.parse(body);
    } catch (e) {
      console.log('Error parsing body:', e);
      return res.json(400);
    }
    return res.json(body);
  })
});

app.get('/rest/items/:id', function (req, res) {
  request('http://go.client.local:4000/items/' + req.params.id, function (error, response, body) {
    try {
      body = JSON.parse(body);
    } catch (e) {
      console.log('Error parsing body:', e);
      return res.json(400);
    }
    return res.json(body);
  })
});

app.listen(3000, function () {
  console.log('Example app listening on port 3000!');
});