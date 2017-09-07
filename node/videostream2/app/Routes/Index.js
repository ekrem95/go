import React, { Component } from 'react';
import io from 'socket.io-client';

export default class Index extends Component {
  constructor() {
    super();
    this.state = { streams: null };
    this.navigate = this.navigate.bind(this);
  }

  componentWillMount() {
    fetch('/streams')
    .then(res => res.json())
    .then(res => {
      if (res.length < 1) {
        this.setState({ streams: null });
      } else {
        this.setState({ streams: res });
      }
    });
  }

  navigate(name) {
    this.props.history.push(('/v/' + name));
  }

  render() {
    return (
      <div>
        <h2>Home</h2>
        <label>Username</label>
        <br />
        <br />
        <input
          onKeyUp={(e) => {
            if (e.keyCode === 13) {
              const name = document.getElementById('name').value;
              this.navigate(name);
            }
          }}

          type="text" id="name"/>
        <button
          onClick={() => {
            const name = document.getElementById('name').value;
            this.navigate(name);
          }}
          >Start</button>
          <br />
          <br />
          {this.state.streams &&
            this.state.streams.map(s => {
              console.log(s);
              return (
                <div
                  key={s}
                  onClick={() => {
                    this.props.history.push('/r/' + s);
                  }}
                  >{s}</div>
              );
            })
          }
      </div>
    );
  }
}
