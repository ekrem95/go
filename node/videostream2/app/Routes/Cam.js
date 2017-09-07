import React, { Component } from 'react';
import io from 'socket.io-client';

export default class Cam extends Component {
  componentWillMount() {

    let socket = io.connect('/');

    if (socket !== undefined) {
      socket.on('dist', data => {
        console.log('dist');
        const img = document.getElementById('play');
        img.src = data;
      });
    }
  }

  render() {
    return (
      <div>
        <h2>Cam</h2>
        <img id="play"/>
      </div>
    );
  }
}
