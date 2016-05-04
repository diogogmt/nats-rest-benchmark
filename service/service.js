var express = require('express');
var app = express();

var servers = process.env.NATS_SERVERS && process.env.NATS_SERVERS.split(',') || ['nats://server.nats:4222'];
var port = process.env.HTTP_PORT || 3000;
console.log('servers: ', servers);

var config = {
  servers: servers,
  json: true
};
var nats = require('nats').connect(config);

var items = [
  {
    id: 'id',
    name: 'name'
  }
];

nats.subscribe('pub.sub', 'service', function(msg) {
  // console.log('NATS pub.sub - msg: ', msg);
  var options = msg.options;
  nats.publish(options.reply, {
    items: items,
    options: {
      requestUUID: options.requestUUID
    }
  })
});

nats.subscribe('req.rep', 'service', function (msg, replyTo) {
  // console.log('NATS req.rep - message: ', msg, ' - replyTo: ', replyTo);
  nats.publish(replyTo, { items: items })
});

app.get('/rest', function (req, res) {
  // console.log('HTTP /rest');
  res.status(200).json(items);
});

app.listen(port, function () {
  console.log('Service listening on port ', port);
});