var express = require('express');
var app = express();
var request = require('request');

var servers = [ 'nats://nats.server.local:4222' ];
var nats = require('nats').connect({ servers: servers });

app.get('/nats/items', function (req, res) {
  nats.request('listItems', null, function(response) {
    res.json(JSON.parse(response));
  });
});

app.get('/nats/items/:id', function (req, res) {
  var itemId = req.params.id;
  var item = {
    id: itemId
  }
  nats.request('itemDetails', JSON.stringify(item), function(response) {
    res.json(JSON.parse(response));
  });
});

app.get('/rest/items', function (req, res) {
  request('http://go.client.local:4000/items', function (error, response, body) {
    if (err)
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


