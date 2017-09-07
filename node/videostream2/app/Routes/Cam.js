import React, { Component } from 'react';
import io from 'socket.io-client';

export default class Cam extends Component {
  constructor() {
    super();
    this.state = { username: null };
  }

  componentWillMount() {
    const username = this.props.location.pathname.split('/').pop();
    this.setState({ username });

    let socket = io.connect('/');

    if (socket !== undefined) {
      socket.on('dist' + username, data => {
        const img = document.getElementById('play');
        img.src = data;
      });
    }
  }

  render() {
    return (
      <div>
        <h2>{this.state.username}'s room</h2>
        <img id="play"/>
      </div>
    );
  }
}
