import React, { Component } from 'react';
import io from 'socket.io-client';

export default class Sender extends Component {
  constructor() {
    super();
    this.state = { username: null };
    this.setup = this.setup.bind(this);
  }

  componentDidMount() {
    this.setup();

  }

  setup() {
    let canvas = document.getElementById('preview');
    let context = canvas.getContext('2d');
    canvas.width = 800;
    canvas.height = 600;
    context.width = canvas.width;
    context.height = canvas.height;
    const video = document.getElementById('video');

    const username = this.props.location.pathname.split('/').pop();
    this.setState({ username });

    let socket = io.connect('/');

    if (socket !== undefined) {
      function loadCam(stream) {
        video.srcObject = stream;
        console.log('Camera connected.');
      }

      function loadFail() {
        console.log('Camera not connected.');
      }

      function viewVideo(video, context) {
        context.drawImage(video, 0, 0, context.width, context.height);
        socket.emit('stream',
          [
            canvas.toDataURL('image/webp'),
            username,
          ]
        );
      }

      (function () {
        navigator.getUserMedia = (
          navigator.getUserMedia || navigator.webkitGetUserMedia ||
          navigator.mozGetUserMedia || navigator.msgGetUserMedia);
        if (navigator.getUserMedia) {
          navigator.getUserMedia({ video: true }, loadCam, loadFail);
        }

        setInterval(() => {
          viewVideo(video, context);
        }, 20);
      }());

      socket.emit('new_stream', username);

    }
  }

  render() {
    return (
      <div>
        <h2>My name is {this.state.username}</h2>
        <video src="" id="video" autoPlay="true"></video>
        <canvas id="preview"></canvas>
      </div>
    );
  }
}
