// React framework
'use strict';

const e = React.createElement;
      
class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = init({ msg: ({ f, args }) => { this.update({ f, args }) } });
  }

  render() {
    return view({ model: this.state, msg: ({ f, args }) => { this.update({ f, args }) } })
  }
  
  update({ f, args }) {
    const { model, cmd = none } = f({ model: this.state, args, msg: ({ f, args }) => { this.update({ f, args }) } })
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


// Init
function init({ msg }) {
  const socket = io('/')



  socket.on('connect', () => {
    console.log('connected')

    socket.emit('chat', 'hi!')
  })

  socket.on('chat', value => {
    // console.log('chat: ', value)
    try {
      const parsedValue = JSON.parse(value)

      const message = {
        content: '' + parsedValue['Msg'],
        color: parseInt(parsedValue['Color'])
      }

      msg({ f: gotChat, args: { message }})
    }
    catch (e) {
      console.error(e)
    }
  
  })


  socket.on('board', value => {
    const board = JSON.parse(value)
    msg({ f: gotBoard, args: { board }})
  })

  window.setTimeout(() => {
    msg({ f: gotTick })
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

function gotTick({ model, msg }) {
  model.socket.emit('tick')

  window.setTimeout(() => {
    msg({ f: gotTick })
  }, 3000)

  return { model }
}

function gotChat({ model, args: { message } }) {
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
function view({ model, msg }) {
  if (model.board) {
    return div({}, [
      viewBoard({ model, board: model.board, msg }),
      div({}, [
        model.items
        .slice(-10)
        .map(item => viewItem({ model, item }))
      ]),
      form({
        onSubmit: e => {
          e.preventDefault()
          msg({ f: chatSubmit })
        },
        className: 'flex'
      }, [
        input({
          onInput: e => msg({ f: chatInput, args: { value: e.target.value }}),
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

function viewBoard({ model, board, msg }) {
  return board.map((line, y) => viewLine({ model, line, y, msg }))
}

function viewLine({ model, line, y, msg }) {
  return div({
    className: 'flex'
  }, line.map((cell, x) => viewCell({ model, cell, x, y, msg })))
}

function viewCell({ model, cell, x, y, msg }) {
  const style = cell !== -1 ? {
    backgroundColor: `hsl(${cell}, 50%, 50%)`
  } : {}

  return div({
    style,
    className: `border w-8 h-8`,
    onClick: () => msg({ f: toggl, args: { x, y } } )
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
