// React framework
'use strict';

const e = React.createElement;
      
class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = init({ component: this });
  }

  render() {
    return view({ model: this.state, component: this })
  }
  
  msg({ f, args }) {
    console.log(f.name, args)
    const { model, cmd = none } = f({ model: this.state, args, component: this })
    this.setState(_log(model))
    console.log('->', snapshot(model))
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


// Init
function init({ component }) {
  const socket = io('/')



  socket.on('connect', () => {
    console.log('connected')

    socket.emit('chat', 'hi!')
  })

  socket.on('chat', msg => {
    console.log('chat: ', msg)
    try {
      const value = JSON.parse(msg)

      const message = {
        content: '' + value['Msg'],
        socketId: '' + value['SocketID'],
        color: parseInt(value['Color'])
      }

      component.msg({ f: gotChat, args: { message }})
    }
    catch (e) {
      console.error(e)
    }
  
  })


  socket.on('board', msg => {
    // console.log('board: ', msg)
    const board = JSON.parse(msg)
    // console.log(board)
    component.msg({ f: gotBoard, args: { board }})
  })

  window.setTimeout(() => {
    component.msg({ f: gotTick })
  }, 3000)


  return {
    board: null,
    socket,
    items: []
  }
}

// Update
function gotBoard({ model, args: { board }}) {
  model.board = board
  return { model }
}

function toggl({ model, args: { x, y }}) {
  model.socket.emit('toggl', JSON.stringify({ x, y }))
  return { model }
}

function gotTick({ model, component }) {
  model.socket.emit('tick')

  window.setTimeout(() => {
    component.msg({ f: gotTick })
  }, 3000)

  return { model }
}

function gotChat({ model, args: { message }, component }) {
  model.items.push({
    id: uuidv4(),
    content: message.content,
    socketId: message.socketId,
    color: message.color
  })

  return { model }
}

function chatInput({ model, args: { value }}) {
  model.newChat = value

  return { model }
}
function chatSubmit({ model }) {
  if (model.newChat) {
    model.socket.emit('chat', model.newChat)
    model.newChat = ''
  }

  return { model }
}

// View
function view({ model, component }) {
  if (model.board) {
    return div({}, [
      viewBoard({ model, board: model.board, component }),
      div({}, [
        model.items
        .slice(-10)
        .map(item => viewItem({ model, item }))
      ]),
      form({
        onSubmit: e => {
          e.preventDefault()
          component.msg({ f: chatSubmit })
        },
        className: 'flex'
      }, [
        input({
          onInput: e => component.msg({ f: chatInput, args: { value: e.target.value }}),
          className: 'p-1 shadow',
          value: model.newChat
        }),
        button({
          className: 'rounded p-2 bg-green-200 shadow'
        }, [ text('+') ])
      ])

    ])
  }
  else {
    return text('Loading...')
  }
}

function viewBoard({ model, board, component }) {
  return board.map((line, y) => viewLine({ model, line, y, component }))
}

function viewLine({ model, line, y, component }) {
  return div({
    className: 'flex'
  }, line.map((cell, x) => viewCell({ model, cell, x, y, component })))
}

function viewCell({ model, cell, x, y, component }) {
  const style = cell !== -1 ? {
    backgroundColor: `hsl(${cell}, 50%, 50%)`
  } : {}

  return div({
    style,
    className: `border w-8 h-8`,
    onClick: () => component.msg({ f: toggl, args: { x, y } } )
  }, [])
}

function viewItem({ model, item }) {
  return div({
    style: {
      color: `hsl(${item.color}, 50%, 30%)`
    }
  }, [
    text(item.content)
  ])
}




// Virtual DOM
function div(attributes, children) {
  return React.createElement('div', attributes, children)
}

function text(str) {
  return React.createElement('span', null, str)
}

function button(attributes, children) {
  return React.createElement('button', attributes, children)
}


function form(attributes, children) {
  return React.createElement('form', attributes, children)
}


function input(attributes, children) {
  return React.createElement('input', attributes, children)
}
