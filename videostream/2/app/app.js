import React, { Component } from 'react';
import { render } from 'react-dom';
import { BrowserRouter, Route, Switch } from 'react-router-dom';

import Index from './Routes/Index';
import Cam from './Routes/Cam';
import Sender from './Routes/Sender';

import style from './style.scss';

class App extends Component {
  render() {
    return (
      <BrowserRouter>
        <Switch>
          <Route exact path="/" component={Index}/>
          <Route path="/r/:name" component={Cam}/>
          <Route path="/v/:name" component={Sender}/>
        </Switch>
      </BrowserRouter>
    );
  }
}

render(<App />, document.getElementById('app'));
