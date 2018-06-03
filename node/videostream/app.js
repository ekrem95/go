const express = require('express');
const app = express();
const http = require('http').Server(app);
const io = require('socket.io')(http);

const Log = require('log');
const log = new Log('debug');

const port = process.env.PORT || 3000;

app.use(express.static(__dirname + '/public'));
app.use('/io', express.static(__dirname + '/node_modules/socket.io-client/dist/'));

app.get('/', (req, res) => res.redirect('index.html'));
app.get('/video', (req, res) => res.redirect('view.html'));

io.on('connection', socket => {
  socket.on('stream', image => socket.broadcast.emit('stream', image));
});

http.listen(port, () => log.info('Listening'));
