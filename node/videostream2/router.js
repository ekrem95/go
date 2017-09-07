const express = require('express');
const router = express.Router();
const path = require('path');
const redis = require('redis');
const client = redis.createClient();

client.on('error', function (err) {
    console.log('Error ' + err);
  });

router.get('/streams', (req, res) => {

    client.smembers('test_online_users', (err, reply) => {
      if (err)return err;

      // res.json({ streams: JSON.parse(reply) });
      res.send(reply);
    });
  });

router.get('*', (req, res) => {
  res.sendFile(path.join(__dirname, './template', 'index.html'));
});

module.exports = router;
