import { init, view } from './App.js'

// React framework
'use strict';

const e = React.createElement;
      
class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = init({ msg: (f, args) => { this.update(f, args) } });
  }

  render() {
    return view({ model: this.state, msg: (f, args) => { this.update(f, args) } })
  }
  
  update(f, args) {
    const { model, cmd = none } = f({ model: this.state, args, msg: (f, args) => { this.update(f, args) } })
    this.setState(model)
    cmd()
    // console.log(f.name, args)
    // console.log('->', snapshot(model))
  }
}

const domContainer = document.querySelector('#app');
ReactDOM.render(e(App), domContainer);


function _log(v) {
  console.log(v)
  return v
}
function none() {}

function snapshot(model) {
  try {
    return JSON.parse(JSON.stringify(model))
  }
  catch (e) {
    return {}
  }
}



// App


