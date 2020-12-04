export { init, view }

// Init
function init({ msg }) {
  const socket = io('/', {
    transports: [ 'websocket' ]
  })

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

      msg(gotChat, { message })
    }
    catch (e) {
      console.error(e)
    }
  
  })


  socket.on('board', value => {
    const board = JSON.parse(value)
    msg(gotBoard, { board })
  })

  window.setTimeout(() => {
    msg(gotTick)
  }, 3000)


  return {
    board: null,
    socket,
    items: [],
    newChat: ''
  }
}

// Update
function gotBoard({ model, args: { board } }) {
  model.board = board
  return { model }
}

function toggl({ model, args: { x, y } }) {
  model.socket.emit('toggl', JSON.stringify({ x, y }))
  return { model }
}

function gotTick({ model, msg }) {
  model.socket.emit('tick')

  window.setTimeout(() => {
    msg(gotTick)
  }, 3000)

  return { model }
}

function gotChat({ model, args: { message } }) {
  model.items.push({
    id: uuidv4(),
    content: message.content,
    color: message.color
  })

  return { model }
}

function chatInput({ model, args: { value } }) {
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
          msg(chatSubmit)
        },
        className: 'flex'
      }, [
        input({
          onInput: e => msg(chatInput, { value: e.target.value }),
          className: 'p-1 shadow',
          value: model.newChat
        }),
        button({
          className: 'rounded p-2 shadow'
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
    onClick: () => msg(toggl, { x, y })
  }, [])
}

function viewItem({ model, item }) {
  return div({
    style: {
      color: `hsl(${item.color}, 50%, 30%)`
    },
    className: 'flex'
  }, [
    div({
      className: 'flex flex-col justify-center'
    }, [
      div({
        style: {
          backgroundColor: `hsl(${item.color}, 50%, 50%)`
        },
        className: 'w-4 h-4 mr-1'
      }, [])
    ]),
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
