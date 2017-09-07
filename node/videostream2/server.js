const express = require('express');
const app = express();
const path = require('path');
const bodyParser = require('body-parser');
const server = require('http').createServer(app);
const io = require('socket.io').listen(server, { pingTimeout: 30000 });
const morgan = require('morgan');
const compression = require('compression');
const redis = require('redis');
const client = redis.createClient();

app.use(compression());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));

// app.use(morgan('tiny'));

app.use(express.static(path.join(__dirname, 'public')));

app.use('/', require('./router.js'));

server.listen(process.env.PORT || 3000);

// Connect

client.on('error', function (err) {
    console.log('Error ' + err);
  });

client.on('connect', function () {
    console.log('Connected to Redis');
  });

client.del(['test_online_users'], (err, reply) => {
  if (err)return err;
  console.log(reply);
});

let history = {};

io.on('connection', socket => {
  console.log('Connected');

  let user;

  socket.on('new_stream', username => {
    console.log(username);
    user = username;

    client.sadd(['test_online_users', user], (err, reply) => {
      if (err) return err;
      console.log(reply);
    });
  });

  socket.on('stream', data => {
    // socket.emit('dist', data);
    socket.broadcast.emit('dist' + data[1], data[0]);
  });

  socket.on('disconnect', () => {

    console.log(user + ' exited');

    if (user) {
      client.srem(['test_online_users', user], (err, reply) => {
        if (err) return err;
        console.log(reply);
      });
    }

  });

});
