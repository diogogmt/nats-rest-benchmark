var express = require('express');
var app = express();
var request = require('request');
var uuid = require('node-uuid');

var servers = [ 'nats://nats01.server.local:4222', 'nats://nats02.server.local:4223', 'nats://nats03.server.local:4224' ];
var nats = require('nats').connect({ servers: servers });

var pendingRequests = {};

var reply = 'proxy';

nats.subscribe(reply, function(msg) {
  msg = JSON.parse(msg);
  var data = JSON.parse(msg.data);
  pendingRequests[msg.requestUUID].res.json(data);
  delete pendingRequests[msg.requestUUID];
});

app.get('/nats/items', function (req, res) {
  var requestUUID = uuid.v4();
  var payload = {
    requestUUID: requestUUID,
    reply: reply
  }
  pendingRequests[requestUUID] = {
    req: req,
    res: res
  }
  nats.publish('listItems', JSON.stringify(payload));
});

app.get('/nats/items/:id', function (req, res) {
  var itemId = req.params.id;
  var item = {
    id: itemId
  }
  var requestUUID = uuid.v4();
  var payload = {
    requestUUID: requestUUID,
    reply: reply,
    data: JSON.stringify(item)
  }
  pendingRequests[requestUUID] = {
    req: req,
    res: res
  }
  nats.publish('itemDetails', JSON.stringify(payload));
});

app.get('/rest/items', function (req, res) {
  request('http://go.client.local:4000/items', function (error, response, body) {
    res.json(body)
  })
});

app.get('/rest/items/:id', function (req, res) {
  request('http://go.client.local:4000/items/' + req.params.id, function (error, response, body) {
    res.json(body)
  })
});

app.listen(3000, function () {
  console.log('Example app listening on port 3000!');
});


